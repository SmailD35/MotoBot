package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"sync"
)

const indexName = "items"

type DBItem struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	PubDate     string `json:"pubDate"`
	Description string `json:"description"`
}

func DB(items []Item) error {
	var wg sync.WaitGroup

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	for _, item := range items {
		wg.Add(1)

		go writeItem(item, &wg, es)
	}

	wg.Wait()

	return nil
}

func writeItem(item Item, wg *sync.WaitGroup, es *elasticsearch.Client)  {
	defer (*wg).Done()

	bdItem := DBItem{
		Author:      item.Author,
		Title:       item.Title,
		Link:        item.Link,
		PubDate:     item.PubDate,
		Description: item.Description,
	}

	b, err := json.Marshal(bdItem)
	if err != nil {
		log.Fatalf("Error marshal item: %s", err)
	}

	itemID := hashMD5(bdItem)

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: itemID,
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%s", res.Status(), itemID)
	} else {
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
}
