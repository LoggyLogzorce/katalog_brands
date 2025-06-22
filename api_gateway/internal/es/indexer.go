package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type Indexer interface {
	SearchProducts(ctx context.Context, query string, size, from int) ([]ProductDoc, int, error)
	SearchBrands(ctx context.Context, query string, size, from int) ([]BrandDoc, int, error)
}

type ESIndexer struct {
	client *elasticsearch.Client
}

func NewIndexer(client *elasticsearch.Client) *ESIndexer {
	return &ESIndexer{client: client}
}

type ProductDoc struct {
	ID          string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
	Status      string  `json:"status"`
}

type BrandDoc struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Status      string `json:"status"`
}

type CategoryDoc struct {
	ID    string `json:"category_id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type SearchHit struct {
	Source json.RawMessage `json:"_source"`
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []SearchHit `json:"hits"`
	} `json:"hits"`
}

func (e *ESIndexer) SearchProducts(ctx context.Context, query string, size, from int) ([]ProductDoc, error) {
	// Строим тело запроса
	body := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2", "category", "description"},
			},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	// Выполняем поиск
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("products"),
		e.client.Search.WithBody(buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("ES search error: %s", res.Status())
	}

	// Парсим ответ
	var r SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	// Преобразуем в []ProductDoc
	out := make([]ProductDoc, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (e *ESIndexer) SearchBrands(ctx context.Context, query string, size, from int) ([]BrandDoc, error) {
	body := map[string]interface{}{
		"from": from, "size": size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2", "description"},
			},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("brands"),
		e.client.Search.WithBody(buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("ES brands search error: %s", res.Status())
	}
	var r SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	out := make([]BrandDoc, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (e *ESIndexer) SearchCategories(ctx context.Context, query string, size, from int) ([]CategoryDoc, error) {
	body := map[string]interface{}{
		"from": from, "size": size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2"},
			},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("categories"),
		e.client.Search.WithBody(buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("ES search error: %s", res.Status())
	}

	var r SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	out := make([]CategoryDoc, len(r.Hits.Hits))
	for i, hit := range r.Hits.Hits {
		if err := json.Unmarshal(hit.Source, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}
