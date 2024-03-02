ifdef CI
	BASE_URL=https://thara.dev
else
	BASE_URL=http://localhost:8000
endif

YEAR=$(shell date +%Y)

PANDOC_OPT=-t html -f gfm+yaml_metadata_block --template=template -V base_url=$(BASE_URL) -V year=$(YEAR)

MD_FILES=$(shell ls -d $$(find ./pages -type f) | sed 's/^\.\/pages/public2/g')
HTML_FILES=$(MD_FILES:.md=.html)

md:
	echo $(MD_FILES)
	echo $(HTML_FILES)

all: $(HTML_FILES)

public2/%.html: pages/%.md
	@mkdir -p $(dir $@)
	pandoc -s $< -o $@ $(PANDOC_OPT)

.PHONY: clean
clean:
	@rm -f $(HTML_FILES)

.PHONY: serve
serve: all
	@cd public2; python3 -m http.server
