package main

import (
	"bytes"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	Title   string `toml:"title"`
	Content template.HTML
	Author  Author `toml:"author"`
}

type Author struct {
	Name   string `toml:"name"`
	Email  string `toml:"email"`
	GitHub string `toml:"github"`
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

		var post PostData
		remainingMd, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		err = mdRenderer.Convert([]byte(remainingMd), &buf)
		if err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
		}
		post.Content = template.HTML(buf.String())

		err = tpl.Execute(w, post)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
