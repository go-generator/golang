package models

import "gopkg.in/mgo.v2/bson"

// Represents a movie, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type Merchant struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	MerchantName        string        `bson:"merchantName" json:"merchantName"`
}
