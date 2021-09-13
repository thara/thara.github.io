# thara.github.io

My personal web site

[![build](https://github.com/thara/thara.github.io/actions/workflows/publish.yml/badge.svg)](https://github.com/thara/thara.github.io/actions/workflows/publish.yml)

## Requirements

- deno 1.x later
- make

## Usage

### Build locally

`make build`

### Build via GitHub Action

The GitHub Action ([build.yml](./.github/workflows/build.yml)) is available to
build and deploy automatially.

- `main` branch (Source)
- `gh-pages` branch (GitHub Pages)

## Design

### Markdown frontmatter

| key   | description |
| ----  | ------------|
| title | page title  |
| date  | page created timestamp (used as order in post list) |
| path  | overwrite destination path |

### Templating

- posts/*.md -> templates/post.ejs
- XXXX.md -> tempaltes/XXX.ejs or default template(tempaltes/layout.ejs)

## Author

Tomochika Hara
