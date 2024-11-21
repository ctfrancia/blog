package post

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
	"path/filepath"
	"strings"
	"time"
)

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

type PostMetaData struct {
	Slug        string
	Title       string    `toml:"title"`
	Author      Author    `toml:"author"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
}

type MetaDataQuerier interface {
	Query() ([]PostMetaData, error)
}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct {
	// Dir is the directory where the markdown files are stored
	Dir string
}

func (fr FileReader) Read(slug string) (string, error) {
	slugPath := filepath.Join(fr.Dir, slug+".md")
	f, err := os.Open(slugPath)
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

func (fr FileReader) Query() ([]PostMetaData, error) {
	postPath := filepath.Join(fr.Dir, "*.md")
	files, err := filepath.Glob(postPath)
	if err != nil {
		return nil, fmt.Errorf("error globbing files: %w", err)
	}

	var posts []PostMetaData
	for _, f := range files {
		f, err := os.Open(f)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		defer f.Close()

		var post PostMetaData
		_, err = frontmatter.Parse(f, &post)
		if err != nil {
			return nil, fmt.Errorf("error parsing frontmatter: %w", err)
		}

		// Extract the slug from the filename
		post.Slug = strings.TrimSuffix(filepath.Base(f.Name()), ".md")

		posts = append(posts, post)
	}

	return posts, nil
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
		postMarkdown, err := sl.Read(slug)
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

type IndexData struct {
	Posts []PostMetaData
}

func IndexHandler(mq MetaDataQuerier, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadata, err := mq.Query()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error querying metadata", http.StatusInternalServerError)
			return
		}

		data := IndexData{
			Posts: metadata,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
