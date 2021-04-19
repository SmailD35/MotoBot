package pkg

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Guid        string   `xml:"guid"`
	Author      string   `xml:"author"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	Category    string   `xml:"category"`
}

func Parser(url string) ([]Item, error) {
	byteValue, err := GetBody(url)
	if err != nil {
		return nil, err
	}

	var rss RSS

	err = xml.Unmarshal(byteValue, &rss)
	if err != nil {
		return nil, err
	}

	return rss.Channel.Items, nil
}
