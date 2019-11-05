package collections

import (
	"context"
	"fmt"
	"time"

	conn "github.com/servicego/src/service/connection"
	structs "github.com/servicego/src/service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

//SaveFormData save form in mongoDB
func SaveFormData(form structs.Form) (structs.Form, error) {
	collection := conn.GetConn().Database("theveloper").Collection("forms")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	form.ID = primitive.NewObjectID()
	//result, err := collection.InsertOne(ctx, form)
	_, err := collection.InsertOne(ctx, form)
	if err != nil {
		return form, err
	}

	return form, nil
}

//GetFieldsFormData function
func GetFieldsFormData(formName string) ([]structs.Field, error) {
	collection := conn.GetConn().Database("theveloper").Collection("forms")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var form structs.Form
	filter := bson.M{"formname": formName}
	fmt.Println("Formname:" + formName)

	err := collection.FindOne(ctx, filter).Decode(&form)

	if err != nil {
		return nil, err
	}
	return form.Fields, nil
}

//GetFormsData function
func GetFormsData() ([]structs.Form, error) {
	fmt.Println("")
	var forms []structs.Form

	collection := conn.GetConn().Database("theveloper").Collection("forms")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var form structs.Form
		cursor.Decode(&form)
		forms = append(forms, form)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return forms, nil
}

func DeleteForm(formName string) error {
	collection := conn.GetConn().Database("theveloper").Collection("forms")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"formname": formName}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
