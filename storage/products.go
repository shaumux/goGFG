package storage

import (
	"context"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"goGFG/logger"
)

var validate *validator.Validate

type Product struct {
	Title string  `json:"title,omitempty" validate:"required"`
	Brand string  `json:"brand,omitempty" validate:"required"`
	Price float32 `json:"price,omitempty" validate:"required"`
	Stock int     `json:"stock,omitempty" validate:"required"`
}

func (p *Product) Validate() error {
	validate = validator.New()
	if err := validate.Struct(p); err != nil {
		return err
	}
	return nil
}

const mapping = `{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
			"properties": {
				"title": {
					"type": "keyword"
				},
				"brand": {
					"type": "keyword"
				},
				"price": {
					"type": "float"
				},
				"stock": {
					"type": "integer"
				}
			}
		}
}`

func (p *Product) IndexProduct() error {
	ctx := context.Background()
	client := GetElastic()
	_, err := client.Index().
		Index("products").
		Id(uuid.Must(uuid.NewV4()).String()).
		BodyJson(p).
		Do(ctx)
	if err != nil {
		// Handle error
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Error Indexing")
		return err
	}
	_, err = client.Flush().Index("products").Do(ctx)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Error Indexing")
	}
	return nil
}

func SearchProduct(searchFields []string, searchText string) *elastic.MultiMatchQuery {
	matchQuery := elastic.NewMultiMatchQuery(searchText)
	for _, searchField := range searchFields {
		matchQuery.Field(searchField)
	}
	return matchQuery
}

func FilterProduct(filterField, filterValue string) *elastic.TermQuery {
	termQuery := elastic.NewTermQuery(filterField, filterValue)
	return termQuery
}

func GetProducts(sortField string, pageNo, pageSize *int, filter *elastic.TermQuery, search *elastic.MultiMatchQuery) (*elastic.SearchResult, error) {
	ctx := context.Background()
	client := GetElastic()
	searchService := client.Search().Index("products")
	if sortField != "" {
		searchService.Sort(sortField, true)
	}

	if pageNo != nil {
		if pageSize != nil {
			searchService.From(*pageNo - 1*(*pageSize))
		} else {
			searchService.From((*pageNo - 1) * 100)
		}
	} else {
		searchService.From(0)
	}

	if pageSize != nil {
		searchService.Size(*pageSize)
	} else {
		searchService.Size(100)
	}

	query := elastic.NewBoolQuery()

	if search != nil {
		query.Must(search)
	} else {
		query.Must(elastic.NewMatchAllQuery())
	}

	if filter != nil {
		query.Filter(filter)
	}

	searchResult, err := searchService.Query(query).Do(ctx)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Error Searching")
	}
	return searchResult, err
}
