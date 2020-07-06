package main

import (
	"crawler/engine"
	"crawler/zhenai/parser"
)

func main() {
	//爬蟲引擎執行
	engine.SimpleEngine{}.Run(engine.Request{
		//改用localhost生成的資料去執行
		Url:        "http://localhost:8080/mock/www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
