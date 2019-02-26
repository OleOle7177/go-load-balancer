package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	sts := parseSettings("some path will be here")

	for _, frontend := range sts.http {
		mux := http.NewServeMux()
		mux.HandleFunc(frontend.path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Handled!")
		})

		listenTo := fmt.Sprintf(":%v", frontend.port)

		wg.Add(1)
		go func() {
			if err := http.ListenAndServe(listenTo, mux); err != nil {
				log.Fatal(err)
			}
		}()
	}
	wg.Wait()
}

// // Lets run the frontned
// func (s *Server) RunFrontendServer(frontend *Frontend) {
// 	if len(frontend.Backends) == 0 {
// 		log.Fatal(errNoBackend.Error())
// 	}
//
// 	host := frontend.Host
// 	port := frontend.Port
// 	address := fmt.Sprintf("%s:%d", host, port)
//
// 	for _, backend := range frontend.Backends {
// 		// Before start the backend let's set a monitor
// 		backend.HeartCheck()
// 	}
//
// 	log.Printf("Run frontend server [%s] at [%s]", frontend.Name, address)
//
// 	// Prepare the mux
// 	httpHandle := http.NewServeMux()
//
// 	httpHandle.HandleFunc(frontend.Route, func(w http.ResponseWriter, r *http.Request) {
// 		s.Lock()
// 		s.Add(1)
// 		s.Unlock()
//
// 		// On a serious problem
// 		defer func() {
// 			if rec := recover(); rec != nil {
// 				log.Println("Err", rec)
// 				http.Error(w, http.StatusText(http.StatusInternalServerError),
// 					http.StatusInternalServerError)
// 			}
// 		}()
//
// 		// Get a channel the already attached to a worker
// 		chanResponse := s.Get(r, frontend)
// 		defer close(chanResponse)
//
// 		r.Close = true
//
// 		// Timeout ticker
// 		ticker := time.NewTicker(frontend.Timeout)
// 		defer ticker.Stop()
//
// 		select {
// 		case result := <-chanResponse:
// 			// We have a response, it's valid ?
// 			for k, vv := range result.Header {
// 				for _, v := range vv {
// 					w.Header().Set(k, v)
// 				}
// 			}
//
// 			s.Lock()
// 			s.Done()
// 			s.Unlock()
//
// 			if result.Upgraded {
// 				if s.Configuration.GeneralConfig.Websocket {
// 					result.HijackWebSocket(w, r)
// 				}
// 			} else {
// 				w.WriteHeader(result.Status)
//
// 				if r.Method != "HEAD" {
// 					w.Write(result.Body)
// 				}
// 			}
// 		case <-r.Cancel:
// 			s.Lock()
// 			s.Done()
// 			s.Unlock()
//
// 		case <-ticker.C:
// 			s.Lock()
// 			s.Done()
// 			s.Unlock()
//
// 			// Timeout
// 			http.Error(w, errTimeout.Error(), http.StatusRequestTimeout)
// 		}
// 	})
//
// 	// Config and start server
// 	server := &http.Server{
// 		Addr:    address,
// 		Handler: httpHandle,
// 	}
//
// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
