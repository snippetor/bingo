package mvc

import (
	"gopkg.in/mgo.v2/bson"
)

type BaseMongoModel struct {
	Id bson.ObjectId `bson:"_id"`
}
