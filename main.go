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

var baseURL string
var serve bool
var port int

func init() {
	flag.StringVar(&baseURL, "base", "", "base URL")
	flag.BoolVar(&serve, "serve", false, "serves the site locally")
	flag.IntVar(&port, "port", 8080, "listen port")
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

		b := fmt.Sprintf(":%d", port)
		fmt.Printf("Listen at %s\n", b)
		log.Fatal(http.ListenAndServe(b, handlers.CompressHandler(r)))
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
