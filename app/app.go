package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/encoder"
	"github.com/codecrafters-io/http-server-starter-go/internal/http"
	log "github.com/sirupsen/logrus"
)

type MyApp struct {
	FileServer FileApp
}

func (app MyApp) UserAgentHandler(rw http.ResponseWriter, r *http.Request) {
	userAgent := r.Header["User-Agent"]
	rw.Write([]byte(userAgent))
}

func (app MyApp) EchoHandler(rw http.ResponseWriter, r *http.Request) {
	contentEncoding, _ := r.Header[string(http.AcceptEncodingHeader)]
	var encodingScheme string
	var encoderInstance encoder.Encoder
	for _, encoding := range strings.Split(contentEncoding, ",") {
		encoderInstance = encoder.EncoderFactory(encoding)

		if encoderInstance != nil {
			encodingScheme = encoding
			break
		}
	}

	msg := r.PathParam["msg"]
	if encoderInstance != nil {
		dataReader, err := encoderInstance.Encode([]byte(msg))
		if err != nil {
			log.Errorf("Unable to encode data: %v", err)
			rw.WriteStatus(http.StatusInternalServerError)
			rw.Write([]byte{})
		}
		rw.SetHeader(string(http.ContentEncodingHeader), encodingScheme)
		io.Copy(rw, dataReader)
		return
	}

	rw.Write([]byte(msg))
}

func (app MyApp) Home(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello There!"))
}

func (app MyApp) FileHandler(rw http.ResponseWriter, r *http.Request) {
	filename, ok := r.PathParam["filename"]
	if !ok {
		rw.WriteStatus(http.StatusNotFound)
		rw.Write([]byte("need to provide filename"))
	}

	if r.Method == "POST" {
		if err := app.FileServer.writeFile(filename, r.Body); err != nil {
			log.Errorf("Could not write to a file: %v", err)
			rw.WriteStatus(http.StatusInternalServerError)
			rw.WriteString("Something Bad Happened to the server ;) Do not send that request")
			return
		}

		rw.WriteStatus(http.StatusCreated)
		rw.WriteString(fmt.Sprintf("Successfully written to a file %s", filename))
		return
	}

	rw.WriteHeader("Content-Type", "application/octet-stream")
	file, err := app.FileServer.getFileContent(filename)
	if err != nil {
		rw.WriteStatus(http.StatusNotFound)
		rw.WriteString(fmt.Sprintf("File not found got %v", err))
		return
	}
	rw.Write(file)
}

type FileApp struct {
	DirectoryPath string
}

func (fa FileApp) getFileContent(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fa.DirectoryPath, filename))
}

func (fa FileApp) writeFileContent(filename string, data []byte) error {
	return os.WriteFile(filepath.Join(fa.DirectoryPath, filename), data, 0o644)
}

func (fa FileApp) writeFile(filename string, r io.Reader) error {
	file, err := os.Create(filepath.Join(fa.DirectoryPath, filename))
	if err != nil {
		return err
	}
	_, err = io.Copy(file, r)
	log.Info("Written to file")
	return err
}
