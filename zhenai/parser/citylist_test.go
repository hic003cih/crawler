package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {

	//防止網站改版或是測試的機器沒有對外的網路,所以把Fetch到的Body存下來成一個檔案,來進行測試比對
	//contents, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	contents, err := ioutil.ReadFile("citylist_test_data.html")
	//有錯就跳出
	if err != nil {
		panic(err)
	}
	//將讀取本地檔案得到的contents byte內容傳給ParseCityList(程式解析器)做解析
	result := ParseCityList(contents)

	//因為有可能測試的機器沒有對外的網路,所以把Fetch到的Body存下來
	//fmt.Printf("%s\n", contents)

	//原本有470個城市,所以訂一個常數為470
	const resultSize = 470
	//先人工去將三個城市的url存下來,用來做比對
	expectedUrls := []string{
		"http://www.zhenai.com/zhenghun/aba",
		"http://www.zhenai.com/zhenghun/akesu",
		"http://www.zhenai.com/zhenghun/alashanmeng",
	}
	//如果result.Requests(取到的城市)數量不等於上面我們預期的470個城市數量
	//則印出錯誤
	if len(result.Requests) != resultSize {
		t.Errorf("result should have %d "+
			"requests; but had %d",
			resultSize, len(result.Requests))
	}
	//將人工抓出來的url做迴圈比對
	//和每個result.Requests(取到的城市)的url對比,如果都沒有
	//則返回錯誤
	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but "+
				"was %s",
				i, url, result.Requests[i].Url)
		}
	}
}
