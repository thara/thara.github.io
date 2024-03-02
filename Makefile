DST=public
TMP_DIR=.tmp
DATA_DIR=$(shell pwd)

ifdef CI
	BASE_URL=https://thara.dev
else
	BASE_URL=http://localhost:8000
endif

YEAR=$(shell date +%Y)

PANDOC_OPT=-t html -f gfm+yaml_metadata_block -V base_url=$(BASE_URL) -V year=$(YEAR) --data-dir=$(DATA_DIR)

MD_FILES=$(shell find ./pages -type f \
				 	| sed 's/^\.\/pages/$(DST)/g' \
				 	| sed "s/posts\/20[0-9][0-9]-[0-1][0-9]-[0-3][0-9]-/posts\//")
HTML_FILES=$(MD_FILES:.md=.html)

debug:
	@echo $(HTML_FILES)

all: $(HTML_FILES)
	@cp -r assets/ $(DST)

$(DST)/posts.html: pages/posts.md
	@mkdir -p $(TMP_DIR)
	@cp -f pages/posts.md $(TMP_DIR)/posts.md
	@for f in $(shell find ./pages/posts -type f | sort -r); do \
		url=`echo $$f | sed "s/.md/.html/" | sed "s/^.\/pages\///" | sed "s/20[0-9][0-9]-[0-1][0-9]-[0-3][0-9]-//"`;\
		pandoc -s $$f --template=list-item-link.md --metadata url=$$url --data-dir=$(DATA_DIR) >> $(TMP_DIR)/posts.md; \
	done
	@mkdir -p $(dir $@)
	@pandoc -s $(TMP_DIR)/posts.md -o $@ $(PANDOC_OPT) --template=base

$(DST)/posts/%.html: pages/posts/20??-??-??-%.md
	@mkdir -p $(dir $@)
	@pandoc -s $< -o $@ $(PANDOC_OPT) --template=post.html
	@page_dir=`echo $@ | sed "s/\.html//"`; \
	mkdir -p $$page_dir ;\
	url=`echo $@ | sed "s/^$(DST)//"`; \
	pandoc -s $< -o $$page_dir/index.html $(PANDOC_OPT) --metadata url=$$url --template=redirect.html

$(DST)/%.html: pages/%.md
	@mkdir -p $(dir $@)
	@pandoc -s $< -o $@ $(PANDOC_OPT) --template=base

.PHONY: clean
clean:
	@rm -rf $(DST)
	@rm -rf $(TMP_DIR)

.PHONY: serve
serve: all
	@cd $(DST); python3 -m http.server
