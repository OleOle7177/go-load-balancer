package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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
	name        string
	client      *http.Client
	pool        *httpBackendPool
	poolCounter chan (struct{})
}

type httpBackendPool struct {
	backends *backendHeap
	mux      sync.Mutex
}

type httpBackend struct {
	weight  int
	proxyTo string
}

func (h *httpProxy) launchServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if val, ok := h.servers[r.Host]; ok {
			// check if there are available connections in pool
			val.poolCounter <- struct{}{}
			defer func() { <-val.poolCounter }()

			// lock backend heap
			val.pool.mux.Lock()
			defer val.pool.mux.Unlock()
			b, _ := val.pool.backends.Pop()
			defer val.pool.backends.Push(b)

			fmt.Println(b.proxyTo)
			req, _ := http.NewRequest(r.Method, b.proxyTo, r.Body)
			req.Host = r.Host
			req.Header.Set("User-Agent", r.UserAgent())
			req.Header.Set("X-Forwarded-For", r.Referer())

			resp, err := val.client.Do(req)
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
			}
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				http.Error(w, "Internal Server Error", 500)
			}

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
