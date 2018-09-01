package scheduler

import "tesla/crawler/parser"

type ReadyNotifier interface {
	Ready(chan parser.Request)
}

type Scheduler interface {
	ReadyNotifier  //接口组合
	Submit(parser.Request)
	//WorkerChan() chan parser.Request
	MakeRequestChan() chan parser.Request
	Run()
}
