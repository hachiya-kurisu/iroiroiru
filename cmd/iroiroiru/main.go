package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"blekksprut.net/iroiroiru"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type PointQuery struct {
	Lat float64 `form:"lat" binding:"required,min=-90,max=90"`
	Lon float64 `form:"lon" binding:"required,min=-180,max=180"`
}

type Token struct {
	Token      string    `bson:"token"`
	GbifID     string    `bson:"gbifID"`
	ValidAfter time.Time `bson:"validAfter"`
	ValidUntil time.Time `bson:"validUntil"`
}

func here(c *gin.Context) {
	var query PointQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": "invalid lat long"},
		)
		return
	}

	occurrences := client.Database("iroiro").Collection("occurrences")

	point := bson.M{
		"type":        "Point",
		"coordinates": []float64{query.Lon, query.Lat},
	}

	pipeline := mongo.Pipeline{}

	pipeline = append(pipeline, bson.D{
		{"$geoNear", bson.D{
			{"near", point},
			{"distanceField", "dist.calculated"},
			{"spherical", true},
			{"maxDistance", 1000},
		}},
	})

	kingdom := c.DefaultQuery("kingdom", "")
	if kingdom != "" {
		pipeline = append(pipeline, bson.D{
			{"$match", bson.D{{"kingdom", kingdom}}},
		})
	}

	pipeline = append(pipeline, bson.D{{"$limit", 25}})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := occurrences.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError, gin.H{"error": "aggregation failed"},
		)
		return
	}
	defer cursor.Close(ctx)

	results := make([]bson.M, 0)
	err = cursor.All(ctx, &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decoding trouble"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func main() {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatalf("failed to connect to mongodb instance: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to ping mongodb instance: %v", err)
	}

	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("failed to disconnect from mongodb")
		}
	}()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost"},
		AllowMethods: []string{
			http.MethodGet,
		},
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "iroiroiru",
			"version": iroiroiru.Version,
		})
	})

	r.GET("/here", here)

	r.Run()
}
