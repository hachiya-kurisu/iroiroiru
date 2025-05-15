package main

import (
	"blekksprut.net/iroiroiru"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func importOccurrences(path string, collection *mongo.Collection) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to read occurrence file %s: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	headers := strings.Split(scanner.Text(), "\t")

	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")

		var rawLatitude, rawLongitude string

		doc := bson.D{}
		for i, raw := range fields {
			v := strings.TrimSpace(raw)
			if v != "" {
				doc = append(doc, bson.E{Key: headers[i], Value: v})
			}
			if i == 0 {
				doc = append(doc, bson.E{Key: "_id", Value: v})
			}

			if headers[i] == "decimalLatitude" && v != "" {
				rawLatitude = v
			}

			if headers[i] == "decimalLatitude" && v != "" {
				rawLongitude = v
			}
		}

		// skip records with no or invalid latitude and/or longitude
		if rawLatitude == "" || rawLongitude == "" {
			fmt.Println("no lat/lon")
			continue
		}

		latitude, err := strconv.ParseFloat(rawLatitude, 64)
		if err != nil {
			fmt.Println("latn")
			continue
		}
		longitude, err := strconv.ParseFloat(rawLongitude, 64)
		if err != nil {
			fmt.Println("tn")
			continue
		}

		point := bson.D{
			{Key: "type", Value: "Point"},
			{Key: "coordinates", Value: bson.A{longitude, latitude}},
		}

		doc = append(doc, bson.E{Key: "location", Value: point})

		_, err = collection.InsertOne(context.TODO(), doc)
		if err != nil {
			fmt.Printf("uh oh: %v\n", err)
			continue
		}

		count++
		if count%10000 == 0 {
			fmt.Printf("%d docs...\n", count)
		}
	}

	fmt.Printf("done! added %d occurrence docs\n", count)
}

func updateMultimedia(path string, collection *mongo.Collection) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to read multimedia file %s: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	headers := strings.Split(scanner.Text(), "\t")

	count := 0
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")

		id := fields[0]

		doc := bson.D{}
		for i, raw := range fields {
			v := strings.TrimSpace(raw)
			if v != "" && headers[i] != "gbifID" {
				doc = append(doc, bson.E{Key: headers[i], Value: v})
			}
		}

		filter := bson.M{"_id": id}
		update := bson.M{"$addToSet": bson.M{"multimedia": doc}}

		result, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			continue
		}

		if result.ModifiedCount > 0 {
			count++
			if count%10000 == 0 {
				fmt.Printf("%d docs...\n", count)
			}
		}
	}
	fmt.Printf("finished adding %d multimedia docs\n", count)
}

func main() {
	occurrencePath := flag.String("o", "", "occurrence data")
	mediaPath := flag.String("m", "", "multimedia data")
	mongoURI := flag.String("u", "mongodb://localhost:27017", "mongodb uri")
	databaseName := flag.String("db", "iroiro", "mongodb database")
	collectionName := flag.String("c", "occurrences", "mongodb collection")
	version := flag.Bool("v", false, "version")

	flag.Parse()

	if *version {
		fmt.Println(iroiroiru.Version)
		return
	}

	if *occurrencePath == "" && *mediaPath == "" {
		flag.Usage()
		fmt.Println()
		fmt.Println("occurrence data (-o) or multimedia data (-m) required")
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

	if *occurrencePath != "" {
		importOccurrences(*occurrencePath, collection)
	}

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "kingdom", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
	}

	_, err = collection.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		log.Fatal(err)
	}

	if *mediaPath != "" {
		updateMultimedia(*mediaPath, collection)
	}
}
