package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	colle "github.com/servicego/src/service/collections"
	conn "github.com/servicego/src/service/connection"
	structs "github.com/servicego/src/service/model"

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
	testaFace()
	fmt.Println("start serviceGO")
	userConnection.Dsn = "mysql:dbname=web;host=127.0.0.1"
	userConnection.User = "root"
	userConnection.Password = ""

	fields = append(fields, structs.Field{Name: "nome", Type: "text", Title: "Nome"})

	router := mux.NewRouter()
	conn.CreateConnection(router)
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

//use mongo ok
func getForm(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	forms, _ = colle.GetFormsData()
	for _, item := range forms {
		if item.FormName == params["formName"] {
			response.Header().Add("content-type", "application/json")
			json.NewEncoder(response).Encode(item)
			return
		}
	}
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(&structs.Form{})
}

//use mongo ok
func getForms(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	forms, err := colle.GetFormsData()
	if err != nil {
		fmt.Println(err.Error())
	}
	json.NewEncoder(response).Encode(forms)
}

//use mongo ok
func deleteForm(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	err := colle.DeleteForm(params["formName"])
	response.Header().Add("content-type", "application/json")
	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`"message" : "` + err.Error() + `"`))
	}
	response.Write([]byte(`"message" : "Sucesfull delete!"`))
}

//use mongo ok
func createForm(response http.ResponseWriter, request *http.Request) {
	var form structs.Form
	_ = json.NewDecoder(request.Body).Decode(&form)
	response.Header().Add("content-type", "application/json")

	forms = append(forms, form)
	form, err := colle.SaveFormData(form)

	if err != nil {
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`"message" : "` + err.Error() + `"`))
	}
	json.NewEncoder(response).Encode(form)
}

func getConn(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	json.NewEncoder(response).Encode(userConnection)
}

func createField(response http.ResponseWriter, request *http.Request) {
	var field structs.Field
	_ = json.NewDecoder(request.Body).Decode(&field)

	if field.Name != "" {
		formName := mux.Vars(request)["formName"]
		for i, formx := range forms {
			if formx.FormName == formName {
				addFieldForm(formx, field, i)
				response.Header().Add("content-type", "application/json")
				json.NewEncoder(response).Encode(forms)
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

//Use mongo ok
func getField(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fields, err := colle.GetFieldsFormData(params["formName"])
	response.Header().Add("content-type", "application/json")

	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	for _, item := range fields {
		if item.Name == params["name"] {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
}

//Use mongo - ok
func getFields(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	fields, err := colle.GetFieldsFormData(params["formName"])
	response.Header().Add("content-type", "application/json")

	if err != nil {
		fmt.Println(err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	} else {
		json.NewEncoder(response).Encode(fields)
	}
}

func deleteField(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	for indexForm, formx := range forms {
		if formx.FormName == params["formName"] {
			for index, item := range formx.Fields {
				if item.Name == params["name"] {
					forms[indexForm].Fields = append(forms[indexForm].Fields[:index], forms[indexForm].Fields[index+1:]...)
					break
				}
			}
			response.Header().Add("content-type", "application/json")
			json.NewEncoder(response).Encode(forms[indexForm])
			return
		}
	}
}
