#!/usr/bin/env sh
set -eu

# Publishes a generated packaging manifest to an external repository.
# Usage:
#   PACKAGING_PUSH_TOKEN=<token> ./scripts/publish-packaging-manifest.sh \
#     <target-repo> <source-file> <target-file> <display-name> <commit-message>
# Example:
#   PACKAGING_PUSH_TOKEN=... ./scripts/publish-packaging-manifest.sh \
#     gobbi9/tap packaging/generated/cr3-keywords.rb Formula/cr3-keywords.rb \
#     "Homebrew formula" "Update cr3-keywords formula for v0.1.0"

if [ "${PACKAGING_PUSH_TOKEN:-}" = "" ]; then
  echo "PACKAGING_PUSH_TOKEN env var is required."
  exit 1
fi

if [ "$#" -ne 5 ]; then
  echo "Usage: $0 <target-repo> <source-file> <target-file> <display-name> <commit-message>"
  exit 1
fi

TARGET_REPO="$1"
SOURCE_FILE="$2"
TARGET_FILE="$3"
DISPLAY_NAME="$4"
COMMIT_MESSAGE="$5"

if [ ! -f "$SOURCE_FILE" ]; then
  echo "Source file not found: $SOURCE_FILE"
  exit 1
fi

git config --global user.name "github-actions[bot]"
git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

echo "Publishing ${DISPLAY_NAME} to ${TARGET_REPO}"
git clone "https://x-access-token:${PACKAGING_PUSH_TOKEN}@github.com/${TARGET_REPO}.git" "${WORK_DIR}/repo"

TARGET_DIR="$(dirname "$TARGET_FILE")"
mkdir -p "${WORK_DIR}/repo/${TARGET_DIR}"
cp "$SOURCE_FILE" "${WORK_DIR}/repo/${TARGET_FILE}"

(
  cd "${WORK_DIR}/repo"
  if git diff --quiet -- "$TARGET_FILE"; then
    echo "No ${DISPLAY_NAME} changes to commit."
  else
    git add "$TARGET_FILE"
    git commit -m "$COMMIT_MESSAGE"
    git push
  fi
)
