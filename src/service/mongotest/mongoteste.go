package mongoteste

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

type Person struct {
	ID       primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	FistName string             `json:"firstname, omitempty" bson:"firstname, omitempty"`
	LastName string             `json:"lastname, omitempty" bson:"lastname, omitempty"`
}

func CreatePersonEndpoint(response http.ResponseWriter, resquest *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(resquest.Body).Decode(&person)
	collection := client.Database("theveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func TestaMongo(router *mux.Router) {
	fmt.Println("Testa mongo")
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))

	if err != nil {
		fmt.Println("Erro ao criar client")
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Erro ao conectar")
		log.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		fmt.Println("Erro ao pingar")
		log.Fatal(err)
	}

	_ = client
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")

	//start mongo
	//$ mongod --dbpath=/Users/jonathancani/data/db

}
