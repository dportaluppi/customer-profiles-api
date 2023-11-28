package profile

import (
	"context"
	"github.com/dportaluppi/customer-profiles-api/pkg/profile"
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

func NewMongoRepository(client *mongo.Client, db string) profile.Repository {
	return &mongoRepository{
		client:     client,
		db:         db,
		collection: "profiles",
	}
}
func (r *mongoRepository) Updater(ctx context.Context, profile *profile.Profile) (*profile.Profile, error) {
	coll := r.client.Database(r.db).Collection(r.collection)

	var objID primitive.ObjectID
	var err error

	isNew := profile.ID == ""
	if isNew {
		objID = primitive.NewObjectID()
		profile.CreatedAt = now()
	} else {
		objID, err = primitive.ObjectIDFromHex(profile.ID)
		if err != nil {
			return nil, err
		}
		profile.UpdatedAt = now()
	}

	update := bson.M{"$set": bson.M{
		"channels":   profile.Channels,
		"attributes": profile.Attributes,
		"createdAt":  profile.CreatedAt,
		"updatedAt":  profile.UpdatedAt,
	}}

	opts := options.Update().SetUpsert(true)

	_, err = coll.UpdateOne(ctx, bson.M{"_id": objID}, update, opts)
	if err != nil {
		return nil, err
	}

	if isNew {
		profile.ID = objID.Hex()
	}

	return profile, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, profileID string) (*profile.Profile, error) {
	coll := r.client.Database(r.db).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	profile := &profile.Profile{}
	err = coll.FindOne(ctx, filter).Decode(profile)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *mongoRepository) Delete(ctx context.Context, profileID string) error {
	coll := r.client.Database(r.db).Collection(r.collection)
	objID, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	_, err = coll.DeleteOne(ctx, filter)
	return err
}

func (r *mongoRepository) GetAll(ctx context.Context, page, limit int) ([]*profile.Profile, int, error) {
	coll := r.client.Database(r.db).Collection(r.collection)
	findOptions := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	cursor, err := coll.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var profiles []*profile.Profile
	for cursor.Next(ctx) {
		var profile profile.Profile
		if err := cursor.Decode(&profile); err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, &profile)
	}

	count, err := coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return profiles, int(count), nil
}
