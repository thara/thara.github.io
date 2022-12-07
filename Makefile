DST := public

.PHONY: build fmt serve clean post

build:
	go build .

run:
	go run .

serve:
	go run . --serve

fmt:
	deno fmt build.ts

clean:
	rm -rf $(DST)/**

post:
	@read -p "post title: " title; \
  (vim pages/posts/`date +%Y-%m-%d`-$$title.md)
