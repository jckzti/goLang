package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//UserConnection usada para guardar os dados da conexão para uso posterior no PHP
type UserConnection struct {
	Dsn      string `json:"dsn"`
	User     string `json:"user"`
	Password string `json:"password"`
}

//Field usada para guardar os dados para servirem de base para construção dinâmica dos campos em tela
type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Title string `json:"title"`
}

var fields []Field
var userConnection UserConnection

// função principal
func main() {

	userConnection.Dsn = "mysql:dbname=web;host=127.0.0.1"
	userConnection.User = "root"
	userConnection.Password = ""

	fields = append(fields, Field{Name: "nome", Type: "text", Title: "Nome"})

	router := mux.NewRouter()
	router.HandleFunc("/conn", GetConn).Methods("GET")
	router.HandleFunc("/fields", GetFields).Methods("GET")
	router.HandleFunc("/field/{name}", GetField).Methods("GET")
	router.HandleFunc("/field", CreateField).Methods("POST")
	router.HandleFunc("/field/{name}", DeleteField).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}

//GetConn function
func GetConn(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userConnection)
}

//CreateField function
func CreateField(w http.ResponseWriter, r *http.Request) {
	var field Field
	_ = json.NewDecoder(r.Body).Decode(&field)

	for i, item := range fields {
		if item.Name == field.Name {
			fields[i] = field
			//json.NewEncoder(w).Encode("{ 'Response':'atualizado'")
			//json.NewEncoder(w).Encode(fields[i])
			json.NewEncoder(w).Encode(fields)
			return
		}
	}
	if field.Type == "" {
		field.Type = "text"
	}
	fields = append(fields, field)
	json.NewEncoder(w).Encode(fields)
}

//GetField function
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

//GetFields function
func GetFields(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(fields)
}

//DeleteField function
func DeleteField(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range fields {
		if item.Name == params["name"] {
			fields = append(fields[:index], fields[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(fields)
	}
}
