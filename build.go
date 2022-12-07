package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

//go:embed pages/*
var files embed.FS

func buildPages(distDirRoot string) error {
	if err := os.RemoveAll(distDirRoot); err != nil {
		return fmt.Errorf("fail to remove dir %s: %v", distDirRoot, err)
	}
	if err := os.Mkdir(distDirRoot, os.ModePerm); err != nil {
		return fmt.Errorf("fail to mkdir %s: %v", distDirRoot, err)
	}

	var pages []Page
	err := walkDir("pages", func(filename, parent string) error {
		distDir := path.Join(distDirRoot, strings.TrimPrefix(parent, "pages"))
		if err := os.MkdirAll(distDir, os.ModePerm); err != nil {
			return fmt.Errorf("fail to mkdir %s: %v", distDir, err)
		}
		var page Page
		if err := loadPage(filename, parent, distDir, &page); err != nil {
			return fmt.Errorf("fail to load page %s/%s: %v ", parent, filename, err)
		}
		pages = append(pages, page)
		return nil
	})
	if err != nil {
		return fmt.Errorf("fail to walk pages dir: %v", err)
	}

	var posts []Page
	for _, page := range pages {
		var layout *template.Template
		if t, ok := templates[page.layout]; ok {
			layout = t
		} else {
			layout = defaultTemplate
			log.Printf("template not found for page layout: %s", page.layout)
		}

		if err := page.write(layout); err != nil {
			return fmt.Errorf("fail to write page %s: %v", page.distPath, err)
		}

		if page.pageType == pageTypePost {
			posts = append(posts, page)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].timestamp > posts[j].timestamp
	})

	var indexPage IndexPage
	if err := loadPage("index.md", "pages", distDirRoot, &indexPage.Page); err != nil {
		return fmt.Errorf("fail to load index page: %v ", err)
	}
	indexPage.Posts = posts
	content, err := getContent(indexPage, templates["index"])
	if err != nil {
		return fmt.Errorf("fail to fill index page: %v", err)
	}
	indexPage.Content = content
	if err := indexPage.write(templates["base"]); err != nil {
		return fmt.Errorf("fail to write index page: %v", err)
	}

	return err
}

func walkDir(dir string, f func(filename, parent string) error) error {
	es, err := files.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("fail to read dir %s: %v", dir, err)
	}
	for _, e := range es {
		if e.IsDir() {
			childDir := path.Join(dir, e.Name())
			if err := walkDir(childDir, f); err != nil {
				return fmt.Errorf("fail to walk dir %s: %v", childDir, err)
			}
		} else {
			if err := f(e.Name(), dir); err != nil {
				return fmt.Errorf("fail to call walk func %s/%s: %v", dir, e.Name(), err)
			}
		}
	}
	return nil
}
