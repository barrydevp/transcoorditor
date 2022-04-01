package mongodb

import (
	"context"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/common"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var logger = common.Logger().WithFields(logrus.Fields{
	"pkg": "store/mongodb",
})

type baseRepo struct {
	Db *mongo.Database
}

type mongodbBackend struct {
	*store.Backend
	c  *mongo.Client
	db *mongo.Database
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

func NewStore() (store.Interface, error) {
	client, err := connectMongo()
	if err != nil {
		return nil, err
	}

	logger.Info("Mongodb is connecting.")

	// ping to ensure we were connected
	err = pingMongo(client)
	if err != nil {
		return nil, err
	}

	logger.Info("Mongodb is connected.")

	db := client.Database(viper.GetString("MONGODB_DB"))

	baseRepo := &baseRepo{
		Db: db,
	}
	backend := &store.Backend{
		SessionImpl:     NewSession(baseRepo),
		ParticipantImpl: NewParticipant(baseRepo),
		ReplsetImpl:     NewReplset(baseRepo),
	}

	return &mongodbBackend{
		Backend: backend,
		c:       client,
		db:      db,
	}, nil
}

func (s *mongodbBackend) Close() {
	if s.c == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), MongoDisconnectTimeout*time.Second)
	defer cancel()

	if err := s.c.Disconnect(ctx); err != nil {
		logger.Error("Oops... Cannot disconnect mongodb! Reason: %v", err)

		return
	}

	logger.Info("Mongodb is disconnected.")
}
