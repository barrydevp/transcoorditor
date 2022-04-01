package mongodb

import (
	"context"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"github.com/hashicorp/raft"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type replsetRepo struct {
	*baseRepo
	col *mongo.Collection
}

func NewReplset(opts *baseRepo) *replsetRepo {

	return &replsetRepo{
		baseRepo: opts,
		col:      opts.Db.Collection("replset"),
	}
}

func (s *replsetRepo) SaveLastLog(log *raft.Log) error {
	_, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{}
		opts := options.FindOneAndReplace().SetReturnDocument(options.After).SetUpsert(true)

		result := s.col.FindOneAndReplace(ctx, filter, log, opts)

		if result.Err() != nil {
			return nil, result.Err()
		}

		return nil, nil
	}, 10)

	if err != nil {
		return err
	}

	return nil
}

func (s *replsetRepo) GetLastLog() (*raft.Log, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{}
		opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(1)

		cursor, err := s.col.Find(ctx, filter, opts)

		if err != nil {
			return nil, err
		}

		var results []*raft.Log
		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return nil, mongo.ErrNoDocuments
		}

		return results[0], nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*raft.Log)

	return r, nil
}
