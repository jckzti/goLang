package mongoteste

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	ID       primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	FistName string             `json:"firstname, omitempty" bson:"firstname, omitempty"`
	LastName string             `json:"lastname, omitempty" bson:"lastname, omitempty"`
}

func TestaMongo() {
	fmt.Println("Testa mongo")
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = client.Connect(ctx)
}
