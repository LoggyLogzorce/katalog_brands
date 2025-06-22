package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type Indexer interface {
	IndexBrand(ctx context.Context, brand BrandDoc) error
	DeleteBrand(ctx context.Context, id string) error
}

type ESIndexer struct {
	client *elasticsearch.Client
}

func NewIndexer(client *elasticsearch.Client) *ESIndexer {
	return &ESIndexer{client: client}
}

type BrandDoc struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Status      string `json:"status"`
}

func (e *ESIndexer) IndexBrand(ctx context.Context, b BrandDoc) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	res, err := e.client.Index(
		"brands",
		bytes.NewReader(data),
		e.client.Index.WithDocumentID(b.ID),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing brand ID=%s: %s", b.ID, res.Status())
	}
	return nil
}

func (e *ESIndexer) DeleteBrand(ctx context.Context, id string) error {
	res, err := e.client.Delete(
		"brands",
		id,
		e.client.Delete.WithContext(ctx),
		e.client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error deleting brand ID=%s: %s", id, res.Status())
	}
	return nil
}
