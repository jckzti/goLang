package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Kagami/go-face"
	"github.com/gorilla/mux"
	colle "github.com/servicego/src/service/collections"
	conn "github.com/servicego/src/service/connection"
	structs "github.com/servicego/src/service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
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
var rec *face.Recognizer

const dataDir = "images"

type Face struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PersonName string             `json:"personName"`
	Descriptor face.Descriptor    `json:"descriptor"`
}

var faceMD Face

//Função principal
func main() {
	var err error
	rec, err = face.NewRecognizer(dataDir)
	if err != nil {
		fmt.Println(err)
	}

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
	setupRoutes()
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

func identificaPessoa(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File\n")
		fmt.Println(err)
		return
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(w, "Successfully Uploaded File\n")
	fmt.Fprintf(w, "Successfully Uploaded File\n")

	faces, err := rec.Recognize(fileBytes)
	if err != nil || len(faces) < 1 {
		fmt.Println("Não foi possível identificar uma face na imagem enviada: %v\n", err)
		fmt.Fprintf(w, "Não foi possível identificar uma face na imagem enviada: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Println("Rostos encontrados na imagem: ====wwww %v\n", len(faces))
	fmt.Fprintf(w, "Rostos encontrados na imagem: %v\n", len(faces))

	//Preenche com todos os descripts de todos os rostos da base de dados
	rostos, err := getFacesData()

	var samples []face.Descriptor
	var pessoas []int32

	for i, f := range rostos {
		//saveImageDescriptor(f.Descriptor)
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		pessoas = append(pessoas, int32(i))
	}

	rec.SetSamples(samples, pessoas)

	if err != nil {
		fmt.Fprintf(w, "Não achou rosto na imagem: %v", err)
		return
	}

	fmt.Fprint(w, "Pessoa identificadas:\n")
	for i, facex := range faces {
		//personID := rec.ClassifyThreshold(facex.Descriptor, 0.1)
		personID := rec.Classify(facex.Descriptor)
		if personID < 0 {
			fmt.Fprintf(w, "Pessoa não encontrada ou não cadastrada no sistema\n")
			continue
		}

		fmt.Println(personID)
		fmt.Println(rostos[personID].PersonName)
		fmt.Fprintf(w, "[%v] - Pessoa encontrada: %v\n", i+1, rostos[personID].PersonName)

	}
	//saveImageDescriptor(faces[0].Descriptor, rostos[avengerID].PersonName) // caso tenha encontrado, cataloga imagem para fazer parte do acervo
}

func enviaPessoa(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(10 << 20)
	nome := r.FormValue("nome")
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	//tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")

	if err != nil {
		fmt.Println(err)
		return
	}
	//defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	//tempFile.Write(fileBytes)

	faces, err := rec.Recognize(fileBytes)
	if err != nil {
		fmt.Fprintf(w, "Não foi possível identificar uma face na imagem enviada: %v\n", err)
		return
	}
	fmt.Fprintf(w, "Rostos encontrados na imagem: %v\n", len(faces))

	saveImageDescriptor(faces[0].Descriptor, nome)

}

func saveImageDescriptor(descriptor face.Descriptor, nome string) error {
	collection := conn.GetConn().Database("recog").Collection("faces")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	faceMD.ID = primitive.NewObjectID()
	faceMD.PersonName = nome
	faceMD.Descriptor = descriptor
	//result, err := collection.InsertOne(ctx, form)
	_, err := collection.InsertOne(ctx, faceMD)
	if err != nil {
		return err
	}

	return nil
}

func getFacesData() ([]Face, error) {
	fmt.Println("")
	var rostos []Face

	collection := conn.GetConn().Database("recog").Collection("faces")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rosto Face
		cursor.Decode(&rosto)
		rostos = append(rostos, rosto)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rostos, nil
}

func setupRoutes() {
	http.HandleFunc("/identifica", identificaPessoa)
	http.HandleFunc("/envia", enviaPessoa)
	http.ListenAndServe(":8080", nil)
}
