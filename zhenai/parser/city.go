package parser

import (
	"crawler/engine"
	"regexp"
)

//把City的url做常量使用
var (
	profileRe = regexp.MustCompile(`<a href="(http://localhost:8080/mock/album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe = regexp.MustCompile(`href="(http://localhost:8080/mock/www.zhenai.com/zhenghun/guangzhou/[^"]+)"`)
)

//返回package engine內的ParseResult types
func ParseCity(contents []byte) engine.ParserResult {
	//使用正則表達式把城市名稱取出
	//正則表達式提取功能
	//將所需要的部分用()框起來
	//re := regexp.MustCompile(cityRe)

	//FindAllString沒有辦法提取,
	//改用FindAllStringSubmatch
	//會返回一個二維的STRING slice
	//每個匹配都佔一個項
	//本身自己會佔一個,後面匹配的也會佔
	//[55151@gmail.com 55151 gmail .com]
	matches := profileRe.FindAllSubmatch(contents, -1)
	result := engine.ParserResult{}
	//把返回的二維打印出來
	//完整的如下
	//[<a href="http://www.zhenai.com/zhenghun/zunyi" data-v-2cb5b6a2>遵义</a> http://www.zhenai.com/zhenghun/zunyi 遵义]

	//m[0]=<a href="http://www.zhenai.com/zhenghun/zunyi" data-v-2cb5b6a2>遵义</a>
	//m[1]=http://www.zhenai.com/zhenghun/zunyi
	//m[2]=遵义
	for _, m := range matches {
		name := string(m[2])
		//把城市的名字用append做為一個items返回出去,把原本的值換成string丟出去
		result.Items = append(result.Items, "User"+name)
		//把URL用append存到Result中返回
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			//使用一個匿名函數,調用ParseProfile,並把抓到的name傳給ParseProfile
			//要不然沒辦法把這邊取到的name傳給固定的ParserFunc形式
			ParserFunc: func(c []byte) engine.ParserResult {
				return ParseProfile(c, name)
			},
		})
	}
	//取第二頁的URL
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		//把URL用append存到Result中返回
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			//將取到的城市url傳到ParseCity func做解析
			ParserFunc: ParseCity,
		})
	}

	return result
}
