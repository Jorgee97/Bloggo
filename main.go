package main

import (
	"github.com/jorgee97/goblogmicro/src"
	"net/http"
)

func main() {
	router := src.NewRouter()

	router.HandleFunc("POST", "/signup", src.SingUp)
	router.HandleFunc("POST", "/login", src.Login)
	router.Handle("GET", "/blog/username/:username", http.HandlerFunc(src.GetArticlesByUsername))
	router.Handle("GET", "/blog/:id", src.JWTAuthentication(http.HandlerFunc(src.GetArticleById)))
	router.Handle("PUT", "/blog/:id", src.JWTAuthentication(http.HandlerFunc(src.UpdateArticleById)))
	router.Handle("DELETE", "/blog/:id", src.JWTAuthentication(http.HandlerFunc(src.DeleteArticleById)))
	router.Handle("POST", "/blog/", src.JWTAuthentication(http.HandlerFunc(src.PostArticle)))
	router.HandleFunc("GET", "/blog/", src.GetAllArticles)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic("Server failed with error: " + err.Error())
	}
}
