package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"blekksprut.net/iroiroiru"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func download(url string, path string) error {
	fmt.Printf("downloading %s to %s\n", url, path)
	return nil

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, resp.Body)
	return err
}

func hash(raw string) string {
	hash := md5.Sum([]byte(raw))
	hashed := hex.EncodeToString(hash[:])

	uri, err := url.Parse(raw)
	if err != nil {
		return hashed
	}
	ext := path.Ext(uri.Path)
	return hashed + ext
}

func main() {
	mongoURI := flag.String("u", "mongodb://localhost:27017", "mongodb uri")
	databaseName := flag.String("db", "iroiro", "mongodb database")
	collectionName := flag.String("c", "occurrences", "mongodb collection")
	output := flag.String("o", "scrape", "output directory")
	version := flag.Bool("v", false, "version")

	flag.Parse()

	if *version {
		fmt.Println(iroiroiru.Version)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(*mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database(*databaseName).Collection(*collectionName)

	pipeline := mongo.Pipeline{
		{{"$unwind", "$multimedia"}},
		{{"$project", bson.M{"identifier": "$multimedia.identifier"}}},
		{{"$match", bson.M{"identifier": bson.M{"$type": "string"}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {

		var result struct {
			Identifier string `bson:"identifier"`
		}

		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("decode error:", err)
			continue
		}

		uri := result.Identifier
		if strings.HasPrefix(uri, "https://inaturalist-open-data.s3") {
			uri = strings.Replace(uri, "original.", "medium.", 1)
		}

		hashed := hash(result.Identifier)
		path := filepath.Join(*output, hashed)

		err = download(uri, path)
		if err != nil {
			fmt.Println("download failed:", err)
		}
	}
}
