package http

import (
	"fmt"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

// type statusCode int

type Request struct {
	Method string
	URL    *URL
	Proto  string
	Header Header
	Body   io.Reader
}

type Header map[string]string

func (r *Request) SetHeader(key string, value string) {
	if r.Header == nil {
		r.Header = make(Header)
	}
	r.Header[key] = value
}

type Handler interface {
	Handle(ResponseWriter, *Request)
}

type ResponseWriter struct {
	statusCode Status
	conn       net.Conn
	Request    // used for headers
}

func (rw ResponseWriter) WriteHeader(key string, value string) {
	rw.Request.SetHeader(key, value)
}

func (rw *ResponseWriter) WriteStatus(status Status) {
	rw.statusCode = status
}

func (rw ResponseWriter) Write(p []byte) (int, error) {
	defer rw.conn.Close()
	rw.WriteHeader("Content-Length", fmt.Sprint(len(p)))
	resp, err := prepareResponse(rw, p)
	if err != nil {
		return 0, err
	}
	// fmt.Println(string(p))
	// fmt.Println(string(resp))
	cnt, err := rw.conn.Write(resp)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func prepareResponse(rw ResponseWriter, p []byte) ([]byte, error) {
	// body, err := io.ReadAll(rw.Request.Body)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// }
	headerString := fmt.Sprintf("%s %d %s\r\n", rw.Proto, rw.statusCode, StatusString[rw.statusCode])
	for key, val := range rw.Header {
		headerString += fmt.Sprintf("%s=%s\r\n", key, val)
	}
	headerString += "\r\n" // separator for the body
	return append([]byte(headerString), p...), nil
}

func ListenAndServe(address string, handler Handler) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("Failed to listen at %s", address)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error on Accepting connection %v\n", err)
		}
		serv := &server{conn, handler}
		fmt.Printf("serving with %v\n", serv)
		go serv.serve()
	}
}
