package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type loadBalancer struct {
	http []*httpProxy
}

type httpProxy struct {
	port    int
	servers map[string]*httpServer
}

type httpServer struct {
	name     string
	client   *http.Client
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

			fmt.Println(val.backends[0].proxyTo)
			req, _ := http.NewRequest(r.Method, val.backends[0].proxyTo, r.Body)
			resp, _ := val.client.Do(req)
			body, _ := ioutil.ReadAll(resp.Body)
			w.Write(body)

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

func createHTTPClient(timeout int) *http.Client {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	return client
}
