package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/handlers"
	"github.com/otiai10/copy"
)

const distDirRoot string = "public"

var baseURL string
var serve bool
var port int
var watch bool

func init() {
	flag.StringVar(&baseURL, "base", "", "base URL")
	flag.BoolVar(&serve, "serve", false, "serves the site locally")
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.BoolVar(&watch, "watch", false, "rebuild sites on changed pages dir")
}

func main() {
	flag.Parse()

	buildAll := func() {
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
	}
	buildAll()

	if serve {
		if watch {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatalf("fail to init fsnotify watcher: %v", err)
			}
			defer watcher.Close()

			go func() {
				fmt.Println("watching pages dir...")
				for {
					select {
					case ev, ok := <-watcher.Events:
						if !ok {
							return
						}
						log.Println("event:", ev)
						//NOTE Edit a file on vim, it creates swap files and rename to the original.
						// `chmod` is the last operation in such process.
						if ev.Has(fsnotify.Chmod) {
							fmt.Println("rebuild pages")
							buildAll()
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						log.Println("error:", err)
					}
				}
			}()

			watcher.Add("pages")
			fmt.Println("watch: pages/")
			if err := walkDir("pages", func(dirname, parent string) error {
				p := path.Join(parent, dirname)
				fmt.Printf("watch: %s/\n", p)
				return watcher.Add(p)
			}); err != nil {
				log.Fatalf("fail to watch pages dir: %v", err)
			}
		}

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
