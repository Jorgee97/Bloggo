package main

import (
	"github.com/jorgee97/bloggo/src"
	"log"
	"net/http"
	"time"
)

func main() {
	router := src.NewRouter()

	router.HandleFunc("POST", "/signup", src.SingUp)
	router.HandleFunc("POST", "/login", src.Login)
	// Public access, returns all username articles while they are not private
	router.Handle("GET", "/blog/username/:username", http.HandlerFunc(src.GetArticlesByUsername))
	// Private access with JWT, for management of the Articles of the user.
	router.Handle("GET", "/blog/dashboard",
		src.JWTAuthentication(http.HandlerFunc(src.GetArticlesByUsernamePrivate)))
	router.Handle("GET", "/blog/:id", http.HandlerFunc(src.GetArticleById))
	router.Handle("PUT", "/blog/:id", src.JWTAuthentication(http.HandlerFunc(src.UpdateArticleById)))
	router.Handle("DELETE", "/blog/:id", src.JWTAuthentication(http.HandlerFunc(src.DeleteArticleById)))
	router.Handle("POST", "/blog/", src.JWTAuthentication(http.HandlerFunc(src.PostArticle)))
	router.HandleFunc("GET", "/blog/", src.GetAllArticles)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
