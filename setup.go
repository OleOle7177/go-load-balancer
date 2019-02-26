package main

type settings struct {
	http []httpProxy
}

type httpProxy struct {
	host     string
	port     string
	backends []httpBackend
}

type httpBackend struct {
	weight  int
	proxyTo string
}

func parseSettings(path string) *settings {
	return &settings{http: []httpProxy{}}
}

// http:
//  - port:80
//    host:"some_host"
//    backends:
//      - weight: 20
//        proxyTo: "http://sports.ru"
