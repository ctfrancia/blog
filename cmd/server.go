package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ctfrancia/blog/pkg/post"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	mux := http.NewServeMux()
	postReader := post.FileReader{
		Dir: "internal/posts",
	}

	postTemplate := template.Must(template.ParseFiles("ui/templates/post.gohtml"))
	mux.HandleFunc("GET /posts/{slug}", post.PostHandler(postReader, postTemplate))

	indexTemplate := template.Must(template.ParseFiles("ui/templates/index.gohtml"))
	mux.HandleFunc("GET /", post.IndexHandler(postReader, indexTemplate))

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	return nil
}
