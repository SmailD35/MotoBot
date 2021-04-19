package main

import (
	"fmt"
	"log"
	"motobot/pkg"
)

var url = "https://lenta.ru/rss/articles/russia"

func main()  {
	items, err := pkg.Parser(url)
	if err != nil {
		log.Fatal(err)
	}

	err = pkg.DB(items)
	if err != nil {
		log.Fatal(err)
	}

	PrintItems(items)
}

func PrintItems(items []pkg.Item)  {
	for _, item := range items {
		fmt.Println("Title: " + item.Title)
		fmt.Println("PubDate: " + item.PubDate)
		fmt.Println("Author: " + item.Author)
		fmt.Println("Link: " + item.Link)
		fmt.Println("Description: " + item.Description)
	}
}
