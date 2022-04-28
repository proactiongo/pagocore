package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// Options is a MongoDB options
type Options struct {
	Hosts    []string
	Database string
	User     string
	Password string
}

// NewMongoDB creates MongoDB instance
func NewMongoDB(opt *Options) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var err error

	mOpt := &options.ClientOptions{
		Hosts: opt.Hosts,
	}

	if opt.User != "" {
		mOpt.SetAuth(options.Credential{
			AuthSource: opt.Database,
			Username:   opt.User,
			Password:   opt.Password,
		})
	}

	client, err := mongo.Connect(ctx, mOpt)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(opt.Database)

	return &MongoDB{
		client: client,
		Db:     db,
	}, nil
}

// MongoDB manages mongo connection
type MongoDB struct {
	client *mongo.Client
	Db     *mongo.Database
}

// Close closes db connection context
func (m *MongoDB) Close() error {
	if m.client != nil {
		return m.client.Disconnect(context.Background())
	}
	return nil
}
