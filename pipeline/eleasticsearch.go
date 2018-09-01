package pipeline

import (
	"context"
	"log"

	"errors"

	"gopkg.in/olivere/elastic.v5"
	"tesla/crawler/parser"
)

func SaveItem(index string) (chan parser.Item, error) {
	// Must turn off sniff in docker
	client, err := elastic.NewClient(elastic.SetSniff(false))

	if err != nil {
		return nil, err
	}

	out := make(chan parser.Item)

	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+ "#%d: %v", itemCount, item)
			itemCount++

			err := Save(client, index, item)

			if err != nil {
				log.Printf("Item Saver: error "+ "saving item %v: %v", item, err)
			}
		}
	}()

	return out, nil
}

func Save(client *elastic.Client, index string, item parser.Item) error {

	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := client.Index().Index(index).Type(item.Type).BodyJson(item)

	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.Do(context.Background())

	return err
}
