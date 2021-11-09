package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName    string = "posts"
	mongoQueryTimeout        = 10 * time.Second
)

type PostData struct {
	Title      string             `bson:"title" json:"title"`
	Body       string             `bson:"body" json:"body"`
	CreateDate time.Time          `bson:"createdate" json:"createdate"`
	IsActive   bool               `bson:"isactive" json:"isactive"`
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	CategoryId int                `bson:"categoryid" json:"categoryid"`
}

// Repository represents an post repository
type PostRepository interface {
	FindAll() ([]*PostData, error)
	FindById(id string) (*PostData, error)
	FindByPagination(page int64, pageSize int64) ([]*PostData, error)
	FindByCategoryPagination(categoryid int, page int, pageSize int) ([]*PostData, error)
	Create(data PostData) (bool, error)
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a mongo API definition repo
func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{collection: db.Collection(collectionName)}, nil
}

func (r *MongoRepository) Add(data PostData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, data)
	if err != nil {
		return false, err
	}
	return true, err
}

// FindByPagination fetches all the API definitions available
func (r *MongoRepository) FindByCategoryPagination(categoryid int64, page int64, pageSize int64) ([]*PostData, error) {
	var result []*PostData

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	options := options.Find()
	options.SetSort(bson.M{"createdate": 1})
	options.SetSkip((page - 1) * pageSize)
	options.SetLimit(pageSize)

	cur, err := r.collection.Find(ctx, bson.M{"categoryid": categoryid, "isactive": true}, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		u := new(PostData)
		if err := cur.Decode(u); err != nil {
			return nil, err
		}

		result = append(result, u)
	}
	return result, cur.Err()
}

// FindByPagination fetches all the API definitions available
func (r *MongoRepository) FindByPagination(page int64, pageSize int64) ([]*PostData, error) {
	var result []*PostData

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	options := options.Find()
	options.SetSort(bson.M{"createdate": 1})
	options.SetSkip((page - 1) * pageSize)
	options.SetLimit(pageSize)

	cur, err := r.collection.Find(ctx, bson.M{"isactive": true}, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		u := new(PostData)
		if err := cur.Decode(u); err != nil {
			return nil, err
		}

		result = append(result, u)
	}

	return result, cur.Err()
}

// FindAll fetches all the API definitions available
func (r *MongoRepository) FindAll() ([]*PostData, error) {
	var result []*PostData

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "createdate", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		u := new(PostData)
		if err := cur.Decode(u); err != nil {
			return nil, err
		}

		result = append(result, u)
	}

	return result, cur.Err()
}

func (r *MongoRepository) FindById(id string) (*PostData, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return r.findOneByQuery(bson.M{"_id": objectId})
}

func (r *MongoRepository) findOneByQuery(query interface{}) (*PostData, error) {
	var result PostData

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	err := r.collection.FindOne(ctx, query).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, err
	}

	return &result, err
}
