package storage

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"goGFG/logger"
	"os"
)

var esClient *elastic.Client

func init() {
	esHost := os.Getenv("ESHOST")
	esPORT := os.Getenv("ESPORT")

	esURI := fmt.Sprintf("http://%s:%s", esHost, esPORT)
	logger.Log.WithFields(logrus.Fields{"esURI": esURI}).Info("Connecting to ElasticSearch")
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(esURI), elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	info, code, err := client.Ping(esURI).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	logger.Log.WithFields(logrus.Fields{"code": code, "version": info.Version.Number}).Info("Elasticsearch returned with code and version")

	exists, err := client.IndexExists("products").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("products").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	esClient = client
}

func GetElastic() *elastic.Client {
	return esClient
}
