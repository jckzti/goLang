package collections

import (
	"context"
	"time"

	conn "github.com/servicego/src/service/connection"
	structs "github.com/servicego/src/service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//SaveUserData save form in mongoDB
func SaveUserData(form structs.Form) (structs.Form, error) {
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
