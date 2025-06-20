package es

import "github.com/elastic/go-elasticsearch/v7"

func NewClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Проверим соединение
	_, err = es.Info()
	if err != nil {
		return nil, err
	}

	return es, nil
}
