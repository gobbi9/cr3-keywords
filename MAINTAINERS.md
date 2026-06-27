# Maintainers

## Release flow

```bash
git push origin main
make release-tag VERSION=X.Y.Z
```

Tagging (`vX.Y.Z`) triggers `.github/workflows/release.yml`, which:

1. Builds binaries for macOS, Linux, and Windows.
2. Adds `docs/man/cr3.1` to release assets.
3. Publishes release assets on GitHub Releases.
4. Generates package manager manifests with SHA256:
   - `packaging/generated/cr3-keywords.rb`
   - `packaging/generated/cr3-keywords.json`
5. Uploads generated manifests as a workflow artifact.
6. Optionally pushes manifests to Homebrew tap and Scoop bucket repos (if configured).

## Package manager distribution

### Templates in this repo

- `packaging/homebrew/cr3-keywords.rb`
- `packaging/scoop/cr3-keywords.json`

### Local generation (manual)

```bash
# from repo root (required)
uv run scripts/generate_packaging.py --version X.Y.Z
```

Outputs:

- `packaging/generated/cr3-keywords.rb`
- `packaging/generated/cr3-keywords.json`

Requirements for generation:

- `dist/cr3_<version>_darwin-arm64`
- `dist/cr3_<version>_darwin-amd64`
- `dist/cr3_<version>_linux-amd64`
- `dist/cr3_<version>_windows-amd64.exe`
- `dist/cr3.1` (man page asset copied from `docs/man/cr3.1`)

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
