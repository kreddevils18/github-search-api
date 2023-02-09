package repository

import (
	"context"
	"github.com/kien-hoangtrung/github-repository/internal/file/entity"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/dbcontext"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	logger     log.ILogger
	collection *mongo.Collection
}

func NewFileRepository(logger log.ILogger, db *dbcontext.DB) *FileRepository {
	collection := db.Client.Database("github-repository").Collection("files")

	return &FileRepository{logger, collection}
}

func (f FileRepository) CreateOne(file entity.FileEntity) error {
	_, err := f.collection.InsertOne(context.TODO(), file)

	return err
}

func (f FileRepository) FindByPath(path string) (bson.D, error) {
	filter := bson.D{{"path", path}}

	var result bson.D
	err := f.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f FileRepository) List() []entity.FileEntity {
	filter := bson.D{{}}
	cursor, err := f.collection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []entity.FileEntity
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	return results
}
