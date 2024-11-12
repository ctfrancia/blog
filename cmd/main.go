package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	fr := FileReader{}
	mux.HandleFunc("GET /posts/{slug}", PostHandler(fr))

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func PostHandler(sl SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		file := fmt.Sprintf("internal/posts/%s.md", slug)
		content, err := sl.Read(file)
		if err != nil {
			fmt.Println(err)
			// TODO: log error according to the context
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Post slug: %s", content)
	}
}
