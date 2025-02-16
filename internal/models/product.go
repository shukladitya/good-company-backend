package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Document struct {
    ID    primitive.ObjectID `bson:"_id" json:"_id"`
    Data  map[string]interface{} `bson:",inline" json:"data"`
}