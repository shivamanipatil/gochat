package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var connections []net.Conn

func main() {
	//channels
	sigs := make(chan os.Signal, 1)
	msg := make(chan string, 1)
	quitInput := make(chan bool)
	quitAcceptLoop := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//Listener
	ln, err := net.Listen("tcp", ":8880")
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		defer ln.Close()
		for {
			select {
			case <-quitAcceptLoop:
				fmt.Println("Out of accept loop")
				return
			default:
				conn, err := ln.Accept()
				if err != nil {
					fmt.Println(err)
				}
				connections = append(connections, conn)
				go handleConnection(conn)
			}
		}
	}()
	go readInput(msg, quitInput)

	for {
		select {
		case <-sigs:
			fmt.Println("Exit:")
			close(quitAcceptLoop)
			close(quitInput)
			return
		case s := <-msg:
			s = "Server : " + s
			sendToAllConnections(connections, s, nil)
		}
	}

}

func readInput(msg chan string, quitInput chan bool) {
	// Receive input in a loop
	reader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-quitInput:
			return
		default:
			s, _ := reader.ReadString('\n')
			msg <- s
		}
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
				return
			}
			fmt.Println(err)
		}
		fmt.Println(string(buff[:n]))
		sendToAllConnections(connections, string(buff[:n]), conn)
	}
}

func sendToAllConnections(connections []net.Conn, msg string, selfConn net.Conn) {
	for _, conn := range connections {
		b := []byte(msg)
		if conn == selfConn {
			continue
		}
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
