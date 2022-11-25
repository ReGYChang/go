package main

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	// Used to create a singleton object of Elasticsearch client.
	// Initialized and exposed through GetClient().
	client *es.Client

	// Used to execute client creation procedure only once.
	once sync.Once
)

func GetClient() (*es.Client, error) {
	var err error

	once.Do(func() {
		cfg := es.Config{
			Addresses: Elasticsearch.Hosts,
			Username:  Elasticsearch.Username,
			Password:  Elasticsearch.Password,
		}
		client, err = es.NewClient(cfg)
		if err != nil {
			return
		}
	})

	return client, err
}

func IndexRequest(client *es.Client, req *esapi.IndexRequest) error {
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	indexResponse := &IndexResponse{}
	if err = json.NewDecoder(res.Body).Decode(indexResponse); err != nil {
		return err
	}

	// error handle
	if res.IsError() {
		return errors.New(indexResponse.Result)
	}

	return nil
}

func SearchRequest(client *es.Client, req *esapi.SearchRequest) ([]*SearchHit, error) {
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	searchResult := &SearchResult{}
	if err = json.NewDecoder(res.Body).Decode(searchResult); err != nil {
		return nil, err
	}

	// error handle
	if res.IsError() {
		return nil, &Error{Status: res.StatusCode, Details: searchResult.Error}
	} else if searchResult.TotalHits() == 0 {
		return nil, errors.New("elastic: found no documents")
	}

	searchResult.Header = res.Header
	return searchResult.Hits.Hits, nil
}

func UpdateRequest(client *es.Client, req *esapi.UpdateRequest) error {
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	updateResponse := &UpdateResponse{}
	if err = json.NewDecoder(res.Body).Decode(updateResponse); err != nil {
		return err
	}

	// error handle
	if res.IsError() {
		return errors.New(updateResponse.Result)
	}

	return nil
}

func DeleteRequest(client *es.Client, req *esapi.DeleteRequest) error {
	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	deleteResponse := &DeleteResponse{}
	if err = json.NewDecoder(res.Body).Decode(deleteResponse); err != nil {
		return err
	}

	// error handle
	if res.IsError() {
		return errors.New(deleteResponse.Result)
	}

	return nil
}
