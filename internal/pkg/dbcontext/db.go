package dbcontext

import (
	"context"
	"fmt"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
}

func NewDB(dsn string) (*DB, error) {
	params := "?maxPoolSize=20&w=majority"
	dsn = dsn + params
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	return &DB{client}, nil
}

func NewMongo(conf *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
	)
	return NewDB(dsn)
}
