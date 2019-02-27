package main

import (
	"fmt"
	"log"
	"net/http"
)

type settings struct {
	http []*httpProxy
}

type httpProxy struct {
	port    int
	servers map[string]*httpServer
}

type httpServer struct {
	name     string
	backends []*httpBackend
}

type httpBackend struct {
	weight  int
	proxyTo string
}

func (h *httpProxy) launchServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if val, ok := h.servers[r.Host]; ok {
			fmt.Printf("should proxy to %s\n", val.name)

		} else {
			fmt.Println("No proxy for this host header")
		}

	})

	listenTo := fmt.Sprintf(":%v", h.port)

	go func() {
		if err := http.ListenAndServe(listenTo, mux); err != nil {
			log.Fatal(err)
		}
	}()
}

func parseSettings(path string) *settings {
	return &settings{http: []*httpProxy{
		&httpProxy{
			port: 5000,
			servers: map[string]*httpServer{
				"sports.ru": &httpServer{
					name:     "sports",
					backends: []*httpBackend{&httpBackend{weight: 1, proxyTo: "http://sports.ru"}},
				},
				"vc.ru": &httpServer{
					name:     "vc",
					backends: []*httpBackend{&httpBackend{weight: 1, proxyTo: "http://vc.ru"}},
				},
			}},
		&httpProxy{
			port: 5005,
			servers: map[string]*httpServer{
				"yandex.ru": &httpServer{
					name:     "yandex",
					backends: []*httpBackend{&httpBackend{weight: 1, proxyTo: "http://yandex.ru"}},
				},
			},
		},
	},
	}
}
