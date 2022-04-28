package mongodb

import (
	"context"
	"github.com/proactiongo/pagocore"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DAO is an interface for DAOs
type DAO interface {
	FetchByID(id string, target interface{}, opts ...*options.FindOneOptions) error
	FetchByIDs(ids []string, target interface{}, opts ...*options.FindOptions) error
	FetchByExIDs(ids []string, target interface{}, opts ...*options.FindOptions) error

	FetchOne(target interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	FetchAll(target interface{}, opts ...*options.FindOptions) error
	FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error

	InsertOne(data interface{}, opts ...*options.InsertOneOptions) (id string, err error)
	InsertMany(rows []interface{}, opts ...*options.InsertManyOptions) (insertedIDs []string, err error)

	UpdateByID(id string, data interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateOne(filter interface{}, data interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	DeleteByID(id string, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	Ctx(seconds uint) (context.Context, context.CancelFunc)
	Err(err error) error
}

// NewDAOMg creates new DAOMg instance with the specified collection
func NewDAOMg(collection *mongo.Collection) *DAOMg {
	return &DAOMg{
		c: collection,
	}
}

// DAOMg is a mongo collection abstraction
type DAOMg struct {
	c *mongo.Collection
}

// FetchByID fetches row by ID to the target
func (d *DAOMg) FetchByID(id string, target interface{}, opts ...*options.FindOneOptions) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	filter := bson.M{"_id": id}
	err := d.C().FindOne(ctx, filter, opts...).Decode(target)

	if err != nil {
		return d.Err(err)
	}

	return nil
}

// FetchByIDs fetches rows by IDs list
func (d *DAOMg) FetchByIDs(ids []string, target interface{}, opts ...*options.FindOptions) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	return d.FetchAllF(target, filter, opts...)
}

// FetchByExIDs fetches rows by exclude IDs list
func (d *DAOMg) FetchByExIDs(ids []string, target interface{}, opts ...*options.FindOptions) error {
	filter := bson.M{"_id": bson.M{"$nin": ids}}
	return d.FetchAllF(target, filter, opts...)
}

// FetchOne fetches one row by the filter
func (d *DAOMg) FetchOne(target interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	ctx, cancel := d.Ctx(1)
	defer cancel()

	err := d.C().FindOne(ctx, filter, opts...).Decode(target)

	if err != nil {
		return d.Err(err)
	}

	return nil
}

// FetchAll fetches all rows from cursor to the target
func (d *DAOMg) FetchAll(target interface{}, opts ...*options.FindOptions) error {
	return d.FetchAllF(target, bson.M{}, opts...)
}

// FetchAllF fetches all rows from cursor with filter to the target
func (d *DAOMg) FetchAllF(target interface{}, filter interface{}, opts ...*options.FindOptions) error {
	ctx, cancel := d.Ctx(10)
	defer cancel()

	cur, err := d.C().Find(ctx, filter, opts...)
	if err != nil {
		return d.Err(err)
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		_ = d.Err(err)
	}(cur, ctx)

	err = cur.All(ctx, target)
	if err != nil {
		return d.Err(err)
	}

	return nil
}

// InsertOne insets row to the collection
func (d *DAOMg) InsertOne(data interface{}, opts ...*options.InsertOneOptions) (id string, err error) {
	ctx, cancel := d.Ctx(3)
	defer cancel()

	doc, err := bson.Marshal(data)
	if err != nil {
		return "", d.Err(err)
	}

	res, err := d.C().InsertOne(ctx, doc, opts...)
	if err != nil {
		return "", d.Err(err)
	}

	id, ok := res.InsertedID.(string)
	if ok {
		return id, nil
	}

	v, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return v.String(), nil
	}

	return "", nil
}

// InsertMany inserts multiple documents to the collection
func (d *DAOMg) InsertMany(rows []interface{}, opts ...*options.InsertManyOptions) (insertedIDs []string, err error) {
	ctx, cancel := d.Ctx(30)
	defer cancel()

	docs := make([]interface{}, len(rows))
	for i, row := range rows {
		doc, err := bson.Marshal(row)
		if err != nil {
			return nil, d.Err(err)
		}
		docs[i] = doc
	}

	res, err := d.C().InsertMany(ctx, docs, opts...)
	if err != nil {
		return nil, d.Err(err)
	}

	insertedIDs = make([]string, len(res.InsertedIDs))
	for i, id := range res.InsertedIDs {
		var ok bool
		var vs string
		var vp primitive.ObjectID

		vs, ok = id.(string)
		if ok {
			insertedIDs[i] = vs
			continue
		}

		vp, ok = id.(primitive.ObjectID)
		if ok {
			insertedIDs[i] = vp.String()
			continue
		}

		insertedIDs[i] = ""
	}

	return insertedIDs, nil
}

// UpdateByID updates one row by ID
func (d *DAOMg) UpdateByID(id string, data interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := d.Ctx(3)
	defer cancel()

	upd := map[string]interface{}{
		"$set": data,
	}
	doc, err := bson.Marshal(upd)
	if err != nil {
		return nil, err
	}

	return d.C().UpdateByID(ctx, id, doc, opts...)
}

// UpdateOne updates one row
func (d *DAOMg) UpdateOne(filter interface{}, data interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := d.Ctx(3)
	defer cancel()

	upd := map[string]interface{}{
		"$set": data,
	}
	doc, err := bson.Marshal(upd)
	if err != nil {
		return nil, err
	}

	return d.C().UpdateOne(ctx, filter, doc, opts...)
}

// DeleteByID deletes one row by ID
func (d *DAOMg) DeleteByID(id string, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return d.DeleteOne(filter, opts...)
}

// DeleteOne deletes one row
func (d *DAOMg) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := d.Ctx(3)
	defer cancel()
	return d.C().DeleteOne(ctx, filter, opts...)
}

// DeleteMany deletes filtered rows
func (d *DAOMg) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := d.Ctx(5)
	defer cancel()
	return d.C().DeleteMany(ctx, filter, opts...)
}

// C returns Collection
func (d *DAOMg) C() *mongo.Collection {
	return d.c
}

// Ctx creates new timeout context
func (d *DAOMg) Ctx(seconds uint) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// Err transforms and log an error if needed
func (d *DAOMg) Err(err error) error {
	if err == nil {
		return nil
	}
	needLog := true
	switch err {
	case mongo.ErrNoDocuments:
		return pagocore.ErrNotFound
	case pagocore.ErrNotFound:
		needLog = false
	}
	if needLog {
		cName := "_unknown_"
		c := d.C()
		if c != nil {
			cName = c.Name()
		}
		log.WithField("dao", cName).Error(err)
	}
	return err
}
