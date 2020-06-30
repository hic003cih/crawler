package parser

import (
	"crawler/fetcher"
	"fmt"
	"testing"
)

func TestParseCityList(t *testing.T) {

	contents, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")

	if err != nil {
		panic(err)
	}
	//因為有可能測試的機器沒有對外的網路,所以把Fetch到的Body存下來
	fmt.Printf("%s\n", contents)

}
