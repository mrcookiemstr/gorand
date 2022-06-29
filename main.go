package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultHostName     = "localhost"
	envHttpPortName     = "GORAND_HTTP_PORT"
	envRandomOrgApiKey  = "GORAND_API_KEY"
	defaultHttpPort     = "8080"
	defaultRandomOrgKey = "665740fa-1cad-4fa5-8bdf-9bbfb2a25b7e"
)

func ProvideHostName() string {
	osHostName, osHostNameErr := os.Hostname()

	if osHostNameErr == nil {
		return defaultHostName
	}

	return osHostName
}

func ProvideHttpRequestPort() string {
	httpPort := os.Getenv(envHttpPortName)

	if httpPort == "" {
		return defaultHttpPort
	}

	return httpPort
}

func ProvideRandomOrgApiKey() string {
	randomOrgKey := os.Getenv(envRandomOrgApiKey)

	if randomOrgKey == "" {
		panic("Random.org key not present")
	}

	return randomOrgKey
}

func main() {

	hostName, httpPort := ProvideHostName(), ProvideHttpRequestPort()
	log.Println("Http server at http://" + hostName + ":" + httpPort)

	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%s", httpPort),
		WriteTimeout: 5 * time.Second,
		Handler:      http.TimeoutHandler(http.HandlerFunc(HandleHttpRequest), time.Second, "Timeout"),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("Http server error ", err)
	}
}
