package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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

// okHeader := "HTTP/1.1 200 OK\r\nContent-Type: text/plain"

func handleConn(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	request := string(buffer)
	fmt.Println(request)

	respOk := "HTTP/1.1 200 OK\r\n\r\n"
	respNotFound := "HTTP/1.1 404 Not Found\r\n\r\n"
	method, path, proto := parseHeader(request)
	fmt.Println(method, path, proto)
	if path == "/" {
		fmt.Fprint(conn, respOk)
	} else if strings.HasPrefix(path, "/echo") {
		echoResponse(conn, path)
	} else {
		fmt.Fprint(conn, respNotFound)
	}
}

func echoResponse(conn net.Conn, path string) {
	message := strings.Split(path, "/")[2]
	fmt.Println(message)
	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
	fmt.Fprint(conn, resp)
}

func parseHeader(request string) (string, string, string) {
	header1 := strings.Split(request, "\r\n")[0]
	// fmt.Println("headers", header1)
	firstHeaders := strings.Split(header1, " ")
	return firstHeaders[0], firstHeaders[1], firstHeaders[2]
}
