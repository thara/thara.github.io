SRC=pages
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

MD_FILES := $(shell find ./$(SRC) -type f \
				 	| sed 's/^\.\/$(SRC)/$(DST)/g' \
				 	| sed "s/posts\/20[0-9][0-9]-[0-1][0-9]-[0-3][0-9]-/posts\//")
HTML_FILES := $(MD_FILES:.md=.html)

POSTS_SRC := $(shell find posts -type f -name '20??-??-??-*.md')
post_slug = $(shell basename '$(1)' .md | sed -E 's/^20[0-9]{2}-[0-1][0-9]-[0-3][0-9]-//')

POST_HTML_FILES := $(foreach f,$(POSTS_SRC),$(DST)/posts/$(call post_slug,$(f)).html)

test_post_slug:
	@echo $(call post_slug,posts/2024-06-15-my-first-post.md)

all: $(POST_HTML_FILES) $(HTML_FILES)

define POST_RULE
$(DST)/posts/$(call post_slug,$(1)).html: $(1)
	./build_post.sh "$$<" "$$@" "$(DST)" "$(PANDOC_OPT)"
endef

$(foreach f,$(POSTS_SRC),$(eval $(call POST_RULE,$(f))))

$(DST)/posts.html: $(SRC)/posts.md $(POST_HTML_FILES)
	@mkdir -p $(TMP_DIR)
	@cp -f $(SRC)/posts.md $(TMP_DIR)/posts.md
	@for f in $(shell find ./posts -type f | sort -r); do \
		url=`echo $$f | sed "s/.md/.html/" | sed "s/20[0-9][0-9]-[0-1][0-9]-[0-3][0-9]-//"`;\
		pandoc -s $$f --template=list-item-link.md --metadata url=$$url --data-dir=$(DATA_DIR) >> $(TMP_DIR)/posts.md; \
	done
	@mkdir -p $(dir $@)
	pandoc -s $(TMP_DIR)/posts.md -o $@ $(PANDOC_OPT) --template=base

$(DST)/%.html: $(SRC)/%.md
	@mkdir -p $(dir $@)
	@pandoc -s $< -o $@ $(PANDOC_OPT) --template=base

.PHONY: clean
clean:
	@rm -f $(HTML_FILES)
	@rm -rf $(DST)/posts
	@rm -rf $(TMP_DIR)

.PHONY: serve
serve: all
	@cd $(DST); python3 -m http.server

.PHONY: debug
debug:
	@echo $(HTML_FILES)

.PHONE: post
post:
	@echo "Title: "
	@read title; \
	date=$(shell date +%Y-%m-%d); \
	echo "---\ntitle: $$title\ndate: '$$date'\npublished: '$$date'\n---" > posts/$$date-$$title.md
