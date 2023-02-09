package repository

import (
	"context"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/dbcontext"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"github.com/kien-hoangtrung/github-repository/internal/repository/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepoRepository struct {
	logger     log.ILogger
	collection *mongo.Collection
}

func NewRepoRepository(logger log.ILogger, db *dbcontext.DB) *RepoRepository {
	collection := db.Client.Database("github-repository").Collection("repositories")
	return &RepoRepository{logger, collection}
}

func (r RepoRepository) CreateOne(repo entity.RepositoryEntity) error {
	_, err := r.collection.InsertOne(context.TODO(), repo)

	return err
}

func (r RepoRepository) FindByName(name string) (bson.D, error) {
	filter := bson.D{{"name", name}}

	var result bson.D
	err := r.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r RepoRepository) FindMany(skip int64, limit int64) ([]entity.RepositoryEntity, int64, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit)

	cursor, err := r.collection.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		return nil, -1, err
	}

	var results []entity.RepositoryEntity
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, -1, err
	}

	totalRecords, err := r.collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, -1, err
	}

	return results, totalRecords, nil
}
