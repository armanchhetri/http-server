package main

import (
	"flag"

	"github.com/codecrafters-io/http-server-starter-go/internal/http"
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
