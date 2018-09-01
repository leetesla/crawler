package scheduler

// *QueuedScheduler实现了engine.Scheduler接口
import (
	"tesla/crawler/parser"
)

type QueuedScheduler struct {
	requestChan chan parser.Request
	//workerChan  chan chan parser.Request
	requestChanChan  chan chan parser.Request
}

func (s *QueuedScheduler) MakeRequestChan() chan parser.Request {
	return make(chan parser.Request)
}

// 提交engine.Request
func (s *QueuedScheduler) Submit(r parser.Request) {
	s.requestChan <- r //Go语言自带隐式解引用
}

func (s *QueuedScheduler) Ready(w chan parser.Request) {
	// 向woker管道中发送woker(chan engine.Request)
	//s.workerChan <- w
	s.requestChanChan <- w //Go语言自带隐式解引用
}


//
func (s *QueuedScheduler) Run() {
	//s.workerChan = make(chan chan parser.Request)
	s.requestChanChan = make(chan chan parser.Request)
	s.requestChan = make(chan parser.Request)

	go func() {
		//requestQ和workerQ只初始化一次，与scheduler一对一
		var requestQ []parser.Request
		var requestChanQ []chan parser.Request

		// for循环
		for {
			var activeRequest parser.Request
			//var activeWorker chan parser.Request
			var activeRequestChan chan parser.Request

			if len(requestQ) > 0 && len(requestChanQ) > 0 {
				//activeWorker = workerQ[0]
				//重点：此处activeRequestChan取自workerQ，来自s.requestChanChan
				activeRequestChan = requestChanQ[0]
				activeRequest = requestQ[0]
			}

			// 关键：select 调度
			// 这就是为什么要从engine.run发过来，再接收回去

			select {
			//1.从s.requestChan不断接收request放到requestQ中
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)

			//case w := <-s.workerChan:
			// 2.从s.requestChanChan不断接收requestChan（worker）放到requestChanQ中
			case w := <-s.requestChanChan:
				requestChanQ = append(requestChanQ, w)

			// 3.如果向activeRequestChan发送activeRequest成功(需要有e.createWorker中request := <-inRequestChan取出)，两方队列各减一
			// 注意：activeRequestChan接收自requestChanQ，requestChanQ中的chan来自s.requestChanChan,
			// 而chan是引用传递，来自e.createWorker中s.Ready(inRequestChan)即 s.requestChanChan <- inRequestChan
			case activeRequestChan <- activeRequest:
				requestChanQ = requestChanQ[1:]
				requestQ = requestQ[1:]
			}

		}
	}()
}
