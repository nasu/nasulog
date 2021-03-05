package main

import (
	"flag"
	"net/http"
)

func main() {
	api := flag.String("api", "", "an api URL which should be inspected.")
	flag.Parse()
	article(*api)
}

func article(api string) {
	body := `{"title":"test2", "content":"test"}`
	http.Post(api+"/article", "application/json", &body)

}
