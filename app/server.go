package main

import (

	// "log"
	// "net"
	// "os"
	// "strings"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
	log "github.com/sirupsen/logrus"
)

func main() {
	directory := flag.String("directory", "", "Supply directory name for the File server")
	flag.Parse()
	var app MyApp

	mux := http.NewMux()
	mux.Register("/", app.Home)
	mux.Register("/user-agent", app.UserAgentHandler)
	mux.Register("/echo/<msg>", app.EchoHandler)

	if *directory != "" {
		app.FileServer = FileApp{*directory}
		mux.Register("/files/<filename>", app.FileHandler)
	}

	http.ListenAndServe("0.0.0.0:4221", mux)
}

type FileApp struct {
	DirectoryPath string
}

func (fa FileApp) getFileContent(filename string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fa.DirectoryPath, filename))
}

type MyApp struct {
	FileServer FileApp
}

func (app MyApp) UserAgentHandler(rw http.ResponseWriter, r *http.Request) {
	userAgent := r.Header["User-Agent"]
	rw.Write([]byte(userAgent))
}

func (app MyApp) EchoHandler(rw http.ResponseWriter, r *http.Request) {
	log.Info(r.Header)
	msg := r.PathParam["msg"]
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
	rw.WriteHeader("Content-Type", "application/octet-stream")
	file, err := app.FileServer.getFileContent(filename)
	if err != nil {
		rw.WriteStatus(http.StatusNotFound)
		rw.WriteString(fmt.Sprintf("File not found got %v", err))
	}
	rw.Write(file)
}
