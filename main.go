package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func main() {

	//取得網站body
	resp, err := http.Get(
		"http://www.zhenai.com/zhenghun",
	)
	//錯誤時輸出錯誤
	if err != nil {
		panic(err)
	}
	//關閉response的body
	defer resp.Body.Close()

	//取response的status狀態
	//如果狀態不是OK則輸出錯誤的訊息並跳出程式
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}
	//因為不一定每個網頁都是gbk編碼
	//用determineEncoding來辨別是哪種編碼
	e := determineEncoding(resp.Body)

	//使用transform庫
	//將上面轉辨別好的e傳入這邊
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())

	//讀取轉換完以後utf8格式的reader
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", all)
	printCityList(all)
}

//猜這個html的document的encoding是甚麼
//這邊把原來的函數包裝一下
//encoding.Encoding要引入golang.org/x/text/encoding才可以用
func determineEncoding(r io.Reader) encoding.Encoding {
	//直接讀io.Reader的話,1024的byte就沒辦法再讀了
	//所以這邊用 bufio.NewReader(r).Peek來裝一下body的資料(Peek窺視)
	//Peek(s)->取前面s位數,如下面的1024就是取前1024
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	//要傳一開始的1024的byte
	//上面bufio.NewReader(r).Peek好以後,就可以直接用
	//會收到一個e->encoding
	//name,和certain ->是否確認
	//name和certain先不用管
	e, _, _ := charset.DetermineEncoding(bytes, "")
	//最後把e->encoding傳回
	return e
}
func printCityList(contents []byte) {
	//使用正則表達式把城市名稱取出
	//
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" [^>]*>([^<]+)</a>`)
	matches := re.FindAll(contents, -1)
	for _, m := range matches {
		for _, subMatch := range m {
			fmt.Printf("%s", subMatch)
		}
		fmt.Printf("%s\n", m)
	}
	fmt.Printf("matches found:%d\n", len(matches))
}
