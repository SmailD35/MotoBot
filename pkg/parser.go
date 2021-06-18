package pkg

import (
	"bytes"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name 	`xml:"channel"`
	Items   []XMLItem   `xml:"item"`
}

type XMLItem struct {
	XMLName     xml.Name `xml:"item"`
	Author      string   `xml:"author"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Category    string   `xml:"category"`
}

func XMLParser(url string) ([]*Item, error) {
	byteValue, err := GetBody(url)
	if err != nil {
		return nil, err
	}

	var rss RSS

	err = xml.Unmarshal(byteValue, &rss)
	if err != nil {
		return nil, err
	}

	var items []*Item

	for _, item := range rss.Channel.Items {
		date, _ := time.Parse(time.RFC1123Z, item.PubDate)
		i := &Item{
			Author:      item.Author,
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			PubDate:     date,
		}
		items = append(items, i)
	}

	return items, nil
}

func ItemPares(url string) (string, error) {
	byteValue, err := GetBody(url)
	if err != nil {
		return "", err
	}

	page, err := goquery.NewDocumentFromReader(bytes.NewReader(byteValue))
	if err != nil {
		return "", err
	}

	text := ""
	page.Find(".js-topic__text").Each(func(i int, s *goquery.Selection) {
		text += s.Find("p").Text()
	})

	return text, nil
}
