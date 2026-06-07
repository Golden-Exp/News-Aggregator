package main

import (
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
)

func main() {
	fmt.Println("STARTING...")
	parser := gofeed.NewParser()

	feedURL := "https://www.puthiyathalaimurai.com/api/v1/collections/tamilnadu"

	feed, err := parser.ParseURL(feedURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Feed Title:", feed.Title)
	fmt.Println("Number of items:", len(feed.Items))
}
