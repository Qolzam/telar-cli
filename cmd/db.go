package cmd

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func checkDB() error {
	dbURL := "mongodb+srv://telar_user:pass@cluster0-l6ojz.mongodb.net/test?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(dbURL)))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if isError(err) {
		return err
	}
	return nil
}
