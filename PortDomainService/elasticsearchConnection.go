package main

import (
	Ports "client-domain-elasticsearch/PortsCommunication/PortsCommunication"
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
)

// ElasticsearchCommunication is used to encapsulate communication with elasticsearch
type ElasticsearchCommunication struct {
	client    *elastic.Client
	batch     *elastic.BulkService
	batchSize int
	indexName string
}

// Initalise the ElasticsearchCommunication structure to set the elasticsearch host, index name and batch size
func (c *ElasticsearchCommunication) Initalise(host string, index string, batchSize int) error {
	client, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		return err
	}
	c.client = client
	c.batch = client.Bulk().Index(index)
	c.batchSize = batchSize
	c.indexName = index
	return nil
}

// GetPort queries the Elasticsearch index for a document with the specified ID and returns the top result
func (c *ElasticsearchCommunication) GetPort(ctx context.Context, id string) (*Ports.Port, error) {
	match := elastic.NewMatchQuery("id", id)
	response, err := c.client.Search(c.indexName).Query(match).Do(ctx)
	if err != nil {
		return &Ports.Port{}, nil
	}

	if response.TotalHits() > 0 {
		retval := Ports.Port{}
		err = json.Unmarshal(response.Hits.Hits[0].Source, &retval)
		if err != nil {
			return &Ports.Port{}, nil
		}
		return &retval, nil
	}
	return &Ports.Port{}, nil
}

// AddDocument adds the specified Port Document to the batch of documents ready to process
// if the batch size exceeds the specified limit, it indexes the batch
func (c *ElasticsearchCommunication) AddDocument(ctx context.Context, port *Ports.Port) error {
	var err error
	c.batch.Add(elastic.NewBulkIndexRequest().Id(port.Id).Doc(port))
	if c.batch.NumberOfActions() > c.batchSize {
		err = c.IndexBatch(ctx)
	}
	return err
}

// IndexBatch indexes the current batch of documents into elasticsearch
func (c *ElasticsearchCommunication) IndexBatch(ctx context.Context) error {
	res, err := c.batch.Do(ctx)
	if err != nil {
		return err
	}
	if res.Errors {
		return errors.New("bulk commit failed")
	}
	return nil
}
