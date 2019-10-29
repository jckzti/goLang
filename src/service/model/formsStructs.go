package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

//Field usada para guardar os dados para servirem de base para construção dinâmica dos campos em tela
type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Title string `json:"title"`
}

//Form usada para guardar os dados para servirem de base para construção dinâmica de forms em tela
type Form struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FormName string             `json:"formName"`
	Fields   []Field            `json:"fields"`
}
