package mongodb

import (
	"Basic/Trainning4/redis/staff/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	staffCollection *mongo.Collection
	teamCollection  *mongo.Collection
}

func (DB *MongoDB) NewDB() {
	var staffCollection *mongo.Collection = config.GetCollection(config.ConnectDB(), "staff")
	var teamCollection *mongo.Collection = config.GetCollection(config.ConnectDB(), "team")
	DB.staffCollection = staffCollection
	DB.teamCollection = teamCollection
}

func (DB *MongoDB) GetName() string {
	return "MongoDB"
}
