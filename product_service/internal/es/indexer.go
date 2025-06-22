package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
)

type IndexerRepository interface {
	IndexProduct(ctx context.Context, p ProductDoc) error
	DeleteProduct(ctx context.Context, id string) error
	IndexCategory(ctx context.Context, c CategoryDoc) error
	DeleteCategory(ctx context.Context, id string) error
}

type repoES struct {
	client *elasticsearch.Client
}

func NewIndexer(client *elasticsearch.Client) IndexerRepository {
	return &repoES{client: client}
}

type ProductDoc struct {
	ID          string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Photo       string  `json:"photo"`
	Status      string  `json:"status"`
}

type CategoryDoc struct {
	ID    string `json:"category_id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

func (e *repoES) IndexProduct(ctx context.Context, p ProductDoc) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	res, err := e.client.Index(
		"products",            // index name
		bytes.NewReader(data), // document body
		e.client.Index.WithDocumentID(p.ID),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing product ID=%s: %s", p.ID, res.Status())
	}
	return nil
}

func (e *repoES) DeleteProduct(ctx context.Context, id string) error {
	res, err := e.client.Delete(
		"products",
		id,
		e.client.Delete.WithContext(ctx),
		e.client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error deleting product ID=%s: %s", id, res.Status())
	}
	return nil
}

func (e *repoES) IndexCategory(ctx context.Context, c CategoryDoc) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	res, err := e.client.Index(
		"categories",
		bytes.NewReader(data),
		e.client.Index.WithDocumentID(c.ID),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error indexing product ID=%s: %s", c.ID, res.Status())
	}
	return nil
}

func (e *repoES) DeleteCategory(ctx context.Context, id string) error {
	res, err := e.client.Delete(
		"categories",
		id,
		e.client.Delete.WithContext(ctx),
		e.client.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("error deleting product ID=%s: %s", id, res.Status())
	}
	return nil
}
