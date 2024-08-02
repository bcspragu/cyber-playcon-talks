#!/bin/bash

# This script turns the presentations into static sites.

set -euo pipefail

npm run tailwind:min

declare -A decks=(
  ["01-architecture-101"]="architecture-101.cybr.lol"
  ["02-llms-101"]="llms-101.cybr.lol"
  ["03-gaming-and-security-102"]="gaming-and-security-102.cybr.lol"
  ["04-cyber-investigations-101"]="cyber-investigations-101.cybr.lol"
  ["05-cyber-security-and-ai-102"]="cyber-security-and-ai-102.cybr.lol"
)

for srcdir in "${!decks[@]}"; do
  destdir="${decks[$srcdir]}"

  title="$(jq -r '.title' "$srcdir/metadata.json")"

  rm -rf "./dist/$destdir"
  mkdir "./dist/$destdir"
  sed \
    -e '/<!-- START REMOVE -->/,/<!-- END REMOVE -->/d' \
    -e "s/{{.Title}}/$title/g" \
    -e 's/out.css/out.min.css/g' \
    index.html > "./dist/$destdir/index.html"
  cp "$srcdir/slides.md" "./dist/$destdir/"
  cp -r "$srcdir/assets" "./dist/$destdir/"
  cp favicon/* "./dist/$destdir/"
  cp "out.min.css" "./dist/$destdir/out.min.css"
  cp "remark-latest.min.js" "./dist/$destdir/remark.js"
done
