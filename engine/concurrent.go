package engine

import "log"

type ConcurrentEngine struct {
	//這邊都要大寫,因為是給外面的人用的
	//在這邊定義一個Scheduler
	Scheduler Scheduler
	//建立Worker的數量
	WorkerCount int
}

//使用者定義Scheduler要什麼,然後你自己去實現
//simple和queued都實現下面幾個功能,就可以自由切換了
type Scheduler interface {
	Submit(Request)
	//我有一個worker要給我哪個schedule
	WorkerChan() chan Request
	Run()
	//用組合的方式把ReadyNotifier抓近來
	ReadyNotifier
}

//原本WorkerReady放在Scheduler裡面,太多功能了,所以把WorkerReady拉出來
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {

	//輸入,把Request傳入
	//in := make(chan Request)

	//輸出,ParseResult傳出
	out := make(chan ParserResult)

	//把in channel傳進去ConfigureWorkerChan
	//e.Scheduler.ConfigureMasterWorkerChan(in)

	//這邊改成用Scheduler裡的Run,來生成wokerChan和requestChan
	//然後執行裡面的go func,等待任務的到來
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		//把out channel和scheduler傳進去
		//去跟e.scheduler要workchannel
		//因為Scheduler也有ReadyNotifier,所以可以直接傳
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	//把每個Request傳到scheduler裡面
	//等in channel和out channel都傳入以後再submit
	for _, r := range seeds {

		//檢查是否有重複
		if isDuplicate(r.Url) {
			log.Printf("Duplicate request:"+"%s", r.Url)
			continue
		}
		e.Scheduler.Submit(r)
	}
	itemCount := 0
	for {
		//要把createWorker輸出的結果out收進來
		result := <-out
		//然後把每個out輸出的resultItem輸出
		for _, item := range result.Items {
			log.Printf("Got item #%d: %v", itemCount, item)
			itemCount++
		}

		//URL dedup

		//然後把每個out輸出的resultRequest,丟給Scheduler
		for _, request := range result.Requests {
			//檢查是否有重複
			if isDuplicate(request.Url) {
				log.Printf("Duplicate request:"+"%s", request.Url)
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

//只管從外面接收傳進來channel來造一個worker
//然後自己建立一個channel
//把後面自己拆出來的ReadyNotifier傳進區
func createWorker(in chan Request, out chan ParserResult, ready ReadyNotifier) {

	//做一個go routine
	go func() {
		for {
			//告訴SCHEDULER我自己的channel已經完成
			ready.WorkerReady(in)
			//從in把值取出來,這個是一個worker的輸入
			//輸入哪裡來?SCHEDULER選擇了你,就會給你發送數據
			//然後把in channel傳給request
			request := <-in
			//呼叫worker來把request拆解
			result, err := worker(request)
			//
			if err != nil {
				continue
			}
			//最後把result傳給out
			//然後這個out channel 會把資料傳到Run裡面的result := <-out
			//將資料打印和執行Scheduler.Submit
			out <- result
		}
	}()
}

//建造檢查url的hash map
var visitedUrls = make( map[string])bool

func isDuplicate(url string) bool {
	//如果傳入的url有在visitedUrls map內,
	if visitedUrls[url]{
		return true
	}
	if visitedUrls[url] =true
	return false
}
