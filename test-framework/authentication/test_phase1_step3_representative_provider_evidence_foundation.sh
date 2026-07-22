#!/usr/bin/env bash
set -Eeuo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"; cd "$root"
v=tools/validation/validate_representative_provider_evidence_bundle.py
fixture=validation/fixtures/authentication/representative-provider-evidence-foundation/synthetic-v1
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT
python3 "$v" "$fixture/bundle.json" >"$tmp/a"; python3 "$v" "$fixture/bundle.json" >"$tmp/b"; cmp "$tmp/a" "$tmp/b"; grep -F 'claim=observation-only' "$tmp/a" >/dev/null
printf 'PASS: valid synthetic bundle and deterministic output\n'
case_copy(){ cp -a "$fixture" "$tmp/$1"; printf '%s\n' "$tmp/$1"; }
expect_fail(){ if python3 "$v" "$2" >"$tmp/$1.out" 2>"$tmp/$1.err"; then echo "FAIL: $1 unexpectedly validated" >&2; exit 1; fi; grep -F 'FAIL:' "$tmp/$1.err" >/dev/null; echo "PASS: rejected $1"; }
mutate(){ python3 - "$1" "$2" <<'PY'
import copy,hashlib,json,sys
p,kind=sys.argv[1:]; d=json.load(open(p))
if kind=='compat': d['compatibility_claim']=True
elif kind=='raw': d['capture']['raw_artifacts_committed']=True
elif kind=='secret': d['provider']['client_secret']='not-allowed'
elif kind=='jwt': d['redaction']['notes'].append('aaaaaaaaaaaaaaaa.bbbbbbbbbbbbbbbb.cccccccccccccccc')
elif kind=='traversal': d['artifacts'][0]['path']='../escape.sanitized.json'
elif kind=='digest': d['artifacts'][0]['sha256']='0'*64
elif kind=='duplicate': d['scenarios'].append(copy.deepcopy(d['scenarios'][0]))
elif kind=='interpret': d['scenarios'][0]['mfa_satisfied']=True
elif kind=='disagree': d['scenarios'][1]['literal_claims']['amr']=['pwd']
open(p,'w').write(json.dumps(d,indent=2,sort_keys=True)+'\n')
PY
}
for spec in 'compat compatibility-claim' 'raw raw-artifacts' 'secret credential-key' 'jwt jwt-shaped' 'traversal path-traversal' 'digest digest-mismatch' 'duplicate duplicate-scenario' 'interpret interpreted-policy' 'disagree claims-disagreement'; do set -- $spec; d="$(case_copy "$2")"; mutate "$d/bundle.json" "$1"; expect_fail "$2" "$d/bundle.json"; done
d="$(case_copy private-key)"; python3 - "$d/bundle.json" "$d/artifacts/claims-no-assurance.sanitized.json" <<'PY'
import hashlib,json,sys
b,a=sys.argv[1:]; x=json.load(open(a)); x['note']='-----BEGIN PRIVATE KEY-----'; open(a,'w').write(json.dumps(x,indent=2,sort_keys=True)+'\n'); d=json.load(open(b)); h=hashlib.sha256(open(a,'rb').read()).hexdigest()
for r in d['artifacts']:
 if r['path']=='artifacts/claims-no-assurance.sanitized.json': r['sha256']=h
open(b,'w').write(json.dumps(d,indent=2,sort_keys=True)+'\n')
PY
expect_fail private-key "$d/bundle.json"
printf 'PASS: representative-provider evidence foundation regression\n'
