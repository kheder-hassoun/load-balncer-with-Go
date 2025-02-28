package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"
)

type Server struct {
	Name              string
	URL               string
	ActiveConnections int
	Mutex             sync.Mutex
	Healthy           bool
}

type Config struct {
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []Server `json:"servers"`
	ListenPort          string   `json:"listenPort"`
}

func loadConfig(file string) (Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func nextServerLeastActive(servers []*Server) *Server {
	leastActiveConnections := -1
	leastActiveServer := servers[0]
	for _, server := range servers {
		server.Mutex.Lock()
		if (server.ActiveConnections < leastActiveConnections || leastActiveConnections == -1) && server.Healthy {
			leastActiveConnections = server.ActiveConnections
			leastActiveServer = server
		}
		server.Mutex.Unlock()
	}

	return leastActiveServer
}

func (s *Server) Proxy() *httputil.ReverseProxy {
	urlObj, _ := url.Parse(s.URL)
	return httputil.NewSingleHostReverseProxy(urlObj)
}

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

	var servers []*Server
	for _, server := range config.Servers {
		servers = append(servers, &Server{
			Name:              server.Name,
			URL:               server.URL,
			ActiveConnections: 0,
			Mutex:             sync.Mutex{},
			Healthy:           true,
		})
	}

	for _, server := range servers {
		go func(s *Server) {
			for range time.Tick(healthCheckInterval) {
				res, err := http.Get(s.URL)
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:", r.Method, r.URL.Path, "Server:", r.Host)
		server := nextServerLeastActive(servers)
		server.Mutex.Lock()
		server.ActiveConnections++
		server.Mutex.Unlock()
		server.Proxy().ServeHTTP(w, r)
		server.Mutex.Lock()
		time.Sleep(15 * time.Second)
		server.ActiveConnections--
		server.Mutex.Unlock()

		log.Println("Request completed:", r.Method, r.URL.Path, "Server:", r.Host, "Server URL:", server.URL)
	})

	log.Println("Starting server on port", config.ListenPort)
	err = http.ListenAndServe(config.ListenPort, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
