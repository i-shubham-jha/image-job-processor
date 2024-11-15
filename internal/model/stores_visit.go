package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type VisitInfo struct {
	StoreID    string   `bson:"store_id"`
	VisitTime  string   `bson:"visit_time"`
	ImageURLs  []string `bson:"image_urls"`
	ImageUUIDs []string `bson:"image_uuids"`
	Perimeters []int64  `bson:"perimeters"`
}

type StoresVisit struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Status string             `bson:"status"`
	Error  string             `bson:"error"`
	Count  int                `bson:"count"`
	Visits []VisitInfo        `bson:"visits"`
}
