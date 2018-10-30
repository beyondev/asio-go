package main

import (
	"fmt"
	"time"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"syscall"
	"net"
)

var connects []net.Conn

func main() {
	iosv := asio.NewIoContext()

	timer := asio.NewDeadlineTimer(iosv)

	scheduleLoop1(timer)

	timer2 := asio.NewDeadlineTimer(iosv)

	scheduleLoop2(timer2)

	socket := asio.NewReactiveSocket(iosv)

	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println(err)
	}

	startAcceptLoop(socket, listen)

	sigintSet := asio.NewSignalSet(iosv, syscall.SIGINT)
	sigintSet.AsyncWait(func(ec error) {
		iosv.Stop()
		sigintSet.Cancel()
	})

	iosv.Run()

	shutdown()
}

func shutdown() {
	for _,c := range connects {
		c.Close()
	}
}

func startAcceptLoop(socket *asio.ReactiveSocket, listen net.Listener) {
	socket.AsyncAccept(listen, func(conn net.Conn, err error) {
		//defer conn.Close()

		if conn == nil {
			fmt.Println("Error connect, nil")
			startAcceptLoop(socket, listen)
			return
		}

		if err != nil {
			conn.Close()
			fmt.Println("Error connect", err)
			startAcceptLoop(socket, listen)
			return
		}

		fmt.Println(conn.RemoteAddr().String())

		//fmt.Println("3sec for io operation...")
		//time.Sleep(time.Second * 3)
		//conn.Close()
		connects = append(connects, conn)

		startRead(socket, conn)
		//startWrite(socket, conn)

		startAcceptLoop(socket, listen)
	})
}

func scheduleLoop1(timer *asio.DeadlineTimer) {
	timer.Cancel()
	timer.ExpiresAt(time.Now().Add(time.Second))
	timer.AsyncWait(func(err error) {
		if err != nil {
			scheduleLoop1(timer)
			return
		}
		fmt.Println("loop 1", time.Now())
		// do something ...
		for i:=0; i<3; i++ {
			fmt.Println("do operation", i)
			time.Sleep(time.Second)
		}
		scheduleLoop1(timer)
	})
}

func scheduleLoop2(timer *asio.DeadlineTimer) {
	timer.Cancel()
	timer.ExpiresAt(time.Now().Add(time.Millisecond * 500))
	timer.AsyncWait(func(err error) {
		if err != nil {
			scheduleLoop2(timer)
			return
		}
		fmt.Println("loop 2", time.Now())
		scheduleLoop2(timer)
	})
}

func startRead(socket *asio.ReactiveSocket, conn net.Conn) {
	buf := make([]byte, 64)
	socket.AsyncRead(conn, buf, func(n int, err error) {
		if err != nil {
			// fmt.Println("Error read", ec.Error)
			conn.Close()
			return
		}
		if n > 0 {
			msg := string(buf[:n])
			fmt.Println(msg)
			socket.AsyncWrite(conn, []byte("hi"), func(n int, err error) {
				if err != nil {
					conn.Close()
				}
			})
		}
		startRead(socket, conn)
	})
}
