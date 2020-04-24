package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func sum1(s1 chan int) {
	sum := 0
	for i := 0; i < 100000; i++ {
		sum = sum + i
	}
	s1 <- sum
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	msg := make(chan string, 1)

	ln, err := net.Listen("tcp", ":8880")
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
			}
			go handleConnection(conn)
		}
	}()

	s1 := make(chan int, 1)
	go readInput(msg)
	go sum1(s1)
loop:
	for {
		select {
		case <-sigs:
			fmt.Println("Exit:")
			break loop
		case s := <-msg:
			fmt.Println(s)
		case i := <-s1:
			fmt.Println(i)
		}

	}
}

func readInput(msg chan string) {
	// Receive input in a loop
	for {
		var s string
		fmt.Scan(&s)
		// Send what we read over the channel
		msg <- s
	}
}

func handleConnection(conn net.Conn) {
	buff := make([]byte, 400)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("THis is the text : ", string(buff[:n]))
	}
}
