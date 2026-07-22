#!/usr/bin/env python3
from __future__ import annotations
import argparse, hashlib, json, re, sys
from datetime import datetime
from pathlib import Path
from typing import Any

SHA=re.compile(r'^[0-9a-f]{64}$'); IDENT=re.compile(r'^[a-z0-9][a-z0-9._-]+$')
JWT=re.compile(r'(?<![A-Za-z0-9_-])[A-Za-z0-9_-]{16,}\.[A-Za-z0-9_-]{16,}\.[A-Za-z0-9_-]{16,}(?![A-Za-z0-9_-])')
PEM_DASHES = "-" * 5
PRIVATE_KEY_LABELS = (
    ("PRIVATE", "KEY"),
    ("RSA", "PRIVATE", "KEY"),
    ("EC", "PRIVATE", "KEY"),
    ("OPENSSH", "PRIVATE", "KEY"),
)
PRIVATE = tuple(
    f"{PEM_DASHES}BEGIN {' '.join(parts)}{PEM_DASHES}"
    for parts in PRIVATE_KEY_LABELS
)
FORBIDDEN={'password','client_secret','client_assertion','client_id','access_token','refresh_token','id_token','authorization_code','code_verifier','cookie','set_cookie','session_cookie','private_key','totp_seed','recovery_code','sub','sid','jti','email','preferred_username','username','name','given_name','family_name','phone_number'}
TOP={'schema_version','bundle_id','evidence_class','claim_status','compatibility_claim','provider','capture','endpoints','artifacts','scenarios','redaction','nonclaims'}
PROVIDER={'label','software','version','configuration_sha256'}; CAPTURE={'captured_at','environment','tool','tool_version','sanitized','raw_artifacts_committed'}
ENDPOINTS={'issuer','discovery_artifact','jwks_artifact'}; ART={'path','kind','sha256'}
SCENARIO={'scenario_id','purpose','authentication_succeeded','literal_claims','claims_artifact','limitations'}; CLAIMS={'acr','amr','auth_time'}; REDACTION={'prohibited_material_absent','review_status','notes'}
KINDS={'discovery-document','jwks-document','literal-assurance-claims','sanitized-observation'}
class V(Exception): pass
def req(c,m):
 if not c: raise V(m)
def exact(v,k,w): req(isinstance(v,dict),f'{w} must be an object'); req(set(v)==k,f'{w} keys differ'); return v
def string(v,w,n=1,x=512): req(isinstance(v,str) and n<=len(v)<=x,f'{w} must be a bounded string'); return v
def digest(v,w): req(isinstance(v,str) and SHA.fullmatch(v),f'{w} must be a lowercase SHA-256'); return v
def scan_keys(v,w):
 if isinstance(v,dict):
  for k,x in v.items():
   req(k.lower().replace('-','_') not in FORBIDDEN,f'{w} contains prohibited key {k!r}'); scan_keys(x,f'{w}.{k}')
 elif isinstance(v,list):
  for i,x in enumerate(v): scan_keys(x,f'{w}[{i}]')
def scan_text(t,w):
 req(not JWT.search(t),f'{w} contains a JWT-shaped value')
 for m in PRIVATE: req(m not in t,f'{w} contains private-key material')
def claims(v,w):
 v=exact(v,CLAIMS,w); req(v['acr'] is None or isinstance(v['acr'],str),f'{w}.acr invalid'); req(isinstance(v['amr'],list) and len(v['amr'])<=32,f'{w}.amr invalid')
 for i,m in enumerate(v['amr']): string(m,f'{w}.amr[{i}]',1,128)
 a=v['auth_time']; req(a is None or (isinstance(a,int) and not isinstance(a,bool) and a>=0),f'{w}.auth_time invalid'); return v
def safe(root,rel,w):
 p=Path(rel); req(not p.is_absolute() and '..' not in p.parts,f'{w} contains path traversal'); req(rel.endswith('.sanitized.json'),f'{w} must end in .sanitized.json')
 q=(root/p).resolve(); req(q.is_relative_to(root.resolve()),f'{w} escapes bundle'); req(q.is_file() and not q.is_symlink(),f'{w} missing or symlinked'); return q

