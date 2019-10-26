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

//Form usada para guardar os dados para servirem de base para construção dinâmica de forms em tela
type Form struct {
	FormName string  `json:"formName"`
	Fields   []Field `json:"field"`
}

var fields []Field
var userConnection UserConnection
var forms []Form

//Função principal
func main() {

	userConnection.Dsn = "mysql:dbname=web;host=127.0.0.1"
	userConnection.User = "root"
	userConnection.Password = ""

	fields = append(fields, Field{Name: "nome", Type: "text", Title: "Nome"})

	router := mux.NewRouter()
	router.HandleFunc("/conn", GetConn).Methods("GET")
	CreateFieldRequests(router)
	CreateFormRequests(router)

	log.Fatal(http.ListenAndServe(":8000", router))

}

//CreateFieldRequests function
func CreateFieldRequests(router *mux.Router) {
	router.HandleFunc("/fields/{formName}", GetFields).Methods("GET")
	router.HandleFunc("/field/{formName}/{name}", GetField).Methods("GET")
	router.HandleFunc("/field/{formName}", CreateField).Methods("POST")
	router.HandleFunc("/field/{formName}/{name}", DeleteField).Methods("DELETE")
}

//CreateFormRequests function
func CreateFormRequests(router *mux.Router) {
	router.HandleFunc("/forms", GetForms).Methods("GET")
	router.HandleFunc("/form/{formName}", GetForm).Methods("GET")
	router.HandleFunc("/form", CreateForm).Methods("POST")
	router.HandleFunc("/form/{formName}", DeleteForm).Methods("DELETE")
}

//GetForm function
func GetForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range forms {
		if item.FormName == params["formName"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Form{})
}

//GetForms function
func GetForms(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(forms)
}

//DeleteForm function
func DeleteForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range forms {
		if item.FormName == params["formName"] {
			forms = append(forms[:index], forms[index+1:]...)
			json.NewEncoder(w).Encode(forms)
			break
		}
		json.NewEncoder(w).Encode(forms)
	}
}

//CreateForm function
func CreateForm(w http.ResponseWriter, r *http.Request) {
	var form Form
	_ = json.NewDecoder(r.Body).Decode(&form)
	for _, formx := range forms {
		if formx.FormName == form.FormName {
			json.NewEncoder(w).Encode(forms)
			return
		}
	}
	forms = append(forms, form)
	json.NewEncoder(w).Encode(forms)
}

//GetConn function
func GetConn(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userConnection)
}

//CreateField function
func CreateField(w http.ResponseWriter, r *http.Request) {
	var field Field
	_ = json.NewDecoder(r.Body).Decode(&field)

	if field.Name != "" {
		formName := mux.Vars(r)["formName"]
		for i, formx := range forms {
			if formx.FormName == formName {
				AddFieldForm(formx, field, i)
				json.NewEncoder(w).Encode(forms)
				return
			}
		}
	}
}

//AddFieldForm function
func AddFieldForm(form Form, field Field, position int) {
	for i, item := range forms[position].Fields {
		if item.Name == field.Name {
			forms[position].Fields[i] = field
			return
		}
	}
	if field.Type == "" {
		field.Type = "text"
	}

	forms[position].Fields = append(forms[position].Fields, field)
}

//GetField function
func GetField(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, item := range forms {
		if item.FormName == params["formName"] {
			for _, item := range item.Fields {
				if item.Name == params["name"] {
					json.NewEncoder(w).Encode(item)
					return
				}
			}
			return
		}
	}

	json.NewEncoder(w).Encode(&Field{})
}

//GetFields function
func GetFields(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range forms {
		if item.FormName == params["formName"] {
			json.NewEncoder(w).Encode(item.Fields)
			return
		}
	}
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
