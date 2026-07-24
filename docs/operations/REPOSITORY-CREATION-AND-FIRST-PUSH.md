# Repository Creation and First Push

## Canonical Naming

- Product: Iron Atlas
- Repository: `iron-atlas`
- Organization: `Iron-Signal-Systems`
- Canonical URL after creation: `https://github.com/Iron-Signal-Systems/atlas`
- Active development branch: `dev`

## Create the Empty Repository

Create a private empty repository in the Iron Signal Systems organization. Do not initialize it with a README, license, or `.gitignore`, because those files already exist in this package.

## Initialize and Push

From the extracted `iron-atlas` directory:

```bash
git init
git checkout -b dev
git add .
git diff --cached --check
git commit -m "establish Iron Atlas phase 0 foundation"
git remote add origin git@github.com:Iron-Signal-Systems/atlas.git
git push -u origin dev
```

After reviewing the repository on GitHub, configure `dev` as the default branch while active development remains pre-release.

## Validate Before Commit

```bash
./tools/validation/phase-gates/validate_phase0_step1.sh
```

Do not commit raw firewall backups, Cisco technical-support output, credentials, or unredacted infrastructure evidence.
