package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	structs "github.com/servicego/src/service/model"
	mongoteste "github.com/servicego/src/service/mongotest"

	"github.com/gorilla/mux"
)

//UserConnection usada para guardar os dados da conexão para uso posterior no PHP
// type UserConnection struct {
// 	Dsn      string `json:"dsn"`
// 	User     string `json:"user"`
// 	Password string `json:"password"`
// }

var fields []structs.Field
var userConnection structs.UserConnection
var forms []structs.Form

//Função principal
func main() {

	userConnection.Dsn = "mysql:dbname=web;host=127.0.0.1"
	userConnection.User = "root"
	userConnection.Password = ""

	fields = append(fields, structs.Field{Name: "nome", Type: "text", Title: "Nome"})

	fmt.Println("antes do teste do mongo")
	mongoteste.TestaMongo()
	fmt.Println("depois do teste do mongo")
	router := mux.NewRouter()
	router.HandleFunc("/conn", getConn).Methods("GET")
	createFieldRequests(router)
	createFormRequests(router)

	log.Fatal(http.ListenAndServe(":8000", router))

}

//CreateFieldRequests function
func createFieldRequests(router *mux.Router) {
	router.HandleFunc("/fields/{formName}", getFields).Methods("GET")
	router.HandleFunc("/field/{formName}/{name}", getField).Methods("GET")
	router.HandleFunc("/field/{formName}", createField).Methods("POST")
	router.HandleFunc("/field/{formName}/{name}", deleteField).Methods("DELETE")
}

//CreateFormRequests function
func createFormRequests(router *mux.Router) {
	router.HandleFunc("/forms", getForms).Methods("GET")
	router.HandleFunc("/form/{formName}", getForm).Methods("GET")
	router.HandleFunc("/form", createForm).Methods("POST")
	router.HandleFunc("/form/{formName}", deleteForm).Methods("DELETE")
}

//GetForm function
func getForm(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range forms {
		if item.FormName == params["formName"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&structs.Form{})
}

func getForms(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(forms)
}

func deleteForm(w http.ResponseWriter, r *http.Request) {
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

func createForm(w http.ResponseWriter, r *http.Request) {
	var form structs.Form
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

func getConn(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userConnection)
}

func createField(w http.ResponseWriter, r *http.Request) {
	var field structs.Field
	_ = json.NewDecoder(r.Body).Decode(&field)

	if field.Name != "" {
		formName := mux.Vars(r)["formName"]
		for i, formx := range forms {
			if formx.FormName == formName {
				addFieldForm(formx, field, i)
				json.NewEncoder(w).Encode(forms)
				return
			}
		}
	}
}

func addFieldForm(form structs.Form, field structs.Field, position int) {
	for i, item := range forms[position].Fields {
		if item.Name == field.Name {
			forms[position].Fields[i] = field
			return
		}
	}
	if field.Type == "" {
		field.Type = "text"
	}

	if field.Title == "" {
		field.Title = field.Name
	}

	forms[position].Fields = append(forms[position].Fields, field)
}

func getField(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(&structs.Field{})
}

func getFields(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range forms {
		if item.FormName == params["formName"] {
			json.NewEncoder(w).Encode(item.Fields)
			return
		}
	}
}

func deleteField(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for indexForm, formx := range forms {
		if formx.FormName == params["formName"] {
			for index, item := range formx.Fields {
				if item.Name == params["name"] {
					forms[indexForm].Fields = append(forms[indexForm].Fields[:index], forms[indexForm].Fields[index+1:]...)
					break
				}
			}
			json.NewEncoder(w).Encode(forms[indexForm])
			return
		}
	}
}
