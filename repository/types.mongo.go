package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository holds the client and context for MongoDB
type MongoRepository struct {
	Client *mongo.Client
	Ctx    context.Context
}
