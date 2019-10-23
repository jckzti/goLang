package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// "Person type" (tipo um objeto)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type UserConnection struct {
	Dsn      string `json:"dsn"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Title string `json:"title"`
}

var people []Person
var fields []Field
var userConnection UserConnection

// função principal
func main() {

	userConnection.Dsn = "mysql:dbname=web;host=127.0.0.1"
	userConnection.User = "root"
	userConnection.Password = ""

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})

	router := mux.NewRouter()
	router.HandleFunc("/contato", GetPeople).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/contato/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/contato/{id}", DeletePerson).Methods("DELETE")

	router.HandleFunc("/conn", GetConn).Methods("GET")
	router.HandleFunc("/fields", GetFields).Methods("GET")
	router.HandleFunc("/field/{name}", GetField).Methods("GET")
	router.HandleFunc("/field", CreateField).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))

}

func GetConn(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userConnection)
}

func CreateField(w http.ResponseWriter, r *http.Request) {
	var field Field
	_ = json.NewDecoder(r.Body).Decode(&field)

	for i, item := range fields {
		if item.Name == field.Name {
			fields[i].Type = field.Type
			fields[i].Value = field.Value
			//json.NewEncoder(w).Encode("{ 'Response':'atualizado'")
			//json.NewEncoder(w).Encode(fields[i])
			json.NewEncoder(w).Encode(fields)
			return
		}
	}
	fields = append(fields, field)
	json.NewEncoder(w).Encode(fields)
}

func GetField(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range fields {
		if item.Name == params["name"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Field{})
}

func GetFields(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(fields)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
