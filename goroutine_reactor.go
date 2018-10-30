package asio

import (
	"reflect"
	)

type GoroutineReactor struct {
	opq 	 chan operation
	//notifies chan os.Signal
	down chan struct{}
}

type operation struct {
	function  interface{}
	argument  []interface{}
}

func NewGoroutineReactor() *GoroutineReactor {
	r := new(GoroutineReactor)
	r.opq = make(chan operation, 128)
	r.down = make(chan struct{}, 1)
	//r.notifies = make(chan os.Signal, 1)
	return r
}

// use GoroutineReactor for each operating system
func (i *IoContext) GetService () ReactorService {
	if i.service == nil {
		i.service = NewGoroutineReactor()
	}
	return i.service
}

func (g *GoroutineReactor) run() {
	for ;; {
		select {
		case <-g.down:
			return
		case op := <-g.opq:
			g.doReactor(op.function, op.argument)
			break
		}
	}
}

func (g *GoroutineReactor) stop() {
	g.down <- struct{}{}
}

func (g *GoroutineReactor) post(op interface{}, args ...interface{}) {
	g.opq <- operation{op, args}
}

//func (g *GoroutineReactor) notify(sig ...os.Signal) {
	//signal.Notify(g.notifies, sig...)
//}

func (g *GoroutineReactor) doReactor(op interface{}, args []interface{}) {
	opv := reflect.ValueOf(op)
	opt := reflect.TypeOf(op)

	if opt == nil {
		println("invalid operation <nil>")
	}

	if opt.Kind() != reflect.Func {
		println("op must be a callback function")
		return
	}

	opNum := opt.NumIn()
	if opNum != len(args) {
		println("invalid arguments", "arguments needs:", opNum)
		return
	}

	opArgs := make([]reflect.Value, opNum)

	for i:=0; i<opt.NumIn(); i++ {
		if args[i] == nil {
			opArgs[i] = reflect.Zero(opt.In(i))
			continue
		}

		if !reflect.TypeOf(args[i]).AssignableTo(opt.In(i)) {
			println("invalid arguments", "wrong args#", i)
			return
		}

		opArgs[i] = reflect.ValueOf(args[i])
	}

	//fmt.Println("args", opArgs)

	opv.Call(opArgs)
}

