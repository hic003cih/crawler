package fetcher

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//把原本在main執行get網站body和檢查網站編碼的功能移過來
//傳入要抓取的網址,並回傳處理過的byte,如果如出錯船error
func Fetch(url string) ([]byte, error) {
	//取得網站body
	resp, err := http.Get(url)
	//錯誤時輸出錯誤
	if err != nil {
		return nil, err
	}
	//關閉response的body
	defer resp.Body.Close()

	//取response的status狀態
	//如果狀態不是OK則回傳nil,並用輔助輸出錯誤的訊息並跳出程式
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}
	//因為不一定每個網頁都是gbk編碼
	//用determineEncoding來辨別是哪種編碼
	e := determineEncoding(resp.Body)

	//使用transform庫
	//將上面轉辨別好的e傳入這邊
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())

	//讀取轉換完以後utf8格式的reader
	return ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
}

//猜這個html的document的encoding是甚麼
//這邊把原來的函數包裝一下
//encoding.Encoding要引入golang.org/x/text/encoding才可以用
func determineEncoding(r io.Reader) encoding.Encoding {
	//直接讀io.Reader的話,1024的byte就沒辦法再讀了
	//所以這邊用 bufio.NewReader(r).Peek來裝一下body的資料(Peek窺視)
	//Peek(s)->取前面s位數,如下面的1024就是取前1024
	//如果是err,為了不讓程式掛掉,直接回傳一個unicode.UTF8,並將error打印出來
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		log.Printf("Fetcher error:%v", err)
		return unicode.UTF8
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
