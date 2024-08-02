#!/bin/bash

# This script turns the presentations into static sites.

set -euo pipefail

nvm use
npm run tailwind:min

declare -a decks=(
  "01-architecture-101"
  "02-llms-101"
  "03-gaming-and-security-102"
  "04-cyber-investigations-101"
  "05-cyber-security-and-ai-102"
)


## now loop through the above array
for deck in "${decks[@]}"
do
  title=(jq -r '.title' "$deck/metadata.json")
  rm -rf "./dist/$deck"
  mkdir "./dist/$deck"
  sed \
    -e '/<!-- START REMOVE -->/,/<!-- END REMOVE -->/d' \
    -e "s/{{.Title}}/$title/g" \
    -e 's/out.css/out.min.css/g' \
    index.html > ".dist/$deck/index.html"
  cp "$deck/slides.md" ".dist/$deck/"
  cp "out.min.css" ".dist/$deck/out.min.css"
  cp "remark-latest.min.js" ".dist/$deck/remark.js"
done

