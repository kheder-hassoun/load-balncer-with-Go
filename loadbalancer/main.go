package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err.Error())
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthCheckInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	servers := initializeServers(config.Servers)
	startHealthChecks(servers, healthCheckInterval)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, servers)
	})

	log.Println("Starting server on port", config.ListenPort)
	err = http.ListenAndServe(config.ListenPort, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
