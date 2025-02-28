package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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

func (s *Server) Proxy() *httputil.ReverseProxy {
	urlObj, _ := url.Parse(s.URL)
	return httputil.NewSingleHostReverseProxy(urlObj)
}

func initializeServers(serverConfigs []Server) []*Server {
	var servers []*Server
	for _, server := range serverConfigs {
		servers = append(servers, &Server{
			Name:              server.Name,
			URL:               server.URL,
			ActiveConnections: 0,
			Mutex:             sync.Mutex{},
			Healthy:           true,
		})
	}
	return servers
}

func startHealthChecks(servers []*Server, interval time.Duration) {
	for _, server := range servers {
		go func(s *Server) {
			for range time.Tick(interval) {
				res, err := http.Get(s.URL)
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}
}
