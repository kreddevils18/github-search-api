package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type RepositoryEntity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Topics      []string           `bson:"topics,omitempty" json:"topics,omitempty"`
}
