package main

import (
	"github.com/jorgee97/goblogmicro/src"
	"net/http"
)

func main() {

	router := src.NewRouter()

	router.HandleFunc("GET", "/blog/:title", src.GetArticleByTitle)
	router.HandleFunc("PUT", "/blog/:title", src.UpdateArticleByTitle)
	router.HandleFunc("DELETE", "/blog/:title", src.DeleteArticleByTitle)
	router.HandleFunc("GET", "/blog/", src.GetAllArticles)
	router.HandleFunc("POST", "/blog/", src.PostArticle)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic("Server failed with error: " + err.Error())
	}
}
