package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/ctfrancia/blog/internal/model"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// humanDate takes a UTC time and returns a formatted string for better human legibility
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.gohtml"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func newPostsCache(p string) (model.IndexData, error) {
	postsPath := filepath.Join(p, "*.md")
	filenames, err := filepath.Glob(postsPath)
	if err != nil {
		return model.IndexData{}, fmt.Errorf("querying for files: %w", err)
	}

	var posts []model.PostMetaData
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return model.IndexData{}, fmt.Errorf("opening file %s: %w", filename, err)
		}
		defer f.Close()
		var post model.PostMetaData
		_, err = frontmatter.Parse(f, &post)
		if err != nil {
			// return nil, fmt.Errorf("parsing frontmatter for file %s: %w", filename, err)
		}
		post.Slug = strings.TrimSuffix(filepath.Base(filename), ".md")
		posts = append(posts, post)
	}

	return model.IndexData{Posts: posts}, nil
}
