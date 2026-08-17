#!/usr/bin/env bash
# Attach a file to an existing Forgejo release. Env: FORGE_TOKEN,
# GITHUB_SERVER_URL, GITHUB_REPOSITORY, and RELEASE_TAG or GITHUB_REF_NAME.
set -euo pipefail

file="$1"
name="$(basename "$file")"
api="${GITHUB_SERVER_URL:?}/api/v1"
repo="${GITHUB_REPOSITORY:?}"
tag="${RELEASE_TAG:-${GITHUB_REF_NAME:?}}"
auth="Authorization: token ${FORGE_TOKEN:?}"

rel_id=$(curl -sSf -H "$auth" "$api/repos/$repo/releases/tags/$tag" \
  | grep -oE '"id":[0-9]+' | head -1 | grep -oE '[0-9]+')
: "${rel_id:?could not resolve release id for tag $tag}"

echo "attaching $name -> $repo release $tag (id $rel_id)"
curl -sSf -H "$auth" -X POST \
  -F "attachment=@${file}" \
  "$api/repos/$repo/releases/$rel_id/assets?name=$name" \
  -o /dev/null -w "  http %{http_code}\n"
