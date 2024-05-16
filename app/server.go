package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error on Accepting connection %v\n", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buffer))
	resp := "HTTP/1.1 200 OK\r\n\r\n"
	conn.Write([]byte(resp))
}
