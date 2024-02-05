package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	BsonID       primitive.ObjectID `bson:"_id"`
	ID *string            `json:"id"`
	Email    *string            `json:"email" validate:"email,required"`
	Password *string            `json:"Password" validate:"required,min=3"`
	Name     *string            `json:"name" validate:"required,min=2,max=100"`

}
