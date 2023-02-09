package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type FileEntity struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RepoName string             `bson:"repo_name,omitempty" json:"repo_name,omitempty"`
	FileName string             `bson:"file_name,omitempty" json:"file_name,omitempty"`
	Path     string             `bson:"path,omitempty" json:"path,omitempty"`
	Content  string             `bson:"content,omitempty" json:"content,omitempty"`
}
