package profile

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/yalochat/go-commerce-components/flat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	db         string
	collection string
}

func NewMongoRepository(client *mongo.Client, db string) *MongoRepository {
	return &MongoRepository{
		client:     client,
		db:         db,
		collection: "attributes",
	}
}

func (r *MongoRepository) Get(ctx context.Context) (Attribute, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return Attribute{}, err
	}
	defer cursor.Close(ctx)

	var attributes []Attribute
	for cursor.Next(ctx) {
		var attribute Attribute
		if err := cursor.Decode(&attribute); err != nil {
			return Attribute{}, err
		}
		attributes = append(attributes, attribute)
	}

	if len(attributes) == 0 {
		return Attribute{}, nil
	}

	return attributes[0], nil
}

func (r *MongoRepository) Updater(ctx context.Context, profile *Profile) (Attribute, error) {
	current, err := r.Get(ctx)
	if err != nil {
		return Attribute{}, err
	}

	f := flat.NewFlattener()
	coll := r.client.Database(r.db).Collection(r.collection)

	opts := options.Update().SetUpsert(true)

	b, err := json.Marshal(profile)
	if err != nil {
		return Attribute{}, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return Attribute{}, err
	}

	fp := f.Flatten(m)

	if current.Attributes == nil {
		current.Attributes = make(Attributes)
	}

	for k, v := range fp {
		strValue := fmt.Sprint(v)

		if _, ok := current.Attributes[k]; !ok {
			current.Attributes[k] = Counter{strValue: 0}
		}

		current.Attributes[k][strValue]++
	}

	update := bson.M{"$set": current}

	id := current.ID
	if current.ID.IsZero() {
		id = primitive.NewObjectID()
	}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": id}, update, opts)
	if err != nil {
		return Attribute{}, err
	}

	return current, nil
}

func (r *MongoRepository) Delete(ctx context.Context, profile *Profile) (Attribute, error) {
	current, err := r.Get(ctx)
	if err != nil {
		return Attribute{}, err
	}

	f := flat.NewFlattener()
	coll := r.client.Database(r.db).Collection(r.collection)

	opts := options.Update().SetUpsert(true)

	var m map[string]interface{}
	if err := mapstructure.Decode(profile, &m); err != nil {
		return Attribute{}, err
	}

	fp := f.Flatten(m)

	for k, v := range fp {
		attr, ok := current.Attributes[k]
		if !ok {
			continue
		}

		strValue := fmt.Sprint(v)

		c, ok := attr[strValue]
		if !ok {
			continue
		}

		if c > 1 {
			current.Attributes[k][strValue]--
		} else {
			delete(current.Attributes[k], strValue)
		}
	}

	update := bson.M{"$set": current}

	_, err = coll.UpdateOne(ctx, bson.M{"_id": current.ID}, update, opts)
	if err != nil {
		return Attribute{}, err
	}

	return current, nil
}
