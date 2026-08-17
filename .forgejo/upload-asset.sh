#!/usr/bin/env bash
# Attach a file to a Forgejo release, creating the release if needed.
# Replaces a same-named asset (409) so re-runs succeed.
# Env: FORGE_TOKEN, GITHUB_SERVER_URL, GITHUB_REPOSITORY,
#      RELEASE_TAG or GITHUB_REF_NAME.
set -euo pipefail

file="$1"
name="$(basename "$file")"
api="${GITHUB_SERVER_URL:?}/api/v1"
repo="${GITHUB_REPOSITORY:?}"
tag="${RELEASE_TAG:-${GITHUB_REF_NAME:?}}"
auth="Authorization: token ${FORGE_TOKEN:?}"

tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT

code=$(curl -sS -o "$tmp" -w '%{http_code}' -H "$auth" "$api/repos/$repo/releases/tags/$tag")
if [ "$code" = "404" ]; then
  curl -sSf -H "$auth" -H 'Content-Type: application/json' \
    -d "{\"tag_name\":\"$tag\",\"name\":\"$tag\"}" \
    "$api/repos/$repo/releases" > "$tmp"
elif [ "$code" != "200" ]; then
  echo "GET release $tag failed: HTTP $code" >&2
  cat "$tmp" >&2
  exit 1
fi

rel_id=$(grep -oE '"id":[0-9]+' "$tmp" | head -1 | grep -oE '[0-9]+')
: "${rel_id:?could not resolve release id for tag $tag}"

delete_named_asset() {
  local assets ids
  assets=$(curl -sSf -H "$auth" "$api/repos/$repo/releases/$rel_id/assets")
  ids=$(NAME="$name" node -e '
    const fs = require("fs");
    const want = process.env.NAME;
    let list = [];
    try { list = JSON.parse(fs.readFileSync(0, "utf8")); } catch {}
    for (const a of Array.isArray(list) ? list : []) {
      if (a && a.name === want && a.id != null) console.log(a.id);
    }
  ' <<<"$assets")
  for aid in $ids; do
    curl -sSf -H "$auth" -X DELETE "$api/repos/$repo/releases/$rel_id/assets/$aid" -o /dev/null || true
  done
}

echo "attaching $name -> $repo release $tag (id $rel_id)"
http=$(curl -sS -o "$tmp" -w '%{http_code}' -H "$auth" -X POST \
  -F "attachment=@${file}" \
  "$api/repos/$repo/releases/$rel_id/assets?name=$name")
if [ "$http" = "409" ]; then
  echo "asset exists; replacing"
  delete_named_asset
  http=$(curl -sS -o "$tmp" -w '%{http_code}' -H "$auth" -X POST \
    -F "attachment=@${file}" \
    "$api/repos/$repo/releases/$rel_id/assets?name=$name")
fi
if [ "$http" != "201" ] && [ "$http" != "200" ]; then
  echo "upload $name failed: HTTP $http" >&2
  cat "$tmp" >&2
  exit 1
fi
echo "  http $http"
