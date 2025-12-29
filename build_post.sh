#!/bin/sh
set -eu

# usage: build_post.sh <src_md> <dst_html> <dst_root> "<pandoc_opt...>"
src="$1"
dst="$2"
dst_root="$3"
pandoc_opt="$4"

# Ensure output dir exists
out_dir="$(dirname "$dst")"
mkdir -p "$out_dir"

# Main html
# shellcheck disable=SC2086
pandoc -s "$src" -o "$dst" $pandoc_opt --template=post.html

# Redirect page: <dst without .html>/index.html
page_dir="${dst%.html}"
mkdir -p "$page_dir"

url="${dst#"$dst_root"}"
# shellcheck disable=SC2086
pandoc -s "$src" -o "$page_dir/index.html" $pandoc_opt --metadata "url=$url" --template=redirect.html
