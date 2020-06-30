package main

import (
	"crawler/engine"
	"crawler/zhenai/parser"
)

func main() {
	//爬蟲引擎執行
	engine.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
