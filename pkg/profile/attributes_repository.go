package profile

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func NewMongoRepository(client *mongo.Client, db string, f Flattener) *MongoRepository {
	return &MongoRepository{
		client:     client,
		db:         db,
		collection: "attributes",
		flattener:  f,
	}
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]Attribute, error) {
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

	return attributes, nil
}

func (r *MongoRepository) Updater(ctx context.Context, profile *Profile) error {
	coll := r.client.Database(r.db).Collection(r.collection)

	b, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	fp := r.flattener.Flatten(m)

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

	b, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	fp := r.flattener.Flatten(m)

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
