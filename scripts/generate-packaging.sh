#!/usr/bin/env sh
set -eu

# Generates packaging manifests from downloaded release artifacts.
# Usage:
#   ./scripts/generate-packaging.sh <version>
# Example:
#   ./scripts/generate-packaging.sh 0.1.0

if [ "${1:-}" = "" ]; then
  echo "Usage: $0 <version-without-v-prefix>"
  exit 1
fi

VERSION="$1"
DIST_DIR="dist"
OUT_DIR="packaging/generated"

if [ ! -f "packaging/homebrew/cr3-keywords.rb.tmpl" ] || [ ! -f "packaging/scoop/cr3-keywords.json.tmpl" ]; then
  echo "Run this script from repository root (cr3-keywords)."
  exit 1
fi

DARWIN_ARM64_FILE="${DIST_DIR}/cr3_${VERSION}_darwin-arm64"
DARWIN_AMD64_FILE="${DIST_DIR}/cr3_${VERSION}_darwin-amd64"
LINUX_AMD64_FILE="${DIST_DIR}/cr3_${VERSION}_linux-amd64"
WINDOWS_AMD64_FILE="${DIST_DIR}/cr3_${VERSION}_windows-amd64.exe"

for f in "$DARWIN_ARM64_FILE" "$DARWIN_AMD64_FILE" "$LINUX_AMD64_FILE" "$WINDOWS_AMD64_FILE"; do
  if [ ! -f "$f" ]; then
    echo "Missing artifact: $f"
    echo "Download release assets into '${DIST_DIR}/' first."
    exit 1
  fi
done

sha256_file() {
  file_path="$1"
  if command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$file_path" | awk '{print $1}'
  elif command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$file_path" | awk '{print $1}'
  else
    echo "No sha256 tool found (need shasum or sha256sum)."
    exit 1
  fi
}

SHA256_DARWIN_ARM64="$(sha256_file "$DARWIN_ARM64_FILE")"
SHA256_DARWIN_AMD64="$(sha256_file "$DARWIN_AMD64_FILE")"
SHA256_LINUX_AMD64="$(sha256_file "$LINUX_AMD64_FILE")"
SHA256_WINDOWS_AMD64="$(sha256_file "$WINDOWS_AMD64_FILE")"

mkdir -p "$OUT_DIR"

sed \
  -e "s/{{VERSION}}/${VERSION}/g" \
  -e "s/{{SHA256_DARWIN_ARM64}}/${SHA256_DARWIN_ARM64}/g" \
  -e "s/{{SHA256_DARWIN_AMD64}}/${SHA256_DARWIN_AMD64}/g" \
  -e "s/{{SHA256_LINUX_AMD64}}/${SHA256_LINUX_AMD64}/g" \
  "packaging/homebrew/cr3-keywords.rb.tmpl" \
  > "${OUT_DIR}/cr3-keywords.rb"

sed \
  -e "s/{{VERSION}}/${VERSION}/g" \
  -e "s/{{SHA256_WINDOWS_AMD64}}/${SHA256_WINDOWS_AMD64}/g" \
  "packaging/scoop/cr3-keywords.json.tmpl" \
  > "${OUT_DIR}/cr3-keywords.json"

echo "Generated:"
echo "  ${OUT_DIR}/cr3-keywords.rb"
echo "  ${OUT_DIR}/cr3-keywords.json"
