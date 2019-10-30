package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	structs "github.com/servicego/src/service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

//Client responsáve pelas ligação com o mongo
var Client *mongo.Client

//Person struct test
type Person struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FistName string             `json:"firstname, omitempty" bson:"firstname, omitempty"`
	LastName string             `json:"lastname, omitempty" bson:"lastname, omitempty"`
}

//GetConn function for manipulation mongo
func GetConn() *mongo.Client {
	if Client == nil {
		CreateConnection(nil)
	}
	return Client
}

func createPersonEndpoint(response http.ResponseWriter, resquest *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(resquest.Body).Decode(&person)
	collection := Client.Database("theveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	person.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`"message" : "` + err.Error() + `"`))
		return
	}
	json.NewEncoder(response).Encode(result)
}

func getPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := Client.Database("theveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

//SaveFormData save form in mongoDB
func SaveFormData(form structs.Form) (structs.Form, error) {
	collection := Client.Database("theveloper").Collection("forms")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	form.ID = primitive.NewObjectID()
	//result, err := collection.InsertOne(ctx, form)
	_, err := collection.InsertOne(ctx, form)
	if err != nil {
		return form, err
	}

	return form, nil
}

//CreateConnection function
func CreateConnection(router *mux.Router) {
	fmt.Println("Testa mongo")
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))

	if err != nil {
		fmt.Println("Erro ao criar Client")
		log.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)
	if err != nil {
		fmt.Println("Erro ao conectar")
		log.Println(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = Client.Ping(ctx, readpref.Primary())

	if err != nil {
		fmt.Println("Erro ao pingar")
		log.Println(err)
	}

	_ = Client

	if router != nil {
		router.HandleFunc("/person", createPersonEndpoint).Methods("POST")
		router.HandleFunc("/people", getPeopleEndpoint).Methods("GET")
	}
	//start mongo
	//$ mongod --dbpath=/Users/jonathancani/data/db

}
