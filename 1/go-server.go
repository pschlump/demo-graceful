package main

// Simple HTTP server  - with 5 second delay

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var Port string = "localhost:8123"
var ShutdownStarted = false

// -------------------------------------------------------------------------------------------------
func respHandlerShutdown(res http.ResponseWriter, req *http.Request) {
	ShutdownStarted = true
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(res, `{"status":"success","msg":"shutdown started"}`)
}

func respHandlerStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	if ShutdownStarted {
		fmt.Fprintf(res, `{"status":"shudown-pending"}`)
	} else {
		fmt.Fprintf(res, `{"status":"success"}`)
	}
}

func respHandlerSlow(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" || req.Method == "GET" {
		req.ParseForm()
		fmt.Println("password:", req.Form.Get("password"))
		fmt.Println("Method:", req.Method)

		time.Sleep(5 * time.Second)

		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(res, `{"status":"success","msg":"slow response"}`)
	} else {
		fmt.Fprintf(res, `
<html>
<body>
This only responds to POST and GET requests.
</body>
</html>
`)
	}
}

// -------------------------------------------------------------------------------------------------
func main() {
	http.HandleFunc("/api/shutdown", respHandlerShutdown)
	http.HandleFunc("/api/status", respHandlerStatus)
	http.HandleFunc("/", respHandlerSlow)

	log.Fatal(http.ListenAndServe(Port, nil))
}
