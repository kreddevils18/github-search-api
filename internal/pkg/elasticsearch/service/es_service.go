package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/config"
	"github.com/kien-hoangtrung/github-repository/internal/pkg/log"
	"strings"
)

type ElasticSearchResult struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source    interface{} `json:"_source"`
			Highlight interface{} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticSearchService struct {
	logger log.ILogger
	client *elasticsearch.Client
}

func NewElasticSearchService(conf *config.Config, logger log.ILogger) (*ElasticSearchService, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{
			conf.ElasticSearch.Address,
		},
		Username: conf.ElasticSearch.Username,
		Password: conf.ElasticSearch.Password,
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, err
	}
	res, err := es.Info()
	if err != nil {
		logger.Fatalf("error getting response: %s", err)
		return nil, err
	}
	if res.IsError() {
		logger.Fatalf("error: %s", res.String())
		return nil, err
	}
	var (
		r map[string]interface{}
	)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Fatalf("Error parsing the response body: %s", err)
	}
	logger.Infof("Client: %s", elasticsearch.Version)
	logger.Infof("Server: %s", r["version"].(map[string]interface{})["number"])
	logger.Info(strings.Repeat("~", 37))

	return &ElasticSearchService{logger: logger, client: es}, nil
}

func (e ElasticSearchService) IndexRequest(index string, documentId string, data any) error {
	temp, err := json.Marshal(data)
	if err != nil {
		e.logger.Fatalf("Error marshaling document: %s", err)
	}
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       bytes.NewReader(temp),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		e.logger.Fatalf("error getting response: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		e.logger.Infof("[%s] Error indexing document ID=%d", res.Status(), documentId)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			e.logger.Infof("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			e.logger.Infof("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return nil
}

func (e ElasticSearchService) Query(index string, query string) (*esapi.Response, error) {
	res, err := e.client.Search(
		e.client.Search.WithIndex(index),
		e.client.Search.WithBody(strings.NewReader(query)),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
