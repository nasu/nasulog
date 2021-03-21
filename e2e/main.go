package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nasu/nasulog/api/domain/article"
)

func main() {
	api := flag.String("api", "", "an api URL which should be inspected.")
	flag.Parse()
	art := createArticle(*api)
	log.Println(art)
}

func createArticle(api string) *article.Article {
	query := `
	mutation {
		createArticle(input: {
			title: "e2e-test",
			content: "This was created by e2e-test",
			tags: ["e2e"],
		}) {
			id
			title
			content
			tags
			created_at
			updated_at
		}
	}
	`
	resp, err := post(api, payload(query))
	if err != nil {
		log.Fatalf("createArticle:post: err=%s", err)
	}

	defer resp.Body.Close()
	data := validate(resp, "createArticle")
	log.Println(data)
	if _, ok := data["createArticle"]; !ok {
		log.Fatalf("createArticle:not found createArticle")
	}
	art := data["createArticle"].(map[string]interface{})
	for _, p := range []string{"id", "title", "content", "tags", "created_at", "updated_at"} {
		if _, ok := art[p]; !ok {
			log.Fatalf("createArticle:not found %s", p)
		}
	}
	createdAt, err := time.Parse(time.RFC3339, art["created_at"].(string))
	if err != nil {
		log.Fatalf("createArticle:failed time parse. erro=%s", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, art["updated_at"].(string))
	if err != nil {
		log.Fatalf("createArticle:failed time parse. erro=%s", err)
	}
	return &article.Article{
		ID:        art["id"].(string),
		Title:     art["title"].(string),
		Content:   art["content"].(string),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func validate(resp *http.Response, method string) map[string]interface{} {
	if resp.StatusCode != 200 {
		log.Fatalf("%s: statusCode=%d", method, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s:body:read: err=%s", method, err)
	}
	bodyMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		log.Fatalf("c%s:json:unmarshal: err=%s", method, err)
	}
	if errors, ok := bodyMap["errors"]; ok {
		log.Fatalf("%s:return errors: errors=%s", method, errors)
	}
	data, ok := bodyMap["data"]
	if !ok {
		log.Fatalf("%s:not found data", method)
	}
	return data.(map[string]interface{})
}

func post(api, body string) (*http.Response, error) {
	req, err := http.NewRequest("POST", api, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	return http.DefaultClient.Do(req)
}

func payload(query string) string {
	type v struct{}
	type p struct {
		OperationName *string
		Query         string
		Variables     v
	}
	pp := p{nil, query, v{}}
	j, err := json.Marshal(&pp)
	if err != nil {
		log.Fatalf("payload:json:marshal: err=%s", err)
	}
	//log.Println(string(j))
	return string(j)
}
