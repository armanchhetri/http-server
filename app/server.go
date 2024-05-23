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
	mux.HandleFunc("/user-agent", app.UserAgentHandler)

	http.ListenAndServe("0.0.0.0:4221", mux)
}

type MyApp struct{}

func (app MyApp) UserAgentHandler(rw http.ResponseWriter, r *http.Request) {
	userAgent := r.Header["User-Agent"]
	rw.WriteHeader("Content-Type", "text/plain")
	// fmt.Println("got: ", userAgent)
	rw.Write([]byte(userAgent))
}
