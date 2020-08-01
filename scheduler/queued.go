package scheduler

import "crawler/engine"

type QueuedScheduler struct {
	requestChan chan engine.Request
	//將每一個不同的小channel集合在一起存到一個大的channel
	//每個小channel就是worker的channel,
	//再從裡面選擇
	workerChan chan chan engine.Request
}

func (e *QueuedScheduler) Submit(r engine.Request) {
	e.requestChan <- r
}

//每個WorkerChannel都有自己的一個channel
func (e *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

//告訴我們有一個worker好了,可以接收request
//然後把這個worker channel送進去s.wokerChan
//workerChannel的型態是chan engine.Request
func (e *QueuedScheduler) WorkerReady(w chan engine.Request) {
	e.workerChan <- w
}

//總控的真正的goroutine
func (e *QueuedScheduler) Run() {
	//這裡先生成channel,讓他們可以做接下來的事情
	//因為要改變s QueuedScheduler內的內容,所以上面要用s *QueuedScheduler,傳址來改變s內的內容
	e.workerChan = make(chan chan engine.Request)
	e.requestChan = make(chan engine.Request)
	//goroutine裡面用一個for 不斷執行
	//兩個channel都要拿,不管類型是甚麼
	go func() {
		//做用來儲存request和worker的對列
		//不要使用requestQ :=,直接做一個類型聲明
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			//當同時request和worker都有人在空閒等待排隊的時候,就執行下面的將對列中的第一個存進activeWorker和activeRequest
			//有值寫入activeWorker就會執行下面的select case activeWorker <- activeRequest
			//都沒有就都不會執行到
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			//因為有可能兩個channel,可能只有一個有,也可能同時有,要分開獨立收channel
			//用select,哪個有的時候,就執行哪個
			select {
			//第一個requestChan chan engine.Request 從s.requestChan拿資料
			//只要Submit有執行,這邊就會執行
			case r := <-e.requestChan:
				//當requestChan友值,放進去requestQ對列
				requestQ = append(requestQ, r)
			//第二個wokerChan chan chan engine.Request 從s.wokerChan拿資料
			//只要WorkerReady有執行,這邊就會執行
			//當wokerChan友值,放進去requestQ對列
			case w := <-e.workerChan:
				workerQ = append(workerQ, w)
			//當activeRequest把值傳入activeWorker
			//把index[0]的worker和request從對列中刪除
			case activeWorker <- activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
