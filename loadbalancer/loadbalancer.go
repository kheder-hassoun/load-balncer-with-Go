package main

import (
	"log"
	"net/http"
	"time"
)

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

func handleRequest(w http.ResponseWriter, r *http.Request, servers []*Server) {
	log.Println("Received request:", r.Method, r.URL.Path, "Server:", r.Host)
	server := nextServerLeastActive(servers)
	server.Mutex.Lock()
	server.ActiveConnections++
	server.Mutex.Unlock()
	server.Proxy().ServeHTTP(w, r)
	server.Mutex.Lock()
	time.Sleep(3 * time.Second)
	server.ActiveConnections--
	server.Mutex.Unlock()

	log.Println("Request completed:", r.Method, r.URL.Path, "Server:", r.Host, "Server URL:", server.URL)
}
