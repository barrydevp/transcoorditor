package mongodb

import (
	"context"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// memory storage
// TBD
type Participant struct {
	*StoreOptions
	col *mongo.Collection
}

func NewParticipant(opts *StoreOptions) *Participant {

	return &Participant{
		StoreOptions: opts,
		col:          opts.Db.Collection("participants"),
	}
}

func (s *Participant) Save(part *schema.Participant) error {
	if _, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		inserted, err := s.col.InsertOne(ctx, part)

		if err != nil {
			return nil, err
		}

		return inserted, nil
	}, 10); err != nil {
		return err
	}

	return nil
}

func (s *Participant) FindById(id string) (*schema.Participant, error) {
	part := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}

		err := s.col.FindOne(ctx, filter).Decode(part)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return part, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	if doc == nil {
		return nil, err
	}

	return part, nil
}

