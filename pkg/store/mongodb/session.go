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

// memory storage
// TBD
type Session struct {
	*StoreBase
	col *mongo.Collection
}

func NewSession(opts *StoreBase) *Session {

	return &Session{
		StoreBase: opts,
		col:       opts.Db.Collection("sessions"),
	}
}

func (s *Session) Save(session *schema.Session) error {
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

func (s *Session) PutById(id string, schemaUpdate *schema.Session) (*schema.Session, error) {
	update := bson.D{}

	if schemaUpdate.State != "" {
		update = append(update, bson.E{"state", schemaUpdate.State})
	}

	if schemaUpdate.Errors != nil {
		update = append(update, bson.E{"errors", &schemaUpdate.Errors})
	}

	if schemaUpdate.UpdatedAt != nil {
		update = append(update, bson.E{"updatedAt", schemaUpdate.UpdatedAt})
	}

	if schemaUpdate.Timeout != 0 {
		update = append(update, bson.E{"timeout", schemaUpdate.Timeout})
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

func (s *Session) FindById(id string) (*schema.Session, error) {
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

func (s *Session) Find(search *schema.SessionSearch) ([]*schema.Session, error) {
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

func (s *Session) UpdateById(id string, schemaUpdate *schema.SessionUpdate) (*schema.Session, error) {
	update := bson.D{}

	if schemaUpdate.State != nil {
		update = append(update, bson.E{"state", schemaUpdate.State})
	}

	if schemaUpdate.Errors != nil {
		update = append(update, bson.E{"errors", &schemaUpdate.Errors})
	}

	if schemaUpdate.UpdatedAt != nil {
		update = append(update, bson.E{"updatedAt", schemaUpdate.UpdatedAt})
	}

	if schemaUpdate.Timeout != nil {
		update = append(update, bson.E{"timeout", schemaUpdate.Timeout})
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

func (s *Session) DeleteById(id string) (*schema.Session, error) {
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
