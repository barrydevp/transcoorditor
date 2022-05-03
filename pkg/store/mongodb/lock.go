package mongodb

import (
	"context"

	"github.com/barrydevp/transcoorditor/pkg/schema"
	"github.com/barrydevp/transcoorditor/pkg/store"
	"github.com/barrydevp/transcoorditor/pkg/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type lockTableRepo struct {
	*baseRepo
	col *mongo.Collection
}

func NewLockTable(opts *baseRepo) *lockTableRepo {

	return &lockTableRepo{
		baseRepo: opts,
		col:      opts.Db.Collection("locktable"),
	}
}

func (s *lockTableRepo) Save(lockEnt *schema.LockEntry) error {
	if _, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"key", lockEnt.Key}}
		doc := &schema.LockEntry{}
		err := s.col.FindOne(ctx, filter).Decode(doc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				doc = nil
			} else {
				return nil, err
			}
		}

		update := bson.D{{"key", lockEnt.Key}, {"owner", lockEnt.Owner}, {"expiredAt", lockEnt.ExpiredAt}}
		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

		doc = &schema.LockEntry{}
		err = s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(doc)

		if err != nil {
			return nil, err
		}

		return doc, nil
	}, 10); err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) Update(lockEnt *schema.LockEntry) error {
	_, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		doc := &schema.LockEntry{}
		filter := bson.D{{"key", lockEnt.Key}, {"owner", lockEnt.Owner}}
		err := s.col.FindOne(ctx, filter).Decode(doc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				doc = nil
				return nil, store.ErrLockNotFound
			} else {
				return nil, err
			}
		}

		doc = &schema.LockEntry{}
		update := bson.D{{"expiredAt", lockEnt.ExpiredAt}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err = s.col.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts).Decode(doc)

		if err != nil {
			return nil, err
		}

		return doc, nil
	}, 10)

	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) Find(key string) (*schema.LockEntry, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		doc := &schema.LockEntry{}
		filter := bson.D{{"key", key}}

		err := s.col.FindOne(ctx, filter).Decode(doc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return doc, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.LockEntry)

	return r, nil
}

func (s *lockTableRepo) FindWithOwner(key string, owner string) (*schema.LockEntry, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		doc := &schema.LockEntry{}
		filter := bson.D{{"key", key}, {"owner", owner}}

		err := s.col.FindOne(ctx, filter).Decode(doc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return doc, nil
	}, 10)

	if err != nil {
		return nil, err
	}

	r, _ := doc.(*schema.LockEntry)

	return r, nil
}

func (s *lockTableRepo) Delete(lockEnt *schema.LockEntry) error {
	_, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		doc := &schema.LockEntry{}
		filter := bson.D{{"key", lockEnt.Key}, {"owner", lockEnt.Owner}}

		err := s.col.FindOneAndDelete(ctx, filter).Decode(doc)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, nil
			}

			return nil, err
		}

		return doc, nil
	}, 10)

	if err != nil {
		return err
	}

	return nil
}

func (s *lockTableRepo) DeleteByOwner(owner string) (int64, error) {
	doc, err := util.WithTimeout(func(ctx context.Context) (interface{}, error) {
		filter := bson.D{{"owner", owner}}

		result, err := s.col.DeleteMany(ctx, filter)

		if err != nil {
			return nil, err
		}

		return result.DeletedCount, nil
	}, 10)

	if err != nil {
		return 0, err
	}

	if err != nil {
		return -1, err
	}

	r, _ := doc.(int64)

	return r, nil
}
