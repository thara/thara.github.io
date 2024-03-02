DST=public2
DATA_DIR=.pandoc

ifdef CI
	BASE_URL=https://thara.dev
else
	BASE_URL=http://localhost:8000
endif

YEAR=$(shell date +%Y)

PANDOC_OPT=-t html -f gfm+yaml_metadata_block -V base_url=$(BASE_URL) -V year=$(YEAR) --data-dir=$(DATA_DIR)

MD_FILES=$(shell ls -d $$(find ./pages -type f) | sed 's/^\.\/pages/$(DST)/g')
HTML_FILES=$(MD_FILES:.md=.html)

md:
	echo $(MD_FILES)
	echo $(HTML_FILES)

all: $(HTML_FILES)
	@cp -r assets/ $(DST)

$(DST)/%.html: pages/%.md
	@mkdir -p $(dir $@)
	pandoc -s $< -o $@ $(PANDOC_OPT)

.PHONY: clean
clean:
	@rm -rf $(DST)

.PHONY: serve
serve: all
	@cd $(DST); python3 -m http.server
