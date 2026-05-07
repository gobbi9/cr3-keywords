# Maintainers

## Release flow

```bash
git push origin main
make release-tag VERSION=0.1.0
```

Tagging (`vX.Y.Z`) triggers `.github/workflows/release.yml`, which:

1. Builds binaries for macOS, Linux, and Windows.
2. Publishes release assets on GitHub Releases.
3. Generates package manager manifests with SHA256:
   - `packaging/generated/cr3-keywords.rb`
   - `packaging/generated/cr3-keywords.json`
4. Uploads generated manifests as a workflow artifact.
5. Optionally pushes manifests to Homebrew tap and Scoop bucket repos (if configured).

## Package manager distribution

### Templates in this repo

- `packaging/homebrew/cr3-keywords.rb.tmpl`
- `packaging/scoop/cr3-keywords.json.tmpl`

### Local generation (manual)

```bash
# from repo root
chmod +x scripts/generate-packaging.sh
./scripts/generate-packaging.sh 0.1.0
```

Outputs:

- `packaging/generated/cr3-keywords.rb`
- `packaging/generated/cr3-keywords.json`

## Configure automatic publishing to external repos

Set these in this repository settings:

### Secrets

- `PACKAGING_PUSH_TOKEN`
  - GitHub PAT with write access to:
    - your Homebrew tap repo
    - your Scoop bucket repo

### Variables (optional)

- `HOMEBREW_TAP_REPO` (default: `gobbi9/tap`)
- `SCOOP_BUCKET_REPO` (default: `gobbi9/scoop-bucket`)

## External repo layout expectations

### Homebrew tap repository

Repo example: `gobbi9/tap`

- Formula path: `Formula/cr3-keywords.rb`

End-user install command:

```bash
brew tap gobbi9/tap https://github.com/gobbi9/tap
brew install cr3-keywords
```

### Scoop bucket repository

Repo example: `gobbi9/scoop-bucket`

- Manifest path: `cr3-keywords.json`

End-user install command:

```powershell
scoop bucket add gobbi9 https://github.com/gobbi9/scoop-bucket
scoop install gobbi9/cr3-keywords
```
