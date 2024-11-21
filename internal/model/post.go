package model

import (
	"time"
)

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
