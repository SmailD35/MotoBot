package pkg

import "time"

type Item struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	Text        string    `json:"text"`
	Signature   []uint64  `json:"signature"`
}

type Aggregation struct {
	Key   interface{}
	Count int64
}

const PutMapping = `
{"settings": {
        "analysis": {
            "filter": {
                "delimiter": {
                    "type": "word_delimiter",
                    "preserve_original": "true"
                },
                "jmorphy2_russian": {
                    "type": "jmorphy2_stemmer",
                    "name": "ru"
                }
            },
            "analyzer": {
                "text_ru": {
                    "tokenizer": "standard",
                    "filter": [
                        "lowercase",
                        "delimiter",
                        "jmorphy2_russian"
                    ]
                }
            }
        }
    },
    "mappings": {
		"dynamic": "strict",
		"properties":{
			"title":{
				"type": "text",
				"analyzer":"text_ru",
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
				"analyzer":"text_ru"
			},
			"pub_date":{
				"type":"date"
			},
			"text":{
				"type":"text",
				"analyzer":"text_ru"
			},
			"signature":{
				"type":"keyword"
			}
		}
  	}
}`

func addKeyWord(field string) string {
	if field == "title" || field == "description" {
		return field + ".keyword"
	}

	return field
}
