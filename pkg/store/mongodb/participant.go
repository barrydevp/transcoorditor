package mongodb

import (
	"context"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	r, _ := doc.(*schema.Participant)

	return r, nil
}

func (s *Participant) FindBySessionId(sessionId string) ([]*schema.Participant, error) {
	var results []*schema.Participant

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}}

		cursor, err := s.col.Find(ctx, filter)

		if err != nil {
			return nil, err
		}

		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}

		return results, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.([]*schema.Participant)

	return r, nil
}

func (s *Participant) FindDupInSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
	dupPart := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}, {"clientId", part.ClientId}, {"requestId", part.RequestId}}

		// duplicate detection by requestId
		// if part.RequestId != "" {
		// 	filter = append(filter, bson.E{"requestId", part.RequestId})
		// }

		err := s.col.FindOne(ctx, filter).Decode(dupPart)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return dupPart, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Participant)

	return r, nil
}

func (s *Participant) UpdateById(id string, partUpdate *schema.ParticipantUpdate) (*schema.Participant, error) {
	update := bson.D{}

	if partUpdate.State != nil {
		update = append(update, bson.E{"state", partUpdate.State})
	}

	if partUpdate.UpdatedAt != nil {
		update = append(update, bson.E{"updatedAt", partUpdate.UpdatedAt})
	}

	if partUpdate.CompensateAction != nil {
		update = append(update, bson.E{"compensateAction", partUpdate.CompensateAction})
	}

	if partUpdate.CompleteAction != nil {
		update = append(update, bson.E{"completeAction", partUpdate.CompleteAction})
	}

	// no changes
	if len(update) == 0 {
		return s.FindById(id)
	}

	if partUpdate.UpdatedAt == nil {
		update = append(update, bson.E{"updatedAt", time.Now()})
	}

	part := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(part)

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

	r, _ := doc.(*schema.Participant)

	return r, nil
}
