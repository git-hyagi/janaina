package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	user          = "admin"
	password      = "bqtwqipFIQ3kOLYh"
	mongo_addr    = "10.0.1.25"
	database      = "admin"
	db_port       = "30901"
	db_connection = "mongodb://" + user + ":" + password + "@" + mongo_addr + ":" + db_port + "/" + database + "?retryWrites=true&w=majority"

	db_connection_timeout_seconds = 30
)

type Mongo struct {
	Conn *mongo.Database
}

func main() {

	test := Mongo{}

	client, err := mongo.NewClient(options.Client().ApplyURI(db_connection))
	if err != nil {
		fmt.Printf("%v", err)
	}

	// set connection timeout interval
	ctx, _ := context.WithTimeout(context.Background(), db_connection_timeout_seconds*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer client.Disconnect(ctx)

	//	db, err := client.ListDatabaseNames(ctx, bson.M{})
	//	if err != nil {
	//		fmt.Printf("%v\n", err)
	//	}
	//
	//fmt.Printf("%v\n", db)

	chat_db := client.Database("chat")
	test.Conn = client.Database("chat")
	collection1 := chat_db.Collection("col1")

	_, err = collection1.InsertOne(ctx, bson.D{
		{Key: "id", Value: "1"},
		{Key: "username", Value: "android"},
		{Key: "timestamp", Value: time.Now().Format("03:04")},
		{Key: "message", Value: "Hello World!"},
	})
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	test.GetAllMessages(ctx, "col1")
	fmt.Println("====================================================")
	test.GetMessagesNewerThan(ctx, "col1", "11:00")

}

func (con *Mongo) GetMessagesNewerThan(ctx context.Context, col, newer string) error {
	collection := con.Conn.Collection(col)
	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", 1}})

	sortCursor, err := collection.Find(ctx, bson.D{
		{"timestamp", bson.D{
			{"$gt", newer},
		}},
	})
	var msgSorted []bson.M
	if err = sortCursor.All(ctx, &msgSorted); err != nil {
		fmt.Println(err)
	}

	for _, i := range msgSorted {
		fmt.Printf("[%v] %v: %v\n", i["timestamp"], i["username"], i["message"])
	}

	return err
}

func (con *Mongo) GetAllMessages(ctx context.Context, col string) error {

	var results []bson.M

	collection := con.Conn.Collection(col)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &results)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for _, line := range results {
		fmt.Printf("[%v] %v: %v\n", line["timestamp"], line["username"], line["message"])
	}

	return err
}
