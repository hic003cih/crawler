package engine

import (
	"crawler/fetcher"
	"log"
)

func Run(seeds ...Request) {
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

		log.Printf("Fetching %s", r.Url)
		//用Fetch把body拿出來
		body, err := fetcher.Fetch(r.Url)
		//不能因為有錯誤就panic中斷,要讓其他爬蟲也可以繼續執行,傳出錯以後用continue繼續執行
		if err != nil {
			log.Printf("Fetcher: error"+"fetching url %s: %v", r.Url, err)
			continue
		}
		//取得body以後,送給r.ParserFunc,得到Requests和Items
		parserResult := r.ParserFunc(body)
		//把parserResult加入到requests
		//parserResult.Requests...表示將所有的值寫入
		requests = append(requests, parserResult.Requests...)
		//把parserResult.Items打印出來
		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}

	}
}
