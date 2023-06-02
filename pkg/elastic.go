package pkg

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(cloudID string, apiKey string) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		// Handle error
		panic(err)
	}
	return es
}
