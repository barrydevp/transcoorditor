package mongodb

import (
	"context"
	"time"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/util"

	// "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sessionRepo struct {
	*baseRepo
	col *mongo.Collection
}

func NewSession(opts *baseRepo) *sessionRepo {

	return &sessionRepo{
		baseRepo: opts,
		col:      opts.Db.Collection("sessions"),
	}
}

func (s *sessionRepo) Save(session *schema.Session) error {
	if _, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		inserted, err := s.col.InsertOne(ctx, session)

		if err != nil {
			return nil, err
		}

		return inserted, nil
	}, 10); err != nil {
		return err
	}

	return nil
}

func (s *sessionRepo) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	update := bson.D{}

	if schemaUpdate.State != "" {
		update = append(update, bson.E{"state", schemaUpdate.State})
	}

	if schemaUpdate.Errors != nil {
		update = append(update, bson.E{"errors", &schemaUpdate.Errors})
	}

	if schemaUpdate.EndAt != nil {
		update = append(update, bson.E{"endAt", schemaUpdate.EndAt})
	}

	if schemaUpdate.UpdatedAt != nil {
		update = append(update, bson.E{"updatedAt", schemaUpdate.UpdatedAt})
	}

	if schemaUpdate.Timeout != 0 {
		update = append(update, bson.E{"timeout", schemaUpdate.Timeout})
	}

	if schemaUpdate.Retries != 0 {
		update = append(update, bson.E{"retries", schemaUpdate.Retries})
	}

	// no changes
	if len(update) == 0 {
		return s.FindById(id)
	}

	if schemaUpdate.UpdatedAt == nil {
		update = append(update, bson.E{"updatedAt", time.Now()})
	}

	session := &schema.Session{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(session)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return session, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Session)

	return r, nil
}

func (s *sessionRepo) FindById(id string) (*schema.Session, error) {
	session := &schema.Session{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}

		err := s.col.FindOne(ctx, filter).Decode(session)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return session, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Session)

	return r, nil
}

func (s *sessionRepo) Find(search *schema.SessionSearch) ([]*schema.Session, error) {
	var results []*schema.Session

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{}

		cursor, err := s.col.Find(ctx, filter)

		if err != nil {
			return nil, err
		}

		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}

		return results, nil
	}, 30)

	if err != nil {
		return nil, err
	}

	r, _ := doc.([]*schema.Session)

	return r, nil
}

func (s *sessionRepo) FindAllUnfinished() ([]*schema.Session, error) {
	var results []*schema.Session

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{
			"state", bson.D{{
				"$in", schema.UnfinishedSessionStates(),
			}},
		}}

		cursor, err := s.col.Find(ctx, filter)

		if err != nil {
			return nil, err
		}

		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}

		return results, nil
	}, 30)

	if err != nil {
		return nil, err
	}

	r, _ := doc.([]*schema.Session)

	return r, nil
}

func (s *sessionRepo) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (*schema.Session, error) {
	update := bson.D{}

	if schemaUpdate.State != nil {
		update = append(update, bson.E{"state", schemaUpdate.State})
	}

	if schemaUpdate.Errors != nil {
		update = append(update, bson.E{"errors", &schemaUpdate.Errors})
	}

	if schemaUpdate.EndAt != nil {
		update = append(update, bson.E{"endAt", schemaUpdate.EndAt})
	}

	if schemaUpdate.UpdatedAt != nil {
		update = append(update, bson.E{"updatedAt", schemaUpdate.UpdatedAt})
	}

	if schemaUpdate.Timeout != nil {
		update = append(update, bson.E{"timeout", schemaUpdate.Timeout})
	}

	if schemaUpdate.Retries != nil {
		update = append(update, bson.E{"retries", schemaUpdate.Retries})
	}

	if schemaUpdate.TerminateReason != nil {
		update = append(update, bson.E{"terminateReason", schemaUpdate.TerminateReason})
	}

	// no changes
	if len(update) == 0 {
		return s.FindById(id)
	}

	if schemaUpdate.UpdatedAt == nil {
		update = append(update, bson.E{"updatedAt", time.Now()})
	}

	session := &schema.Session{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(session)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return session, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Session)

	return r, nil
}

func (s *sessionRepo) DeleteById(id string) (*schema.Session, error) {
	session := &schema.Session{}

	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"id", id}}

		err := s.col.FindOneAndDelete(ctx, filter).Decode(session)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return session, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.Session)

	return r, nil
}
