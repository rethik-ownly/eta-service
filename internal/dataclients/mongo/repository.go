package mongo

import (
	"context"

	"github.com/nutanalabs/eta-service/internal/config"
	"github.com/nutanalabs/rapido-mongo-go/mongo"
	"github.com/nutanalabs/rapido-mongo-go/mongo/results"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	FindOne(collectionName string, filter interface{}, opts *options.FindOneOptions) *mongo2.SingleResult
	FindMany(collectionName string, filter interface{}, queryResponse *results.QueryResponse, opts *options.FindOptions) error
	InsertOne(collectionName string, document interface{}) (*mongo2.InsertOneResult, error)
	UpdateOne(collectionName string, filter interface{}, update interface{}, opts *options.UpdateOptions) (*mongo2.UpdateResult, error)
	CountDocuments(collectionName string, filter interface{}, opts *options.CountOptions) (int64, error)
	DeleteOne(collectionName string, filter interface{}, opts *options.DeleteOptions) (*mongo2.DeleteResult, error)
	FindOneAndDelete(collectionName string, filter interface{}, opts *options.FindOneAndDeleteOptions) *mongo2.SingleResult
	AggregatePipeline(collectionName string, pipeline interface{}, queryResponse *results.QueryResponse, opts *options.AggregateOptions) error
}

type repositoryImpl struct {
	config   *config.Config
	database *mongo.Database
}

func NewMongoRepository(config *config.Config, mongoClient Client) Repository {
	return &repositoryImpl{
		config:   config,
		database: mongoClient.GetDatabase(),
	}
}

func (r *repositoryImpl) FindOne(collectionName string, filter interface{}, opts *options.FindOneOptions) *mongo2.SingleResult {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.FindOne(ctx, filter, opts)
}

func (r *repositoryImpl) FindMany(collectionName string, filter interface{}, queryResponse *results.QueryResponse, opts *options.FindOptions) error {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.Find(ctx, filter, queryResponse, opts)
}

func (r *repositoryImpl) InsertOne(collectionName string, document interface{}) (*mongo2.InsertOneResult, error) {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.InsertOne(ctx, document)
}

func (r *repositoryImpl) UpdateOne(collectionName string, filter interface{}, update interface{}, opts *options.UpdateOptions) (*mongo2.UpdateResult, error) {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.UpdateOne(ctx, filter, update, opts)
}

func (r *repositoryImpl) CountDocuments(collectionName string, filter interface{}, opts *options.CountOptions) (int64, error) {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.CountDocuments(ctx, filter, opts)
}

func (r *repositoryImpl) DeleteOne(collectionName string, filter interface{}, opts *options.DeleteOptions) (*mongo2.DeleteResult, error) {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.DeleteOne(ctx, filter, opts)
}

func (r *repositoryImpl) FindOneAndDelete(collectionName string, filter interface{}, opts *options.FindOneAndDeleteOptions) *mongo2.SingleResult {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.FindOneAndDelete(ctx, filter, opts)
}

func (r *repositoryImpl) AggregatePipeline(collectionName string, pipeline interface{}, queryResponse *results.QueryResponse, opts *options.AggregateOptions) error {
	collection := r.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), r.config.GetMongoReadTimeout())
	defer cancel()
	return collection.Aggregate(ctx, pipeline, queryResponse, opts)
}
