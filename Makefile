SRC := .
DST := public

.PHONY: build clean serve

ifndef CI
BUILD_OPTS=--base-url 'http://0.0.0.0:4507'
endif

build:
	deno run --allow-read --allow-write --unstable build.ts $(BUILD_OPTS)

fmt:
	deno fmt build.ts

bin/file_server:
	deno install --allow-net --allow-read --root . https://deno.land/std/http/file_server.ts

serve: bin/file_server
	./bin/file_server $(DST)

clean:
	rm -rf $(DST)/**
