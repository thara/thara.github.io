package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v2"
)

type PageType int

const (
	_ PageType = iota
	pageTypePost
)

type Page struct {
	pageType PageType

	Title string
	Site  SiteConfig

	Path string

	PageTitle string
	Content   string

	DatePageCreated string

	distPath string
	layout   string

	timestamp int64
}

func loadPage(filename, parent, distDir string, page *Page) error {
	p := path.Join(parent, filename)

	src, err := files.Open(p)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v ", p, err)
	}
	defer src.Close()

	b, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("fail to read %s: %v ", p, err)
	}
	content := string(b)

	pageTitle := "thara.dev"
	distFilename := filename
	layout := "base"

	ext := path.Ext(filename)
	switch ext {
	case ".md", ".markdown":
		bs := bytes.SplitN(b, []byte("---"), 3)

		f := make(map[string]interface{})
		if err := yaml.Unmarshal(bs[1], &f); err != nil {
			return fmt.Errorf("fail to parse frontmatter of %s: %v ", filename, err)
		}
		body := bs[2]

		if v, ok := f["title"]; ok {
			pageTitle = v.(string)
		}

		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(body), &buf); err != nil {
			return fmt.Errorf("fail to parse markdown of %s: %v ", filename, err)
		}

		content = string(buf.Bytes())

		if v, ok := f["path"]; ok {
			distFilename = v.(string)
		} else {
			distFilename = strings.TrimSuffix(filename, ext)
		}

		if filepath.Base(distDir) == "posts" {
			layout = "post"
			date, name, ts, err := parsePageName(distFilename)
			if err != nil {
				return fmt.Errorf("fail to parse page name %s: %v", distFilename, err)
			}
			distFilename = name
			page.DatePageCreated = date

			page.pageType = pageTypePost
			page.timestamp = ts
		}
	}

	distPath := path.Join(distDir, distFilename)

	page.layout = layout
	page.PageTitle = pageTitle
	page.Content = content
	page.distPath = distPath

	page.Site = siteCfg
	page.Title = fmt.Sprintf("%s | %s", pageTitle, siteCfg.SiteTitle)

	page.Path = strings.TrimPrefix(distPath, distDirRoot)

	return nil
}

func (p *Page) write(t *template.Template) error {
	dist, err := os.Create(p.distPath)
	if err != nil {
		return fmt.Errorf("fail to create %s: %v ", p.distPath, err)
	}

	if err := t.Execute(dist, p); err != nil {
		return fmt.Errorf("template execution error: %v", err)
	}
	return nil
}

func parsePageName(filename string) (date, title string, ts int64, err error) {
	d := filename[:10]
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		return "", "", 0, fmt.Errorf("fail to parse time %s: %v", d, err)
	}
	return d, filename[11:], t.UnixMicro(), nil
}

func getContent[T any](p T, t *template.Template) (string, error) {
	var b bytes.Buffer
	if err := t.Execute(&b, p); err != nil {
		return "", fmt.Errorf("template execution error: %v", err)
	}
	return string(b.Bytes()), nil
}
