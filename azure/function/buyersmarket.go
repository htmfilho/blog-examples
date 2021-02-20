package main

import (
	"buyersmarket/offer"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Entry point
func main() {
	httpPort := getHTTPPort()
	runHTTPServer(httpPort)
}

// We could fix the port number, but cloud environments normally require
// some flexibility on defining the server port. This is how it would work
// in Azure.
func getHTTPPort() int {
	httpPort := 8080
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		httpPort, err := strconv.Atoi(val)
		if err == nil {
			return httpPort
		}
	}
	return httpPort
}

// Calls the handlers and starts the HTTP server
func runHTTPServer(httpPort int) {
	http.HandleFunc("/api/offer", offer.GetOffer)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}
