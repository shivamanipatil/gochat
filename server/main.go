package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var connections []net.Conn

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
			connections = append(connections, conn)
			go handleConnection(conn)
		}
	}()
	go readInput(msg)
loop:
	for {
		select {
		case <-sigs:
			fmt.Println("Exit:")
			break loop
		case s := <-msg:
			fmt.Println("Current connections are: ", connections)
			sendToAllConnections(connections, s)
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
			if err == io.EOF {
				connections = removeConn(connections, conn)
				fmt.Println(conn, " Disconnected")
				fmt.Println("COnnections : ", connections)
				break
			}
			fmt.Println(err)
		}
		fmt.Println(string(buff[:n]))
	}
}

func sendToAllConnections(connections []net.Conn, msg string) {
	for _, conn := range connections {
		b := []byte(msg)
		n, err := conn.Write(b)
		if err != nil {
			fmt.Println("Couldnt write to : ", conn)
		}
		if n != len(b) {
			fmt.Println("Could'nt write whole message")
		}
	}
}

func removeConn(connections []net.Conn, conn net.Conn) []net.Conn {
	j := 0
	q := make([]net.Conn, len(connections))
	for _, n := range connections {
		if n != conn {
			q[j] = n
			j++
		}
	}
	q = q[:j]
	return q
}
