//把scheduler的功能實作出來
package scheduler

import "crawler/engine"

type SimpleScheduler struct {
	//放一個workerchannel,然後結構是engine.Request
	//用來放傳進來的engine.Request
	workerChan chan engine.Request
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	//要WorkerChan把s.workerChan返回
	return s.workerChan
}

//WorkerReady不做事情,但為了interface的一致,還是要實現
func (s *SimpleScheduler) WorkerReady(chan engine.Request) {

}
func (s *SimpleScheduler) Run() {
	//把workerChan做出來
	s.workerChan = make(chan engine.Request)
}

//把engine.Request的資料傳到workerchan
//要直接改變struct內的內容,因此必須是星號,變成指針類型
func (s *SimpleScheduler) Submit(r engine.Request) {
	//在這邊要用go routine給他並發執行
	//concurrent的Run會一直丟從out channel接收到的值過來
	//如果不用並發執行,整個程式會卡住
	go func() { s.workerChan <- r }()
}
