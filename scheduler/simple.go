//把scheduler的功能實作出來
package scheduler

import "crawler/engine"

type SimpleScheduler struct {
	//放一個workerchannel,然後結構是engine.Request
	//用來放傳進來的engine.Request
	workerchan chan engine.Request
}

//接收channel
//要直接改變struct內的內容,因此必須是星號,變成指針類型
func (s *SimpleScheduler) ConfigureMasterWorkerChan(c chan engine.Request) {
	s.workerchan = c
}

//把engine.Request的資料傳到workerchan
//要直接改變struct內的內容,因此必須是星號,變成指針類型
func (s *SimpleScheduler) Submit(r engine.Request) {
	s.workerchan <- r
}
