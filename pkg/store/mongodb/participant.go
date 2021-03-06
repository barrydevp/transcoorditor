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
type participantRepo struct {
	*baseRepo
	col *mongo.Collection
}

func NewParticipant(opts *baseRepo) *participantRepo {

	return &participantRepo{
		baseRepo: opts,
		col:      opts.Db.Collection("participants"),
	}
}

func normalizeParticipant(p *schema.Participant) {
	if p.CompensateAction != nil {
		p.CompensateAction.Data = TryConvertBsonDToM(p.CompensateAction.Data)
	}

	if p.CompleteAction != nil {
		p.CompleteAction.Data = TryConvertBsonDToM(p.CompleteAction.Data)
	}
}

func (s *participantRepo) Save(part *schema.Participant) error {
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

func (s *participantRepo) PutBySessionAndId(sessionId string, id int64, partUpdate *schema.Participant) (*schema.Participant, error) {
	update := bson.D{}

	if partUpdate.State != "" {
		update = append(update, bson.E{"state", partUpdate.State})
	}

	if partUpdate.ClientId != "" {
		update = append(update, bson.E{"clientId", partUpdate.ClientId})
	}

	if partUpdate.RequestId != "" {
		update = append(update, bson.E{"requestId", partUpdate.RequestId})
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
		return s.FindBySessionAndId(sessionId, id)
	}

	if partUpdate.UpdatedAt == nil {
		update = append(update, bson.E{"updatedAt", time.Now()})
	}

	part := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}, {"id", id}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(part)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		normalizeParticipant(part)
		return part, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Participant)

	return r, nil
}

func (s *participantRepo) FindBySessionAndId(sessionId string, id int64) (*schema.Participant, error) {
	part := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}, {"id", id}}

		err := s.col.FindOne(ctx, filter).Decode(part)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		normalizeParticipant(part)
		return part, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Participant)

	return r, nil
}

func (s *participantRepo) FindBySessionId(sessionId string) ([]*schema.Participant, error) {
	var results []*schema.Participant

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}}

		cursor, err := s.col.Find(ctx, filter)

		if err != nil {
			return nil, err
		}

		// @TODO: convert bson.D of compensateAction.data to bson.M?
		// if err := cursor.All(ctx, &results); err != nil {
		// 	return nil, err
		// }

		for cursor.Next(ctx) {
			part := &schema.Participant{}
			if err := cursor.Decode(&part); err != nil {
				return nil, err
			}

			normalizeParticipant(part)
			results = append(results, part)
		}

		return results, nil
	}, 30)

	if err != nil {
		return nil, err
	}

	r, _ := doc.([]*schema.Participant)

	return r, nil
}

func (s *participantRepo) FindDupInSession(sessionId string, part *schema.Participant) (*schema.Participant, error) {
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

func (s *participantRepo) UpdateBySessionAndId(sessionId string, id int64, partUpdate *schema.ParticipantUpdate) (*schema.Participant, error) {
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
		return s.FindBySessionAndId(sessionId, id)
	}

	if partUpdate.UpdatedAt == nil {
		update = append(update, bson.E{"updatedAt", time.Now()})
	}

	part := &schema.Participant{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}, {"id", id}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(part)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		normalizeParticipant(part)

		return part, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Participant)

	return r, nil
}

func (s *participantRepo) CountBySessionId(sessionId string) (int64, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}}

		count, err := s.col.CountDocuments(ctx, filter)

		if err != nil {
			return nil, err
		}

		return count, nil
	}, 10)

	if err != nil {
		return -1, err
	}

	r, _ := doc.(int64)

	return r, nil
}

func (s *participantRepo) DeleteBySessionId(sessionId string) (int64, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"sessionId", sessionId}}

		result, err := s.col.DeleteMany(ctx, filter)

		if err != nil {
			return nil, err
		}

		return result.DeletedCount, nil
	}, 10)

	if err != nil {
		return -1, err
	}

	r, _ := doc.(int64)

	return r, nil
}
