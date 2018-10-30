package main

import (
	"net"
					"syscall"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"time"
	"io"
	"fmt"
)

const COUNT = 10000

var (
	connList = make([]net.Conn, COUNT)
	iosv = asio.NewIoContext()
	as = asio.NewReactiveSocket(iosv)
	index = 0
	stop = false
)

func main() {
	go doDial()

	sigset := asio.NewSignalSet(iosv, syscall.SIGINT)
	sigset.AsyncWait(func(err error) {
		iosv.Stop()
		sigset.Cancel()
	})

	iosv.Run()

	close()
}

func close() {
	stop = true
	time.Sleep(time.Second)
	for i:=0; i<=index && i<COUNT; i++ {
		connList[i].Close()
	}
}

func doDial() {
	for i:=0; i<COUNT && !stop; i++ {
		conn, err := net.Dial("tcp", ":8888")
		if err != nil {
			fmt.Println("Error net dial", err)
			return
		}

		connList[index] = conn
		index ++

		go doWrite(conn)
	}
}

func doWrite(conn io.Writer) {
	time.Sleep(time.Second)
	conn.Write([]byte("hello"))
	if !stop {
		doWrite(conn)
	}
}