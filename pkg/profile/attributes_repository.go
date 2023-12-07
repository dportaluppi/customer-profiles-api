package profile

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/yalochat/go-commerce-components/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAttributeUnreachable = errors.New("attribute name unreachable")
	ErrValuesUnreachable    = errors.New("attribute values unreachable")
	ErrValueConvert         = errors.New("cannot convert attribute value into string")
)

type Flattener interface {
	Flatten(in map[string]any) map[string]any
}

type MongoRepository struct {
	client     *mongo.Client
	db         string
	collection string
	flattener  Flattener
}

func NewMongoRepository(ctx context.Context, client *mongo.Client, db string, f Flattener) *MongoRepository {
	r := &MongoRepository{
		client:     client,
		db:         db,
		collection: "attributes",
		flattener:  f,
	}

	if err := createIndexes(ctx, r); err != nil {
		logging.FromContext(ctx).Fatalf("could not create the mongodb index. %s", err)
	}

	return r
}

func createIndexes(ctx context.Context, r *MongoRepository) error {
	account := mongo.IndexModel{
		Keys: bson.M{"attribute": 1},
	}

	combinedQuery := mongo.IndexModel{
		Keys: bson.D{
			{Key: "attribute", Value: 1},
			{Key: "value", Value: 1},
		},
	}

	var indexes []mongo.IndexModel
	indexes = append(indexes, account, combinedQuery)

	_, err := r.client.Database(r.db).Collection(r.collection).Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *MongoRepository) GetAll(ctx context.Context) (Attributes, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	pipeline := bson.A{
		bson.D{
			{"$group", bson.D{
				{"_id", "$attribute"},
				{"values", bson.D{
					{"$push", "$value"},
				}},
			}},
		},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	result := make(Attributes)

	for _, res := range results {
		attr, ok := res["_id"].(string)
		if !ok {
			return nil, ErrAttributeUnreachable
		}

		values, ok := res["values"].(primitive.A)
		if !ok {
			return nil, ErrValuesUnreachable
		}

		var stringValues []string
		for _, v := range values {
			s, ok := v.(string)
			if !ok {
				return nil, ErrValueConvert
			}

			stringValues = append(stringValues, s)
		}

		result[attr] = stringValues
	}

	return result, nil
}

func (r *MongoRepository) Updater(ctx context.Context, profile *Profile) error {
	coll := r.client.Database(r.db).Collection(r.collection)

	fp, err := r.flatProfile(profile)
	if err != nil {
		return err
	}

	var updates []mongo.WriteModel

	for k, v := range fp {
		filter := bson.D{
			{"attribute", k},
			{"value", v},
		}

		update := mongo.NewUpdateOneModel()
		update.SetFilter(filter)
		update.SetUpdate(bson.D{
			{"$set", bson.D{
				{"attribute", k},
				{"value", v},
			}},
			{"$inc", bson.D{
				{"count", 1},
			}},
		})
		update.SetUpsert(true)

		updates = append(updates, update)
	}

	_, err = coll.BulkWrite(ctx, updates)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) Delete(ctx context.Context, profile *Profile) error {
	coll := r.client.Database(r.db).Collection(r.collection)

	fp, err := r.flatProfile(profile)
	if err != nil {
		return err
	}

	var updates []mongo.WriteModel

	for k, v := range fp {
		filter := bson.D{
			{"attribute", k},
			{"value", v},
		}

		update := mongo.NewUpdateOneModel()
		update.SetFilter(filter)
		update.SetUpdate(bson.D{
			{"$inc", bson.D{
				{"count", -1},
			}},
		})

		updates = append(updates, update)
	}

	_, err = coll.BulkWrite(ctx, updates)
	if err != nil {
		return err
	}

	_, err = coll.DeleteMany(ctx, bson.D{{"count", 0}})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) flatProfile(profile *Profile) (map[string]any, error) {
	b, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	delete(m, "profileId")
	delete(m, "createdAt")
	delete(m, "updatedAt")

	fp := r.flattener.Flatten(m)

	return fp, nil
}
