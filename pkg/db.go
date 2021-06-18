package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"os"

	"log"
)

const indexName = "items"

type DBItem struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	PubDate     string `json:"pubDate"`
	Description string `json:"description"`
}

const (
	dateFormat   = "dd.MM.YYYY"
	dateInterval = "1d"
)

type ESClient struct {
	client *elastic.Client
}

func NewESClient() (ESClient, error) {
	client, err := elastic.NewSimpleClient(elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)))
	if err != nil {
		return ESClient{}, err
	}

	es := ESClient{client: client}

	exists, err := es.client.
		IndexExists(indexName).
		Do(context.Background())
	if err != nil {
		log.Println("IndexExists() error", err)
	}

	if exists {
		return es, nil
	}

	_, err = es.client.CreateIndex(indexName).BodyString(PutMapping).Do(context.Background())
	if err != nil {
		return ESClient{}, err
	}

	return es, nil
}

func (es *ESClient) PutItems(items []*Item) error {
	for _, item := range items {
		_, err := es.client.Index().
			Index(indexName).
			Id(item.Link).
			BodyJson(&item).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *ESClient) GetItems(count int) ([]*Item, error) {
	query := elastic.NewMatchAllQuery()
	res, err := es.client.Search().
		Index(indexName).
		Query(query).
		Size(count).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var items []*Item

	for _, hit := range res.Hits.Hits {
		var t *Item
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}

		items = append(items, t)
	}

	return items, nil
}

func (es *ESClient) Search(key string) ([]Item, error) {
	queryString := elastic.NewQueryStringQuery(key)

	res, err := es.client.Search().
		Index(indexName).
		Query(queryString).
		Do(context.Background())
	if err != nil {
		return []Item{}, err
	}

	if res.Hits.TotalHits.Value == 0 {
		return []Item{}, errors.New("Nothing found by Search() ")
	}

	var items []Item

	for _, hit := range res.Hits.Hits {
		var t Item
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}

		items = append(items, t)
	}

	return items, nil
}

func (es *ESClient) TermAggregationByField(field string) ([]Aggregation, error) {
	aggregationQuery := elastic.NewTermsAggregation().
		Field(addKeyWord(field)).
		Size(30).
		OrderByCountDesc()

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, aggregationQuery).
		Do(context.Background())
	if err != nil {
		return []Aggregation{}, err
	}

	rawMsg := result.Aggregations[indexName]

	ar := elastic.AggregationBucketKeyItems{}

	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return nil, err
	}

	var termsAggregations []Aggregation

	for _, item := range ar.Buckets {
		termsAggregations = append(termsAggregations, Aggregation{
			Key:   item.Key,
			Count: item.DocCount,
		})
	}
	return termsAggregations, nil
}

func (es *ESClient) CardinalityAggregationByField(field string) (float64, error) {
	aggregationQuery := elastic.NewCardinalityAggregation().Field(addKeyWord(field))

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, aggregationQuery).
		Do(context.Background())
	if err != nil {
		return 0, err
	}

	rawMsg := result.Aggregations[indexName]

	var ar elastic.AggregationValueMetric

	err = json.Unmarshal(rawMsg, &ar)
	if err != nil {
		return 0, err
	}

	return *ar.Value, nil
}

func (es *ESClient) DateAggregation() ([]Aggregation, error) {
	dailyAggregation := elastic.NewDateHistogramAggregation().
		Field("pub_date").
		CalendarInterval(dateInterval).
		Format(dateFormat)

	result, err := es.client.Search().
		Index(indexName).
		Aggregation(indexName, dailyAggregation).
		Do(context.Background())
	if err != nil {
		return []Aggregation{}, err
	}

	hist, found := result.Aggregations.Histogram(indexName)
	if !found {
		return []Aggregation{}, errors.New("Nothing found using DateAggregation ")
	}

	var dateHistogramAggregations []Aggregation

	for _, bucket := range hist.Buckets {
		dateHistogramAggregations = append(dateHistogramAggregations, Aggregation{
			Key:   *bucket.KeyAsString,
			Count: bucket.DocCount,
		})
	}

	return dateHistogramAggregations, nil
}
