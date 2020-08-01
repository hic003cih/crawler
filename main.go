package main

import (
	"crawler/engine"
	"crawler/scheduler"
	"crawler/zhenai/parser"
)

func main() {
	//爬蟲引擎執行
	e := engine.ConcurrentEngine{
		//把Scheduler取出
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
	}
	//抓全部城市
	/* e.Run(engine.Request{
		//改用localhost生成的資料去執行
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	}) */

	//抓單個城市
	e.Run(engine.Request{
		//改用localhost生成的資料去執行
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun/guangzhou",
		ParserFunc: parser.ParseCity,
	})
}
