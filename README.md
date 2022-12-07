# thara.github.io

My personal web site, hosting at https://thara.dev

[![build](https://github.com/thara/thara.github.io/actions/workflows/publish.yml/badge.svg)](https://github.com/thara/thara.github.io/actions/workflows/publish.yml)

## Requirements

- go 1.x later
- make

## Usage

### Serve the generated site locally

`make serve`

### Deploy via GitHub Action

The GitHub Action ([publish.yml](./.github/workflows/publish.yml)) is available to
build and deploy automatially.

- `master` branch (Source)
- `gh-pages` branch (GitHub Pages)

## Note

### Markdown frontmatter

| key   | description |
| ----  | ------------|
| title | page title  |
| date  | page created timestamp (used as order in post list) |
| path  | overwrite destination path |

### Templating

- pages/posts/*.md -> templates/post.html
- pages/*.md -> templates/base.html

## Author

Tomochika Hara
