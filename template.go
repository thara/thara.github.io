package main

import (
	"log"
	"net/url"
	"text/template"
	"time"
)

const baseURL string = "http://localhost:8080"

var templates map[string]*template.Template
var defaultTemplate *template.Template

func init() {
	templates = make(map[string]*template.Template)

	t, err := template.ParseFiles("templates/base.html", "templates/_head.html", "templates/_footer.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	templates["base"] = t
	defaultTemplate = t

	t, err = template.ParseFiles("templates/post.html", "templates/_head.html", "templates/_footer.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	templates["post"] = t

	t, err = template.ParseFiles("templates/_index.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	templates["index"] = t
}

var siteCfg SiteConfig = func() SiteConfig {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("url parse error: %v", err)
	}
	return SiteConfig{
		SiteTitle: "thara.dev",
		Author:    "thara",
		BaseURL:   u,
		Year:      time.Now().Year(),
	}
}()

type SiteConfig struct {
	SiteTitle string
	Author    string
	BaseURL   *url.URL
	Year      int
}

type IndexPage struct {
	Page
	Posts []Page
}
