package mongodb

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreOptions struct {
	Db *mongo.Database
	L  *logrus.Entry
}
