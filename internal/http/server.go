package http

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type server struct {
	conn    net.Conn
	handler Handler
}

func (s *server) serve() {
	r := bufio.NewReader(s.conn)
	req, err := parseRequest(r)
	if err != nil {
		log.Errorf("Could not parse request %v", err)
		return
	}

	resp := Request{
		URL:   req.URL,
		Proto: req.Proto,
	}
	rw := ResponseWriter{
		statusCode: StatusOk,
		Request:    resp,
		conn:       s.conn,
	}
	rw.WriteHeader("Content-Type", "text/plain")
	// fmt.Println(req, rw)
	s.handler.ServeHTTP(rw, req)
	err = s.conn.Close()
	if err == nil {
		log.Error("Handler Function should write something to the wire")
	}
}

func parseRequest(r *bufio.Reader) (*Request, error) {
	firstLine, isPrefix, err := r.ReadLine()
	if isPrefix || err != nil {
		return nil, errors.New("Malformed HTTP request")
	}
	httpInfo := strings.Split(string(firstLine), " ")
	if len(httpInfo) != 3 {
		return nil, errors.New("Malformed HTTP request")
	}
	method, urlStr, proto := httpInfo[0], httpInfo[1], httpInfo[2]

	pathQuery := strings.Split(urlStr, "?")

	path := pathQuery[0]
	var queryString string
	if len(pathQuery) == 2 {
		queryString = pathQuery[1]
	}

	url := URL{
		Raw:         urlStr,
		Path:        path,
		QueryString: queryString,
	}

	fmt.Println(method, path, proto)
	header := make(map[string]string)

	for {
		line, isPrefix, err := r.ReadLine()
		if isPrefix || err != nil {
			return nil, errors.New("Malformed HTTP request")
		}
		nextLine := string(line)
		if nextLine == "" {
			break
		}
		headerKeyVal := strings.SplitN(string(nextLine), ":", 2)
		if len(headerKeyVal) != 2 {
			fmt.Println(nextLine)
			return nil, fmt.Errorf("Wrong header format %v", nextLine)
		}
		header[headerKeyVal[0]] = strings.TrimSpace(headerKeyVal[1])

	}

	req := &Request{
		Method: method,
		Proto:  proto,
		Header: Header(header),
		Body:   r,
		URL:    &url,
	}
	return req, nil
}
