package user

import (
	"context"
	"github.com/dportaluppi/customer-profiles-api/pkg/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	client     *mongo.Client
	db         string
	collection string
}

func NewMongoRepository(client *mongo.Client, db string) user.Repository {
	return &mongoRepository{
		client:     client,
		db:         db,
		collection: "users",
	}
}
func (r *mongoRepository) Upsert(ctx context.Context, user *user.User) (*user.User, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	var objID primitive.ObjectID
	var err error

	isNew := user.ID == ""
	if isNew {
		objID = primitive.NewObjectID()
		user.CreatedAt = now()
	} else {
		objID, err = primitive.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
		user.UpdatedAt = now()
	}

	update := bson.M{"$set": bson.M{
		"channels":   user.Channels,
		"attributes": user.Attributes,
		"createdAt":  user.CreatedAt,
		"updatedAt":  user.UpdatedAt,
	}}

	opts := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, update, opts)
	if err != nil {
		return nil, err
	}

	if isNew {
		user.ID = objID.Hex()
	}

	return user, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, userID string) (*user.User, error) {
	coll := r.client.Database(r.db).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	user := &user.User{}
	err = coll.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *mongoRepository) Delete(ctx context.Context, userID string) error {
	coll := r.client.Database(r.db).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	_, err = coll.DeleteOne(ctx, filter)
	return err
}

func (r *mongoRepository) GetAll(ctx context.Context, page, limit int) ([]*user.User, int, error) {
	coll := r.client.Database(r.db).Collection(r.collection)
	findOptions := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	cursor, err := coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*user.User
	for cursor.Next(ctx) {
		var user user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, int(count), nil
}

func (r *mongoRepository) ExecuteQuery(ctx context.Context, query map[string]any, currentPage, perPage int) ([]*user.User, int, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	// parse query to bson.M and set pagination
	mongoQuery := bson.M(query)
	findOptions := options.Find().
		SetSkip(int64((currentPage - 1) * perPage)).
		SetLimit(int64(perPage))

	cursor, err := coll.Find(ctx, mongoQuery, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []*user.User
	for cursor.Next(ctx) {
		var user user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		results = append(results, &user)
	}

	totalItems, err := coll.CountDocuments(ctx, mongoQuery)
	if err != nil {
		return nil, 0, err
	}

	return results, int(totalItems), nil
}

func (r *mongoRepository) ExecutePipeline(ctx context.Context, pipeline map[string]any, currentPage, perPage int) ([]*user.User, int, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	// Convertir la consulta map a bson.D para la etapa de match
	matchStage := bson.D{{"$match", bson.D{{"$expr", pipeline}}}}

	// Pipeline para contar los documentos que coinciden con el filtro
	countPipeline := mongo.Pipeline{
		matchStage,
		{{"$count", "total"}},
	}

	countCursor, err := coll.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer countCursor.Close(ctx)

	// Estructura para almacenar el resultado del conteo
	var countResult struct{ Total int }
	if countCursor.Next(ctx) {
		if err := countCursor.Decode(&countResult); err != nil {
			return nil, 0, err
		}
	}

	// Agregar etapas de paginación al pipeline original
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

	var results []*user.User
	for cursor.Next(ctx) {
		var user user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		results = append(results, &user)
	}

	return results, countResult.Total, nil
}
