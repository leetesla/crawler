package main

import (
	"tesla/crawler/config"
	"tesla/crawler/pipeline"
	"tesla/crawler/engine"
	"tesla/crawler/scheduler"
	"tesla/crawler/parser"
	"tesla/crawler/parser/spider"
)

func main()  {
	itemChan, err := pipeline.SaveItem(config.ElasticIndex)

	if err != nil {
		panic(err)
	}

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,//100
		ItemChan:         itemChan,
	}

	request := parser.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: parser.NewParser(spider.ParseCityList, config.ParseCityList),
	}

	// fmt.Println("%T", e)
	e.Run(request)
}