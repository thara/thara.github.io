default: serve

build:
  go build .

sitegen baseURL:
  go run . --base {{baseURL}}

serve:
  go run . --serve
