package src

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

var DBConnection *mongo.Client
var ArticleCollection *mongo.Collection

func init()  {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongodb.")
	DBConnection = client
	ArticleCollection = DBConnection.Database("personalweb").Collection("blog")
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	var articles []*Article
	cur, err := ArticleCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var article Article
		err := cur.Decode(&article)
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, &article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	if err := json.NewEncoder(w).Encode(articles); err != nil {
		panic("Error.")
	}
}

func PostArticle(w http.ResponseWriter, r *http.Request) {
	var newArticle Article
	post, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request wasn't well formated", 400)
	}

	if err := json.Unmarshal(post, &newArticle); err != nil {
		http.Error(w, "Request wasn't well formated", 400)
	}

	_, err = ArticleCollection.InsertOne(context.TODO(), newArticle)
	if err != nil {
		http.Error(w, "Internal error, We couldn't insert the new post to the db.", 500)
	}

	w.WriteHeader(http.StatusOK)
}

func GetArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := Param(r.Context(), "title")
	filter := bson.D{{"title", title}}

	var article Article
	err := ArticleCollection.FindOne(context.TODO(), filter).Decode(&article)
	if err != nil {
		http.Error(w, "The item that you are looking for doesn't exist.", 404)
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(article); err != nil {
		panic("Error encoding Article to JSON format.")
	}
}

func UpdateArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := Param(r.Context(), "title")
	filter := bson.D{{"title", title}}


	var newArticle Article
	post, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request wasn't well formated", 400)
	}

	if err = json.Unmarshal(post, &newArticle); err != nil {
		http.Error(w, "Request wasn't well formated", 400)
	}

	doc := ArticleCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": newArticle})
	if doc.Err() != nil {
		http.Error(w, "We couldn't update the item, sorry for the inconvenience.", 500)
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := Param(r.Context(), "title")
	filter := bson.D{{"title", title}}

	_, err := ArticleCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "We couldn't delete the item", 500)
	}

	w.WriteHeader(http.StatusOK)
}


