package http

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_BUFFER_SIZE = 1024
)

type server struct {
	conn    net.Conn
	handler Handler
}

// It is a wrapper around buffered Reader
// Buffer grows as data is further read
type BodyReader struct {
	reader        *bufio.Reader
	mu            sync.Mutex
	contentLength int
	buffer        []byte
	readCount     int // How much is already read
	closed        bool
}

// reads unread data into the provided slice of byte
// returns EOF if the Read() call has already return content-length number of bytes
func (cr *BodyReader) Read(b []byte) (int, error) {
	if cr.readCount == cr.contentLength {
		cr.closed = true
		return 0, io.EOF
	}

	start := cr.readCount
	end := min(len(cr.buffer), cr.contentLength)
	cr.mu.Lock()
	cnt := copy(b, cr.buffer[start:end])
	cr.mu.Unlock()
	cr.readCount += cnt
	return cnt, nil
}

func NewBodyReader(r *bufio.Reader, contentLength int) *BodyReader {
	cr := BodyReader{
		reader:        r,
		buffer:        make([]byte, DEFAULT_BUFFER_SIZE),
		contentLength: contentLength,
	}
	// no need to do anything if there is not content-length header set. It means there is no body
	// if called Read on no data, it is blocked
	if contentLength > 0 {
		// make first read synchronously
		_, err := cr.reader.Read(cr.buffer)
		if err == nil {
			go cr.backGroundRead()
		}
	}

	return &cr
}

func (cr *BodyReader) backGroundRead() {
	for {
		tempBuff := make([]byte, DEFAULT_BUFFER_SIZE)
		n, err := cr.reader.Read(tempBuff)
		if cr.closed || err != nil {
			// log.Infof("Stopped background read %v", err)
			return
		}
		cr.mu.Lock()
		cr.buffer = append(cr.buffer, tempBuff[:n]...)
		cr.mu.Unlock()

	}
}

func (s *server) serve() {
	r := bufio.NewReader(s.conn)
	req, err := parseRequest(r)
	if err != nil {
		log.Errorf("Could not parse request %v", err)
		return
	}
	var contentLengthInt int
	contentLength, ok := req.Header["Content-Length"]
	if ok {
		contentLengthInt, err = strconv.Atoi(contentLength)
		if err != nil {
			contentLengthInt = 0
		}
	} else {
		contentLengthInt = 0
	}
	req.Body = NewBodyReader(r, contentLengthInt)

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

	log.Infof("%s %s %s", method, path, proto)
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
			return nil, fmt.Errorf("Wrong header format %v", nextLine)
		}
		header[headerKeyVal[0]] = strings.TrimSpace(headerKeyVal[1])

	}

	req := &Request{
		Method: method,
		Proto:  proto,
		Header: Header(header),
		URL:    &url,
	}
	return req, nil
}
