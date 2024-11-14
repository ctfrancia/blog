package main

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	fr := FileReader{}
	postTemplate := template.Must(template.ParseFiles("cmd/templates/post.gohtml"))
	mux.HandleFunc("GET /posts/{slug}", PostHandler(fr, postTemplate))

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

type SlugReader interface {
	Read(slug string) (string, error)
}
type PostData struct {
	Title   string
	Content template.HTML
	Author  string
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

func PostHandler(sl SlugReader, tpl *template.Template) http.HandlerFunc {
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		pp := fmt.Sprintf("internal/posts/%s.md", slug)
		postMarkdown, err := sl.Read(pp)
		if err != nil {
			log.Println(err)
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		err = mdRenderer.Convert([]byte(postMarkdown), &buf)
		if err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
		}

		err = tpl.Execute(w, PostData{
			Content: template.HTML(buf.String()),
			Author:  "Jon Calhoun",
			Title:   "My Blog",
		})
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
