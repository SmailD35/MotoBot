package main

import (
	"fmt"
	"github.com/SmailD35/MotoBot/pkg"
)

const count = 30
const searchWord = "людей"

var url = "https://lenta.ru/rss/articles/russia"

func main() {
	db, err := pkg.NewESClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	newItems, err := pkg.XMLParser(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	existsItems, err := db.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	existsItemsMap := map[string]*pkg.Item{}
	for _, item := range existsItems {
		existsItemsMap[item.Link] = item
	}

	var newItemsForPars []*pkg.Item

	for _, newItem := range newItems {
		if _, ok := existsItemsMap[newItem.Link]; !ok {
			newItemsForPars = append(newItemsForPars, newItem)
		}
	}

	for _, item := range newItemsForPars {
		text, err := pkg.ItemPares(item.Link)
		if err != nil {
			continue
		}
		item.Text = text
	}

	err = db.PutItems(newItemsForPars)
	if err != nil {
		fmt.Println(err)
		return
	}

	search, err := db.Search(searchWord)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Found by word \"%s\"\n", searchWord)
	for _, i := range search {
		fmt.Println("Title: ", i.Title)
		fmt.Println("Link: ", i.Link, "\n")
	}

	termAggregation, err := db.TermAggregationByField("author")
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("\n\nTerm aggregation by %s\n", "author")
	for _, t := range termAggregation {
		fmt.Printf("%s: %v\n", "Author", t.Key)
		fmt.Printf("Topics count: %d\n\n", t.Count)
	}

	cardinalityAggregation, err := db.CardinalityAggregationByField("author")
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Printf("\n\nCardinality aggregation by %s: %.0f\n\n", "author", cardinalityAggregation)

	dataHistogram, err := db.DateAggregation()
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("\n\nDate histogram aggregation")
	for _, t := range dataHistogram {
		fmt.Printf("Key: %v\n", t.Key)
		fmt.Printf("Count: %d\n\n", t.Count)
	}
}
