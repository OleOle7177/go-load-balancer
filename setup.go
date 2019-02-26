package main

type settings struct {
	http []*httpProxy
}

type httpProxy struct {
	port     int
	path     string
	backends []*httpBackend
}

type httpBackend struct {
	weight  int
	proxyTo string
}

func parseSettings(path string) *settings {
	return &settings{http: []*httpProxy{
		&httpProxy{
			port:     5000,
			path:     "/",
			backends: []*httpBackend{&httpBackend{weight: 1, proxyTo: "http://sports.ru"}},
		},
		&httpProxy{
			port:     5005,
			path:     "/s",
			backends: []*httpBackend{&httpBackend{weight: 1, proxyTo: "http://sports.ru"}},
		},
	},
	}
}

// http:
//  - port:80
//    host:"some_host"
//    backends:
//      - weight: 20
//        proxyTo: "http://sports.ru"
