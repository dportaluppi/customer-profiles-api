package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const accountIDKey = "accountId"

// MongoRepository is a generic repository for MongoDB.
type MongoRepository[T Entity] struct {
	client     *mongo.Client
	db         string
	collection string
}

// NewMongoRepository creates a new instance of MongoRepository.
func NewMongoRepository[T Entity](client *mongo.Client, db, collection string) Repository[T] {
	return &MongoRepository[T]{
		client:     client,
		db:         db,
		collection: collection,
	}
}

func (r *MongoRepository[T]) Upsert(ctx context.Context, accountId string, entity T) (T, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	var objID primitive.ObjectID
	var err error

	isNew := entity.GetID() == ""
	if isNew {
		objID = primitive.NewObjectID()
		entity.SetID(objID.Hex())
		entity.SetCreatedAt(time.Now())
	} else {
		objID, err = primitive.ObjectIDFromHex(entity.GetID())
		if err != nil {
			return *new(T), err
		}
		entity.SetUpdatedAt(time.Now())
	}

	update := bson.M{"$set": entity}
	opts := options.Update().SetUpsert(true)

	filter := bson.M{"_id": objID, accountIDKey: accountId}
	_, err = coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return *new(T), err
	}

	return entity, nil
}

// GetByID finds an entity by its ID.
func (r *MongoRepository[T]) GetByID(ctx context.Context, accountID, id string) (T, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return *new(T), err
	}

	filter := bson.M{"_id": objID, accountIDKey: accountID}

	var result = *new(T)

	err = coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Delete removes an entity by its ID.
func (r *MongoRepository[T]) Delete(ctx context.Context, accountID, id string) error {
	coll := r.client.Database(r.db).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, accountIDKey: accountID}
	_, err = coll.DeleteOne(ctx, filter)
	return err
}

// GetAll retrieves all entities, with pagination.
func (r *MongoRepository[T]) GetAll(ctx context.Context, accountID string, page, limit int) ([]T, int, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	findOptions := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))

	filter := bson.M{accountIDKey: accountID}
	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []T

	for cursor.Next(ctx) {
		var entity = *new(T)
		if err := cursor.Decode(&entity); err != nil {
			return nil, 0, errors.WithStack(err)
		}
		results = append(results, entity)
	}

	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return results, int(count), nil
}

// ExecuteQuery executes a query and returns a slice of entities with pagination.
func (r *MongoRepository[T]) ExecuteQuery(
	ctx context.Context,
	accountID string,
	query map[string]any,
	currentPage,
	perPage int,
) ([]T, int, error) {

	coll := r.client.Database(r.db).Collection(r.collection)

	mongoQuery := bson.M(query)
	mongoQuery[accountIDKey] = accountID

	findOptions := options.Find().
		SetSkip(int64((currentPage - 1) * perPage)).
		SetLimit(int64(perPage))

	cursor, err := coll.Find(ctx, mongoQuery, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []T

	for cursor.Next(ctx) {
		var entity = *new(T)
		if err := cursor.Decode(&entity); err != nil {
			return nil, 0, err
		}
		results = append(results, entity)
	}

	totalItems, err := coll.CountDocuments(ctx, mongoQuery)
	if err != nil {
		return nil, 0, err
	}

	return results, int(totalItems), nil
}

// ExecutePipeline executes an aggregation pipeline and returns a slice of entities with pagination.
func (r *MongoRepository[T]) ExecutePipeline(
	ctx context.Context,
	accountID string,
	pipeline map[string]any,
	currentPage,
	perPage int,
) ([]T, int, error) {

	coll := r.client.Database(r.db).Collection(r.collection)

	pipeline[accountIDKey] = accountID
	matchStage := bson.D{{"$match", bson.D{{"$expr", pipeline}}}}

	countPipeline := mongo.Pipeline{
		matchStage,
		{{"$count", "total"}},
	}

	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	var countResult struct{ Total int }
	if countCursor.Next(ctx) {
		if err := countCursor.Decode(&countResult); err != nil {
			return nil, 0, err
		}
	}

	pipelineWithPagination := mongo.Pipeline{
		matchStage,
		{{"$skip", int64((currentPage - 1) * perPage)}},
		{{"$limit", int64(perPage)}},
	}

	cursor, err := coll.Aggregate(ctx, pipelineWithPagination)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []T

	for cursor.Next(ctx) {
		var entity = *new(T)
		if err := cursor.Decode(&entity); err != nil {
			return nil, 0, err
		}
		results = append(results, entity)
	}

	return results, countResult.Total, nil
}
