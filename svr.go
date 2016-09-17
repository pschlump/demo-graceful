package main

//
// Assignment
// 	● The project should be written in Go using only the standard library.
// 	● When launched it should monitor a given port and wait for http connections.
// 	● When a connection is made it should not respond immediately, but rather hold the socket open for 5 seconds and then respond.
// 	● The connection will be a POST request containing a value that the client wishes to be hashed.
// 	● The hashing algorithm should be SHA512.
// 	● The result should be returned base64 encoded.
// 	● The software should be able to process multiple connections simultaneously.
// 	● If, instead of a password request the software receives a "graceful shutdown" request, it should allow any remaining password requests to complete, reject any new requests, and shutdown.
// 	● No additional password requests should be allowed when shutdown is pending.
//
// For instance, given the request (generated by curl):
// 		curl —data "password=angryMonkey" http://localhost:8080
// Your program should return
//	 	ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
// To submit your work please invite ppg to a public GitHub project, and be sure to document any configuration or run instructions.
//

//
// TODO
//	1. Pull out simple server that has a 1 sec delay - use that for the 5 second delay	-- done --
//		1. Remember to change delay to 5 sec. 											-- done --
//	2. Check "curl" results in a "POST" 												-- done --
//  3. Configuration option to set host:port 											-- done --
//	4. Build StringHash512 - that uses SHA512 - test it with this value					-- done --
//	5. Pull base64 encode from the AES/SRP stuff										-- done --
//
// TODO Tests
//	1. Check on multiple connections simultaneously	(See:test3)							-- done --
//	2. Graceful shutdown (See: test1, test2)											-- done --
//	3. Test correct return results (See: test4)											-- done --
// 		curl —data "password=angryMonkey" http://localhost:8080/api/graceful-shutdown
//	4. Also allow for catching of a "signal", $ kill -HUP 22182, (See: test5)			-- done --
//	5. Testing of shutdown process	(See:test3)											-- done --
//	6. Pull in libary and create ./godebug with colors, windows, LF						-- done --
//	7. CLI options with -D for debug													-- done --
//
// Code breakdown
//	Component							Time Est			Test Time Est		Actual			Actual Test
//	----------------					---------			-------------		------			-----------
//		./HashString						15min			25min				8min			5min
//		./ReadCfg							20min			45min				10min			9min
//		./Graceful						2hrs			4hrs					2:29min
//	main.go									30min			30min				22min			44min
//
//	Makefile - with examples and tests					2hrs					2:35						<<Note:this includes a bunch of testing>>
//	Documentation - 					1hrs											-- done --
//		Edit 							1hrs											-- done --
//
// ===========================================================================================================
//	Sums								4:05          	6:40                    5:39            57			+= 6:35
//
// Estimate Total Project Time: Approx: 7-10hrs
// Actual Total Project Time: 6.5 hrs
//

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pschlump/demo-graceful/Graceful"
	"github.com/pschlump/demo-graceful/HashString"
	"github.com/pschlump/demo-graceful/ReadCfg"
	"github.com/pschlump/demo-graceful/godebug"
)

func SetHeadersForJSON(www http.ResponseWriter, req *http.Request) {
	www.Header().Set("Content-Type", "application/json")
	SetHeadersNoCache(www, req)
}

func SetHeadersNoCache(www http.ResponseWriter, req *http.Request) {
	www.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	www.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	www.Header().Set("Expires", "0")                                         // Proxies.
}

func createRespHandlerShutdown(wg *WithGrace.WithGrace) func(www http.ResponseWriter, req *http.Request) {
	return func(www http.ResponseWriter, req *http.Request) {
		SetHeadersForJSON(www, req)
		fmt.Fprintf(www, `{"status":"success","msg":"shutdown started"}`)
		wg.GracefulShutdownServer()
	}
}

func respHandlerStatus(www http.ResponseWriter, req *http.Request) {
	SetHeadersForJSON(www, req)
	fmt.Fprintf(www, `{"status":"success"}`)
}

func createRespHandlerSlow(SleepSeconds time.Duration, wg *WithGrace.WithGrace) func(www http.ResponseWriter, req *http.Request) {
	return func(www http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" || req.Method == "GET" {
			req.ParseForm()
			Password := req.Form.Get("password")
			Fmt := req.Form.Get("fmt")

			if godebug.DebugOn("db1") {
				fmt.Println("password:", req.Form.Get("password"))
				fmt.Println("Method:", req.Method)
			}

			if Password != "" {

				time.Sleep(SleepSeconds * time.Second)

				hh := HashStrings512.HashByte512([]byte(Password))
				ee := base64.StdEncoding.EncodeToString(hh)

				if Fmt == "JSON" {
					SetHeadersForJSON(www, req)
					fmt.Fprintf(www, `{"status":"success","msg":"slow response","encoded":%q}`, ee)
				} else {
					SetHeadersNoCache(www, req)
					fmt.Fprintf(www, "%s", ee)
				}

				return
			}

			http.Error(www, "Password is a required, non-empty parameter", http.StatusBadRequest)
		}
		http.Error(www, "This only responds to GET and POST requests.", http.StatusMethodNotAllowed)
	}
}

var Cfg = flag.String("cfg", "./cfg.json", "Configuration File, default './cfg.json'")
var Debug = flag.String("debug", "", "Debug Flags")
var Version = flag.Bool("version", false, "Show the version")
var ThisOne = flag.Bool("this-one", false, "ignored flag - used for testing")

func init() {
	flag.StringVar(Cfg, "c", "./cfg.json", "Configuration File, default './cfg.json'")
	flag.StringVar(Debug, "D", "", "Debug Flags")
	flag.BoolVar(Version, "v", false, "Show the version")
}

// -------------------------------------------------------------------------------------------------
func main() {

	flag.Parse()

	if *Version {
		fmt.Printf("Version: 0.1.9\n")
		os.Exit(0)
	}

	cfg := ReadCfg.ReadCfg(*Cfg)
	godebug.SetDebugFlags(*Debug)

	// func NewWithGraceListener(netName, laddr string) (rv *WithGrace, err error) {
	wg, err := WithGrace.NewWithGraceListener("tcp", cfg.HostPort)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/graceful_shutdown", createRespHandlerShutdown(wg))
	http.HandleFunc("/api/shutdown", createRespHandlerShutdown(wg))
	http.HandleFunc("/api/status", respHandlerStatus)
	http.HandleFunc("/", createRespHandlerSlow(cfg.SleepTime, wg))

	// log.Fatal(http.ListenAndServe(cfg.HostPort, nil))
	err = wg.ListenAndServeGracefully()
	if err != nil {
		fmt.Printf("Message: %s\n", err)
	}
	wg.WaitForTheEnd()
}

/* vim: set noai ts=4 sw=4: */
