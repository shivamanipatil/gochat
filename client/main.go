package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	msg := make(chan string, 1)
	quitInput := make(chan bool)
	quitConn := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var name string
	fmt.Println("Enter your name: ")
	fmt.Scanf("%s", &name)

	conn, err := net.Dial("tcp", ":8880")
	if err != nil {
		fmt.Println("Couldnt connect")
		os.Exit(1)
	}

	go readInput(msg, quitInput)
	go handleConnection(conn, quitConn, sigs)

	for {
		select {
		case <-sigs:
			fmt.Println("Exit:")
			close(quitInput)
			close(quitConn)
			conn.Close()
			return
		case s := <-msg:
			s = name + " : " + s
			_, err := conn.Write([]byte(s))
			if err != nil {
				fmt.Println("Couldn't write :")
			}
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

func handleConnection(conn net.Conn, quitConn chan bool, sig chan os.Signal) {
	buff := make([]byte, 400)
	for {
		select {
		case <-quitConn:
			return
		default:
			n, err := conn.Read(buff)
			if err != nil {
				sig <- syscall.SIGINT
				return
			}
			fmt.Println(string(buff[:n]))
		}
	}
}
