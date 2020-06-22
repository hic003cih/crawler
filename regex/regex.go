package main

import (
	"fmt"
	"regexp"
)

//把這句話後面的email找出來
//這邊要用``才可以直接用enter換行,用""不可以
const text = `My email is 55151@gmail.com@WS
email1 is abc@def.org
email2 is gggxcq@qq.com
email3 is jjweq@qq.com.tw
`

func main() {
	//得到一個正則表達式
	// . ->表示可以是字母也可以是數字
	// + ->表示一個或者多個
	// * ->表示零個或者多個
	//中間的 . 不能直接用\. ,因為會被go誤認,要用\\. 兩個\\
	//也可以用將 ".+@.+\\..+" 用 `.+@.+\\..+`  就不會被GO誤認
	//.+會使前面的My email is也被輸出出來
	//因此在+前面[a-zA-Z0-9],表示必須是大小寫和數字才會匹配
	//後面也必須一樣處裡
	//第二個[a-zA-Z0-9]加上.表示允許多個.
	/* re := regexp.MustCompile(`[a-zA-Z0-9]+@[a-zA-Z0-9.]+\.[a-zA-Z0-9]+`) */

	//正則表達式提取功能
	//將所需要的部分用()框起來
	re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9]+)(\.[a-zA-Z0-9.]+)`)
	//用一個text去找有沒有這個string
	//也可以用byte去找(re.Find())
	//有找到的話返回一個match string
	// match := re.FindString(text)

	//如果要找到多行用FindAllString
	//傳入string,和要找到幾個,-1表示找到所有
	/* match := re.FindAllString(text, -1) */

	//FindAllString沒有辦法提取,
	//改用FindAllStringSubmatch
	//會返回一個二維的STRING slice
	//本身自己會佔一個,後面匹配的也會佔
	//[55151@gmail.com 55151 gmail .com]
	match := re.FindAllStringSubmatch(text, -1)
	fmt.Println(match)
	for _, m := range match {
		fmt.Println(m)
	}
}
