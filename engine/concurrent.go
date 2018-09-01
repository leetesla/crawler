package engine

import (
	"tesla/crawler/downloader"
	"tesla/crawler/parser"
	"log"
	"tesla/crawler/scheduler"
)

//type Processor func(Request) (ParseResult, error)
type ConcurrentEngine struct {
	Scheduler   scheduler.Scheduler //*QueuedScheduler
	ItemChan    chan parser.Item
	WorkerCount int
	//RequestProcessor processor.Processor//engine.Worker(func) Processor
}

// 运行
func (e *ConcurrentEngine) Run(seeds ...parser.Request) {
	// 典型的值传递类型：数组
	// 典型的引用传递类型：interface, slice，map和channel
	outParseResultChan := make(chan parser.ParseResult)

	e.Scheduler.Run()

	// 创建WorkerCount个worker
	for i := 0; i < e.WorkerCount; i++ {
		// e.Scheduler也实现了ReadyNotifier接口
		//requestChan := e.Scheduler.WorkerChan() //
		inRequestChan := e.Scheduler.MakeRequestChan() //


		//1.通过createWorker创建goroutine，每个worker其实就是一个goroutine，含有一个inRequestChan（称之为worker更好）
		//2.所有的worker共用scheduler（含有requestChan、requestChanChan）和outParseResultChan
		//3.每个createWorker操作都会向scheduler.requestChanChan发送自己的inRequestChan，告诉已准备好，一般会阻塞成一个buf队列
		//4.Scheduler.Submit会不断向scheduler.requestChan发送request(4.1/4.2)
		e.createWorker(inRequestChan, outParseResultChan, e.Scheduler)
	}

	//4.1初始化种子
	for _, r := range seeds {
		if isDuplicate(r.Url) {
			continue
		}
		e.Scheduler.Submit(r)
	}

	//4.2不断循环
	for {
		result := <-outParseResultChan//来自createWorker中的ProcessRequest

		for _, item := range result.Items {
			go func(i parser.Item) {
				e.ItemChan <- i
			}(item)
		}

		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}

	}
}


// 创建worker
//TODO:修改r
//TODO:修改createWorker的名字
func (e *ConcurrentEngine) createWorker(inRequestChan chan parser.Request,
	outParseResultChan chan parser.ParseResult, s scheduler.ReadyNotifier) {
	go func() {

		//不断注册，不断接收；注册一次，接收一次；处理完请求之后继续去注册，接收，处理请求
		for {
			// 即 e.Scheduler.requestChanChan <- inRequestChan
			// 1.1先把inRequestChan发送到requestChanChan并调度到[]requestChanQ中排好队
			// 类似注册，取水
			s.Ready(inRequestChan)

			// 1.2这一步会阻塞，直到Scheduler通过调度将request发送到inRequestChan中
			// request来自e.Scheduler.Submit
			request := <-inRequestChan

			parseResult, err := ProcessRequest(request)

			if err != nil {
				continue
			}

			// out是个传参过来的channel（引用类型）；在Run里面接收
			outParseResultChan <- parseResult
		}
	}()
}


// 处理请求，解析，返回解析后的结果;Fetch+Parse
func ProcessRequest(r parser.Request) (parser.ParseResult, error) {
	body, err := downloader.Fetch(r.Url)

	if err != nil {
		log.Printf("Fetcher: error "+"fetching url %s: %v", r.Url, err)
		return parser.ParseResult{}, err
	}

	// r.Parser=engine.Parser（指针接收者）
	// r.Parser.Parse
	//在具体parser里指定下一级的parser方法
	return r.Parser.Parse(body, r.Url), nil
}



// 判断去重
var visitedUrls = make(map[string]bool)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}