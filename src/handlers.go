package src

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Article struct {
	Id			primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
	Title       string 				`json:"title"`
	Description string 				`json:"description"`
	Content     string 				`json:"content"`
	Username 	string 				`json:"username"`
	Private 	bool   				`json:"private"`
}

type User struct {
	Email string 	`json:"email"`
	Username string `json:"username"`
	Password string	`json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// TODO: This must be an Environment variable
var JwtKey = []byte("donotinvademything")

var DBConnection *mongo.Client
var ArticleCollection *mongo.Collection
var UserCollection *mongo.Collection

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
	UserCollection = DBConnection.Database("personalweb").Collection("users")
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validating non existence of the new user
	filter := bson.D{{"username", user.Username}}
	exist, _ := UserCollection.CountDocuments(context.TODO(), filter)
	if exist > 0 {
		http.Error(w, "An user with that username already exist.", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	_, err = UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request)  {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := bson.D{{"username", user.Username}}
	var storedUser User
	err := UserCollection.FindOne(context.TODO(), filter).Decode(&storedUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expireTime := time.Now().Add(time.Hour)
	claims := Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "authentication",
		Value: tokenString,
		Expires: expireTime,
	})
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var articles []*Article
	cur, err := ArticleCollection.Find(context.TODO(), bson.D{{"private", false }})
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

func GetArticlesByUsername(w http.ResponseWriter, r *http.Request) {
	username := Param(r.Context(), "username")
	filter := bson.D{{"username", username}, { "private", false }}
	var articles []*Article
	cur, err := ArticleCollection.Find(context.TODO(), filter)
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
	username := r.Context().Value("username").(string)

	post, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request wasn't well formated", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(post, &newArticle); err != nil {
		http.Error(w, "Request wasn't well formated", http.StatusBadRequest)
		return
	}

	newArticle.Username = username
	_, err = ArticleCollection.InsertOne(context.TODO(), newArticle)
	if err != nil {
		http.Error(w, "Internal error, We couldn't insert the new post to the db.", http.StatusInternalServerError)
		return
	}
}

func GetArticleById(w http.ResponseWriter, r *http.Request) {
	id := Param(r.Context(), "id")
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id",  objId}}

	var article Article
	err := ArticleCollection.FindOne(context.TODO(), filter).Decode(&article)
	if err != nil {
		http.Error(w, "The item that you are looking for doesn't exist.", 404)
		return
	}

	if err = json.NewEncoder(w).Encode(article); err != nil {
		panic("Error encoding Article to JSON format.")
		return
	}
}

func UpdateArticleById(w http.ResponseWriter, r *http.Request) {
	id := Param(r.Context(), "id")
	objId, _ := primitive.ObjectIDFromHex(id)
	username := r.Context().Value("username").(string)
	filter := bson.D{{"_id",  objId}, { "username", username}}

	var article Article
	post, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request wasn't well formated", 400)
		return
	}

	if err = json.Unmarshal(post, &article); err != nil {
		http.Error(w, "Request wasn't well formated", 400)
		return
	}

	article.Username = username
	doc := ArticleCollection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": article})
	if doc.Err() != nil {
		http.Error(w, "We couldn't update the item, sorry for the inconvenience.", 500)
		return
	}
}

func DeleteArticleById(w http.ResponseWriter, r *http.Request) {
	id := Param(r.Context(), "id")
	objId, _ := primitive.ObjectIDFromHex(id)
	username := r.Context().Value("username").(string)
	filter := bson.D{{"_id",  objId}, { "username", username}}

	_, err := ArticleCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, "We couldn't delete the item", 500)
		return
	}
}