def validate(bundle:Path)->str:
 req(bundle.is_file(),f'bundle missing: {bundle}'); root=bundle.parent.resolve(); raw=bundle.read_text(); scan_text(raw,str(bundle))
 try: d=json.loads(raw)
 except json.JSONDecodeError as e: raise V(f'invalid bundle JSON: {e}')
 scan_keys(d,'bundle'); d=exact(d,TOP,'bundle'); req(d['schema_version']==1,'schema_version must be 1'); bid=string(d['bundle_id'],'bundle_id',8,128); req(IDENT.fullmatch(bid),'bundle_id invalid')
 req(d['evidence_class'] in {'synthetic','controlled-sanitized'},'evidence_class unsupported'); req(d['claim_status']=='observation-only','claim_status must be observation-only'); req(d['compatibility_claim'] is False,'compatibility_claim must be false')
 p=exact(d['provider'],PROVIDER,'provider'); [string(p[k],f'provider.{k}',1,128) for k in ('label','software','version')]; digest(p['configuration_sha256'],'provider.configuration_sha256')
 c=exact(d['capture'],CAPTURE,'capture'); string(c['captured_at'],'capture.captured_at',20,64); req(c['captured_at'].endswith('Z'),'captured_at must end in Z')
 try: datetime.fromisoformat(c['captured_at'][:-1]+'+00:00')
 except ValueError as e: raise V('captured_at invalid') from e
 req(c['environment']=='disposable','capture environment must be disposable'); string(c['tool'],'capture.tool',1,128); string(c['tool_version'],'capture.tool_version',1,64); req(c['sanitized'] is True,'capture.sanitized must be true'); req(c['raw_artifacts_committed'] is False,'raw_artifacts_committed must be false')
 e=exact(d['endpoints'],ENDPOINTS,'endpoints'); req(string(e['issuer'],'endpoints.issuer',8,512).startswith('https://'),'issuer must use HTTPS')
 req(isinstance(d['artifacts'],list) and 1<=len(d['artifacts'])<=128,'artifacts invalid'); amap={}
 for i,r in enumerate(d['artifacts']):
  r=exact(r,ART,f'artifacts[{i}]'); rel=string(r['path'],f'artifacts[{i}].path',1,256); req(rel not in amap,f'duplicate artifact {rel}'); req(r['kind'] in KINDS,f'artifact kind invalid: {rel}'); expected=digest(r['sha256'],f'artifacts[{i}].sha256'); q=safe(root,rel,f'artifacts[{i}].path'); b=q.read_bytes()
  try: text=b.decode(); obj=json.loads(text)
  except (UnicodeDecodeError,json.JSONDecodeError) as x: raise V(f'artifact invalid: {rel}: {x}')
  scan_text(text,rel); scan_keys(obj,rel); req(hashlib.sha256(b).hexdigest()==expected,f'artifact digest mismatch: {rel}'); amap[rel]=(r['kind'],q)
 req(e['discovery_artifact'] in amap and amap[e['discovery_artifact']][0]=='discovery-document','discovery reference invalid'); req(e['jwks_artifact'] in amap and amap[e['jwks_artifact']][0]=='jwks-document','JWKS reference invalid')
 req(isinstance(d['scenarios'],list) and 1<=len(d['scenarios'])<=64,'scenarios invalid'); ids=set()
 for i,s in enumerate(d['scenarios']):
  s=exact(s,SCENARIO,f'scenarios[{i}]'); sid=string(s['scenario_id'],f'scenarios[{i}].scenario_id',3,128); req(IDENT.fullmatch(sid),f'scenario id invalid: {sid}'); req(sid not in ids,f'duplicate scenario_id {sid}'); ids.add(sid); string(s['purpose'],f'scenarios[{i}].purpose'); req(isinstance(s['authentication_succeeded'],bool),'authentication_succeeded invalid'); literal=claims(s['literal_claims'],f'scenarios[{i}].literal_claims'); ref=string(s['claims_artifact'],f'scenarios[{i}].claims_artifact',1,256); req(ref in amap and amap[ref][0]=='literal-assurance-claims',f'claims reference invalid: {ref}'); req(json.loads(amap[ref][1].read_text())==literal,f'scenario claims differ from artifact: {ref}'); req(isinstance(s['limitations'],list) and 1<=len(s['limitations'])<=32,'limitations invalid'); [string(x,f'limitation[{j}]') for j,x in enumerate(s['limitations'])]
 r=exact(d['redaction'],REDACTION,'redaction'); req(r['prohibited_material_absent'] is True,'prohibited_material_absent must be true'); req(r['review_status'] in {'synthetic','reviewed'},'review_status invalid'); req(isinstance(r['notes'],list) and len(r['notes'])<=32,'redaction notes invalid'); [string(x,f'redaction.notes[{i}]') for i,x in enumerate(r['notes'])]
 req(isinstance(d['nonclaims'],list) and 4<=len(d['nonclaims'])<=32,'nonclaims invalid'); n=' '.join(str(x).lower() for x in d['nonclaims'])
 for phrase in ('not a compatibility claim','not production readiness','no provider-semantic inference','no raw tokens or secrets'): req(phrase in n,f'nonclaims missing: {phrase}')
 return f"PASS: representative-provider evidence bundle {bid}; artifacts={len(amap)} scenarios={len(ids)} class={d['evidence_class']} claim={d['claim_status']}"
def main():
 a=argparse.ArgumentParser(); a.add_argument('bundle',type=Path); p=a.parse_args()
 try: print(validate(p.bundle.resolve())); return 0
 except (OSError,V) as e: print(f'FAIL: {e}',file=sys.stderr); return 1
if __name__=='__main__': raise SystemExit(main())
