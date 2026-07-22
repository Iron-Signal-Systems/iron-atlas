#!/usr/bin/env python3
from pathlib import Path
import subprocess,sys
R=Path(__file__).resolve().parents[2]; BASE='e7824049852855f15d26686600fc42802b8a38ff'; p=f=0
def check(n,c):
 global p,f
 print(('PASS' if c else 'FAIL')+': '+n); p+=bool(c); f+=not c
def txt(x): return (R/x).read_text()
def compact(value): return " ".join(value.split())
required=['docs/architecture/REPRESENTATIVE-PROVIDER-EVIDENCE-FOUNDATION.md','docs/testing/REPRESENTATIVE-PROVIDER-EVIDENCE-FOUNDATION-TESTING.md','validation/schemas/representative-provider-evidence-bundle.schema.json','validation/fixtures/authentication/representative-provider-evidence-foundation/synthetic-v1/README.md','validation/fixtures/authentication/representative-provider-evidence-foundation/synthetic-v1/bundle.json','tools/validation/validate_representative_provider_evidence_bundle.py','tools/validation/validate_phase1_step3_representative_provider_evidence_foundation.py','test-framework/authentication/test_phase1_step3_representative_provider_evidence_foundation.sh','tools/validation/phase-gates/validate_phase1_step3_representative_provider_evidence_foundation.sh']
for x in required: check('required file '+x,(R/x).is_file())
a=txt(required[0]); t=txt(required[1]); s=txt(required[2]); g=txt(required[-1]); run=txt('test-framework/run_all.sh'); repo=txt('tools/validation/validate_repository.sh'); road=txt('docs/roadmap/IMPLEMENTATION-ROADMAP.md'); plan=txt('docs/roadmap/PHASE-GATE-PLAN.md'); trace=txt('docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md'); sec=txt('docs/security/MFA-AND-AUTHENTICATION-ASSURANCE-REQUIREMENTS.md'); acc=txt('docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md'); read=txt('README.md'); docs=txt('docs/README.md'); change=txt('CHANGELOG.md')
for q in ('observation-only','compatibility_claim','Raw provider traffic','No named-provider compatibility gate exists','digest-bound observation bundle'): check('architecture contains '+q,q in a)
for q in ('JWT-shaped values','path traversal','deterministic validation output','does not establish compatibility'): check('testing contains '+q,q in compact(t))
for q in ('"compatibility_claim"','"const": false','"controlled-sanitized"','"literal_claims"','"raw_artifacts_committed"'): check('schema contains '+q,q in s)
check('signed provider-neutral boundary is an ancestor',subprocess.run(['git','merge-base','--is-ancestor',BASE,'HEAD'],cwd=R).returncode==0)
check('phase gate exact-revalidates signed provider-neutral boundary',BASE in g)
check('phase gate describes evidence foundation scope','evidence contract and sanitized observation boundary' in g)
check('complete framework registers foundation static validation','representative-provider evidence-foundation static validation' in run)
check('complete framework registers foundation regression','representative-provider evidence-foundation regression' in run)
check('complete framework exact-revalidates architecture alignment', '2347d21f779768f40496a93cb1d9140cc3b6e0ce' in run and 'revalidate_architecture_alignment_static_checkpoint' in run and 'revalidate_architecture_alignment_regression_checkpoint' in run)
check('complete framework exact-revalidates provider-neutral assurance evidence', 'e7824049852855f15d26686600fc42802b8a38ff' in run and 'revalidate_provider_neutral_static_checkpoint' in run and 'revalidate_provider_neutral_regression_checkpoint' in run)
check('complete framework does not run historical status validators on successor documentation', 'run "architecture and roadmap alignment static validation" python3' not in run and 'run "provider-neutral assurance-evidence static validation" python3' not in run)
check('repository validation registers foundation validator','validate_phase1_step3_representative_provider_evidence_foundation.py' in repo)
check('repository validation exact-revalidates historical status contracts', '2347d21f779768f40496a93cb1d9140cc3b6e0ce' in repo and 'e7824049852855f15d26686600fc42802b8a38ff' in repo and 'isolated_python_validator_revalidate' in repo)
helper=txt('tools/validation/lib/isolated_gate_revalidation.sh'); helper_test=txt('test-framework/phase-gates/test_isolated_gate_revalidation.sh')
check('isolated helper supports Python validators','isolated_python_validator_revalidate()' in helper)
check('isolated helper Python behavior is regression tested','successful isolated Python validator returns success' in helper_test and 'isolated Python validator rejects parent traversal' in helper_test)
check('testing records historical checkpoint isolation','Historical checkpoint isolation' in t)
check('roadmap identifies active foundation','Representative-provider evidence foundation is the active bounded Step 3 checkpoint' in road)
check('gate plan keeps compatibility planned','validate_phase1_step3_representative_provider_compatibility.sh' in plan)
check('gate plan identifies active foundation','validate_phase1_step3_representative_provider_evidence_foundation.sh' in plan)
check('traceability identifies foundation status','Representative-provider evidence-foundation implementation status' in trace)
check('security identifies foundation','Representative-provider evidence foundation' in sec)
check('acceptance template records foundation','Representative-provider evidence-foundation checkpoint' in acc)
check('README identifies active foundation','representative-provider evidence foundation is active' in read)
check('documentation index identifies active foundation','representative-provider evidence foundation is the active bounded' in docs)
check('changelog records foundation candidate','representative-provider evidence-foundation candidate' in change)
check('formal Step 3 acceptance record remains absent',not (R/'docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md').exists())
check('named-provider compatibility executable gate remains absent',not (R/'tools/validation/phase-gates/validate_phase1_step3_representative_provider_compatibility.sh').exists())
v=R/'tools/validation/validate_representative_provider_evidence_bundle.py'; b=R/'validation/fixtures/authentication/representative-provider-evidence-foundation/synthetic-v1/bundle.json'
check('committed synthetic evidence bundle validates',subprocess.run([sys.executable,str(v),str(b)],cwd=R).returncode==0)
manifest=set(txt('FILE-MANIFEST.txt').splitlines())
for x in required: check('repository registration contains '+x,x in manifest)
print(f'\nPASS checks: {p}\nFAIL checks: {f}'); raise SystemExit(bool(f))
