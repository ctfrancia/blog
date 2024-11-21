package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/ctfrancia/blog/internal/model"
	"github.com/ctfrancia/blog/pkg/post"
)

type application struct {
	templateCache map[string]*template.Template
	postDir       string
	errorLog      *log.Logger
	infoLog       *log.Logger
	file          post.FileReader
	topicsCache   model.IndexData
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache("ui/templates")
	if err != nil {
		return fmt.Errorf("creating template cache: %w", err)
	}

	topicsCache, err := newPostsCache("internal/posts")
	if err != nil {
		return fmt.Errorf("creating posts cache: %w", err)
	}

	app := &application{
		postDir:       "internal/posts",
		templateCache: templateCache,
		errorLog:      errorLog,
		infoLog:       infoLog,
		file:          post.FileReader{Dir: "internal/posts"},
		topicsCache:   topicsCache,
	}

	srv := &http.Server{
		Addr:    ":3030",
		Handler: app.routes(),
	}

	fmt.Printf("Server is running on port %s", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("starting server error: %w", err)
	}

	return nil
}
