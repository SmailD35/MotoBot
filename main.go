package main

import (
	"context"
	"fmt"
	"github.com/SmailD35/MotoBot/pkg"
	"github.com/SmailD35/MotoBot/proto"
	"google.golang.org/grpc"
	"log"
	"sync"
)

const count = 30
const searchWord = "людей"
const grpcAddr = "[::]:50051"

var url = "https://lenta.ru/rss/articles/russia"

func main() {
	db, err := pkg.NewESClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcConn, err := grpc.Dial(
		grpcAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("error connecting to grpc")
	}
	defer grpcConn.Close()

	fmt.Println("moto")

	client := proto.NewItemServiceClient(grpcConn)

	newItems, err := pkg.XMLParser(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	existingItems, err := db.GetItems(count)
	if err != nil {
		fmt.Println(err)
		return
	}

	existingItemsMap := map[string]*pkg.Item{}
	for _, item := range existingItems {
		existingItemsMap[item.Link] = item
	}

	var newItemsForPars []*pkg.Item

	for _, newItem := range newItems {
		if _, ok := existingItemsMap[newItem.Link]; !ok {
			newItemsForPars = append(newItemsForPars, newItem)
		}
	}

	wg := &sync.WaitGroup{}
	for _, item := range newItemsForPars {
		wg.Add(1)
		go func(i *pkg.Item, wg *sync.WaitGroup) {
			text, err := pkg.ItemPares(i.Link)
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}
			i.Text = text
			wg.Done()
			return
		}(item, wg)
	}
	wg.Wait()

	for _, item := range newItemsForPars {
		signature, err := client.GetSignature(context.Background(), &proto.ItemText{Text: item.Text})
		if err != nil {
			fmt.Println(err)
			return
		}
		item.Signature = signature.Signature
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
