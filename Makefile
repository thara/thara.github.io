DST=public2
TMP_DIR=tmp
DATA_DIR=.pandoc

ifdef CI
	BASE_URL=https://thara.dev
else
	BASE_URL=http://localhost:8000
endif

YEAR=$(shell date +%Y)

PANDOC_OPT=-t html -f gfm+yaml_metadata_block -V base_url=$(BASE_URL) -V year=$(YEAR) --data-dir=$(DATA_DIR)

MD_FILES=$(shell find ./pages -type f | sed 's/^\.\/pages/$(DST)/g')
HTML_FILES=$(MD_FILES:.md=.html)

POST_MD_FILES=$(shell ls -d $$(find ./pages/posts -type f) | sed 's/^\.\/pages/$(DST)/g')

all: $(HTML_FILES)
	@cp -r assets/ $(DST)

$(DST)/index.html: pages/index.md
	@mkdir -p $(dir $@)
	@pandoc -s $< -o $@ $(PANDOC_OPT) --template=template

$(DST)/posts.html: pages/posts.md
	@mkdir -p $(TMP_DIR)
	@cp -f pages/posts.md $(TMP_DIR)/posts.md
	@for f in $(shell find ./pages/posts -type f | sort -r); do \
		url=`echo $$f | sed "s/.md/.html/" | sed "s/^.\/pages\///"`;\
		pandoc -s $$f --template=template-listitem-link.md --metadata url=$$url --data-dir=$(DATA_DIR) >> $(TMP_DIR)/posts.md; \
	done
	@mkdir -p $(dir $@)
	@pandoc -s $(TMP_DIR)/posts.md -o $@ $(PANDOC_OPT) --template=template

$(DST)/%.html: pages/%.md
	@mkdir -p $(dir $@)
	@pandoc -s $< -o $@ $(PANDOC_OPT) --template=template-post
	@page_dir=`echo $@ | sed "s/\.html//"`; \
	mkdir -p $$page_dir ;\
	url=`echo $@ | sed "s/^$(DST)//"`; \
	pandoc -s $< -o $$page_dir/index.html $(PANDOC_OPT) --metadata url=$$url --template=template-redirect.html

.PHONY: clean
clean:
	@rm -rf $(DST)
	@rm -rf $(TMP_DIR)

.PHONY: serve
serve: all
	@cd $(DST); python3 -m http.server
