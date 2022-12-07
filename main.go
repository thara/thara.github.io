package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/otiai10/copy"
)

const distDirRoot string = "public"

var serve bool
var baseURL string

func init() {
	flag.BoolVar(&serve, "serve", false, "serves the site locally.")
	flag.StringVar(&baseURL, "base", "http://localhost:8080", "base URL")
}

func main() {
	flag.Parse()

	if err := buildPages(distDirRoot); err != nil {
		log.Fatalf("fail to buildPages: %v", err)
	}

	es, err := os.ReadDir("assets")
	if err != nil {
		log.Fatalf("fail to read asset dir: %v", err)
	}
	for _, e := range es {
		dir := e.Name()
		if err := copy.Copy("assets/"+dir, path.Join(distDirRoot, dir)); err != nil {
			log.Fatalf("fail to copy assets: %v", err)
		}
	}

	if serve {
		r := http.NewServeMux()
		r.Handle("/", handlers.LoggingHandler(os.Stdout, http.StripPrefix("/", http.FileServer(&restrictedFileSystem{http.Dir("public")}))))

		fmt.Println("Listen at: localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", handlers.CompressHandler(r)))
	}
}

type restrictedFileSystem struct {
	fs http.FileSystem
}

func (r *restrictedFileSystem) Open(path string) (http.File, error) {
	f, err := r.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		idx := filepath.Join(path, "index.html")
		if _, err := r.fs.Open(idx); err != nil {
			if ex := f.Close(); ex != nil {
				return nil, ex
			}
			return nil, err
		}
	}

	return r.fs.Open(path)
}
