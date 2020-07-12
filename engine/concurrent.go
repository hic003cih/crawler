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
type Scheduler interface {
	Submit(Request)
	ConfigureMasterWorkerChan(chan Request)
	WorkerReady(chan Request)
	Run()
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
		//把輸入和輸出chan傳入
		createWorker(in, out)
	}

	//把每個Request傳到scheduler裡面
	//等in channel和out channel都傳入以後再submit
	for _, r := range seeds {
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
		//然後把每個out輸出的resultRequest,丟給Scheduler
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}
func createWorker(in chan Request, out chan ParserResult) {
	//做一個go routine
	go func() {
		for {
			//告訴SCHEDULER我已經完成

			//從in把值取出來,這個是一個worker的輸入
			//輸入哪裡來?SCHEDULER選擇了你,就會給你發送數據
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
