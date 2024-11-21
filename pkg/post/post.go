package post

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
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
