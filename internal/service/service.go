package service

import (
	"context"
	"retail_pulse/internal/db"
	"retail_pulse/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const db_name string = "retail_pulse"
const collection_name string = "stores_visits"

type StoresVisitService struct {
	client *mongo.Client
}

// NewStoresVisitService creates a new instance of StoresVisitService
func NewStoresVisitService() *StoresVisitService {
	return &StoresVisitService{
		client: db.GetMongoClient(),
	}
}

// InsertStoresVisit inserts a new StoresVisit into the database and returns its ID
func (svs *StoresVisitService) InsertStoresVisit(storesVisit model.StoresVisit) (primitive.ObjectID, error) {
	collection := svs.client.Database(db_name).Collection(collection_name)

	// Insert the document and get the result
	result, err := collection.InsertOne(context.TODO(), storesVisit)
	if err != nil {
		return primitive.ObjectID{}, err // Return an empty ObjectID and the error
	}

	// Return the inserted ID
	return result.InsertedID.(primitive.ObjectID), nil
}

// FindStoresVisitByID fetches a StoresVisit by its ID from the database
func (svs *StoresVisitService) FindStoresVisitByID(id primitive.ObjectID) (*model.StoresVisit, error) {
	collection := svs.client.Database(db_name).Collection(collection_name)

	var storesVisit model.StoresVisit
	filter := bson.M{"_id": id}

	// Find the document by ID
	err := collection.FindOne(context.TODO(), filter).Decode(&storesVisit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}

	return &storesVisit, nil
}

// GetStatusAndErrorByID fetches the status and error message for a StoresVisit by its ID
func (svs *StoresVisitService) GetStatusAndErrorByID(id primitive.ObjectID) (string, string, string, error) {
	collection := svs.client.Database(db_name).Collection(collection_name)

	// Create a variable to hold the result
	result := struct {
		Status        string `bson:"status"`
		Error         string `bson:"error"`
		FailedStoreID string `bson:"failed_store_id"`
	}{}

	filter := bson.M{"_id": id}
	projection := bson.M{"status": 1, "error": 1, "failed_store_id": 1} // fields to include

	// Find the document by ID with projection
	err := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", "", mongo.ErrNoDocuments // No document found
		}
		return "", "", "", err // Return the error if something else went wrong
	}

	// Return the status and error message
	return result.Status, result.Error, result.FailedStoreID, nil
}
