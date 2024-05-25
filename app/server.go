package main

import (

	// "log"
	// "net"
	// "os"
	// "strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
)

func main() {
	var app MyApp
	mux := http.NewMux()
	mux.Register("/user-agent", app.UserAgentHandler)
	mux.Register("/echo/<msg>", app.EchoHandler)

	http.ListenAndServe("0.0.0.0:4221", mux)
}

type MyApp struct{}

func (app MyApp) UserAgentHandler(rw http.ResponseWriter, r *http.Request) {
	userAgent := r.Header["User-Agent"]
	// rw.WriteHeader("Content-Type", "text/plain")
	// fmt.Println("got: ", userAgent)
	rw.Write([]byte(userAgent))
}

func (app MyApp) EchoHandler(rw http.ResponseWriter, r *http.Request) {
	msg := r.PathParam["msg"]

	// fmt.Println("got: ", userAgent)
	rw.Write([]byte(msg))
}
