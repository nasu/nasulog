// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Article struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type NewArticle struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type Tag struct {
	Name string `json:"name"`
}
