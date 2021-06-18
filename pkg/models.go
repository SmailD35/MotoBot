package pkg

import "time"

type Item struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	Category    string    `json:"category"`
	Text        string    `json:"text"`
	Signature   []uint64  `json:"signature"`
}

type Aggregation struct {
	Key   interface{}
	Count int64
}

const PutMapping = `
	{
	  "properties":{
			"title":{
				"type": "text",
				"analyzer":"russian",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        		}
			},
			"author": {
				"type":"keyword"
			},
			"link":{
				"type":"keyword"
			},
			"description":{
				"type":"text",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        	},
				"analyzer":"russian"
			},
			"pub_date":{
				"type":"date"
			},
			"text":{
				"type":"text",
				"analyzer":"russian"
			}
	 }
}`

func addKeyWord(field string) string {
	if field == "title" || field == "description" {
		return field + ".keyword"
	}

	return field
}
