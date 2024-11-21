package main

import (
	"bytes"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/ctfrancia/blog/pkg/post"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"net/http"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.renderIndex(w, r, app.topicsCache)
}

func (app *application) post(w http.ResponseWriter, r *http.Request) {
	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
			),
		),
	)

	slug := r.PathValue("slug")
	postMarkdown, err := app.file.Read(slug)
	if err != nil {
		fmt.Println(err)
		app.errorLog.Println(err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var post post.PostData
	remainingMd, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
	if err != nil {
		http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = mdRenderer.Convert([]byte(remainingMd), &buf)
	if err != nil {
		panic(err)
	}
	post.Content = template.HTML(buf.String())

	app.render(w, r, "post.gohtml", post)
}
