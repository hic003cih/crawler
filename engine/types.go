package engine

type Request struct {
	//Url要大寫開頭因為要給別人用
	Url string
	//做一個函數,輸入是一個[]byte,輸出是一個ParserResult
	ParserFunc func([]byte) ParserResult
}

type ParserResult struct {
	Requests []Request
	//interface{}表示任何類型都可以傳入,都可以當作Items
	Items []interface{}
}

//暫時任命一個不做事情的Parser,用來返回一個空的Result
func NilParser([]byte) ParserResult {
	return ParserResult{}
}
