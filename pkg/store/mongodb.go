package store

import (
	"context"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store/mongodb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	c  *mongo.Client
	db *mongo.Database
	l  *logrus.Entry

	session     *mongodb.Session
	participant *mongodb.Participant
}

const (
	MongoConnectTimeout    = 10
	MongoPingTimeout       = 2
	MongoDisconnectTimeout = 3
)

func connectMongo() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoConnectTimeout*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGODB_URI")))

	return client, err
}

func pingMongo(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoPingTimeout*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())

	return err
}

func NewMongoDBStore() (Interface, error) {
	l := common.Logger().WithFields(logrus.Fields{
		"pkg": "store/mongodb",
	})

	client, err := connectMongo()
	if err != nil {
		return nil, err
	}

	l.Info("Mongodb is connecting.")

	// ping to ensure we were connected
	err = pingMongo(client)
	if err != nil {
		return nil, err
	}

	l.Info("Mongodb is connected.")

	db := client.Database(viper.GetString("MONGODB_DB"))

	storeOpts := &mongodb.StoreOptions{
		Db: db,
		L:  l,
	}

	return &MongoDB{
		c:           client,
		db:          db,
		l:           l,
		session:     mongodb.NewSession(storeOpts),
		participant: mongodb.NewParticipant(storeOpts),
	}, nil
}

func (s *MongoDB) Close() {
	if s.c == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), MongoDisconnectTimeout*time.Second)
	defer cancel()

	if err := s.c.Disconnect(ctx); err != nil {
		s.l.Error("Oops... Cannot disconnect mongodb! Reason: %v", err)

		return
	}

	s.l.Info("Mongodb is disconnected.")
}

func (s *MongoDB) Session() Session {
	return s.session
}

func (s *MongoDB) Participant() Participant {
	return s.participant
}
