package engine

import (
	"crawler/fetcher"
	"log"
)

type SimpleEngine struct{}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	//把seeds放進requests裡
	for _, r := range seeds {
		requests = append(requests, r)
	}
	//只要有requests就執行
	for len(requests) > 0 {
		r := requests[0]
		//把第一個request拿出來,每個都去做Fetch處理
		requests = requests[1:]

		//調用拆出去的worker來執行parser
		parserResult, err := worker(r)

		//如果err為空就繼續執行
		if err != nil {
			continue
		}
		//把parserResult加入到requests
		//parserResult.Requests...表示將所有的值寫入
		requests = append(requests, parserResult.Requests...)
		//把parserResult.Items打印出來
		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}

	}
}

//把Fetch和Parser提出來做成一個worker,用來準備做並發爬蟲
func worker(r Request) (ParserResult, error) {
	log.Printf("Fetching %s", r.Url)
	//用Fetch把body拿出來
	body, err := fetcher.Fetch(r.Url)
	//不能因為有錯誤就panic中斷,要讓其他爬蟲也可以繼續執行,傳出錯以後用continue繼續執行
	if err != nil {
		log.Printf("Fetcher: error"+"fetching url %s: %v", r.Url, err)
		//因為要回傳ParserResult結構,所以沒辦法直接傳nil,得傳空的結構
		return ParserResult{}, err
	}
	//取得body以後,送給r.ParserFunc,得到Requests和Items
	return r.ParserFunc(body), nil
}
