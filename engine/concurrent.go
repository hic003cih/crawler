package engine

import (
	"fmt"
)

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
}

func (e ConcurrentEngine) Run(seeds ...Request) {
	//把每個Request傳到scheduler裡面
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	//輸入,把Request傳入
	in := make(chan Request)
	//輸出,ParseResult傳出
	out := make(chan ParserResult)
	for i := 0; i < e.WorkerCount; i++ {
		//把輸入和輸出chan傳入
		createWorker(in, out)
	}

	for {
		//要把createWorker輸出的結果out收進來
		result := <-out
		//然後把每個out輸出的resultItem輸出
		for _, item := range result.Items {
			fmt.Printf("Got item: %v", item)
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
			//從in把值取出來
			request := <-in
			//呼叫worker來把request拆解
			result, err := worker(request)
			//
			if err != nil {
				continue
			}
			//最後把result傳給out
			out <- result
		}
	}()
}
