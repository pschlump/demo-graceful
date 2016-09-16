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
//	 	ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7s u2A+gf7Q==
// To submit your work please invite ppg to a public GitHub project, and be sure to document any configuration or run instructions.
//

// TODO
//	1. Pull out simple server that has a 1 sec delay - use that for the 5 second delay
//		1. Remember to change delay to 5 sec.
//	2. Check "curl" results in a "POST"
//  3. Configuration option to set host:port
//	4. Build StringHash512 - that uses SHA512 - test it with this value
//	5. Pull base64 encode from the AES/SRP stuff
//
// TODO Tests
//	1. Check on multiple connections simultaneously
//	2. Graceful shutdown
// 		curl —data "password=angryMonkey" http://localhost:8080/api/graceful-shutdown
//		Also allow for catching of a "signal"
//	3. Testing of shutdown process
//
// Code breakdown
//	Component							Time Est			Test Time Est		Actual			Test
//	----------------					---------			-------------		------			----
//		./HashString						15min			25min				8min			5min
//		./Encode64							15min			25min
//		./ReadCfg							20min			45min
//		./Graceful						2hrs			4hrs
//	main.go									30min			30min
//	Makefile - with examples and tests					2hrs
//	Documentation - 					1hrs
//		Edit 							1hrs
//
// Estimate Total Project Time: Approx: 9-14hrs
// Actual Total Project Time: 10.2 hrs
//

func main() {
}
