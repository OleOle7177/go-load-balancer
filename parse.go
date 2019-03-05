package main

import (
	yaml "gopkg.in/yaml.v2"
)

type stsLoadBalancer struct {
	HTTP []*stsHTTPPort `yaml:"http"`
}

type stsHTTPPort struct {
	Port    int              `yaml:"port"`
	Servers []*stsHTTPServer `yaml:"servers"`
}

type stsHTTPServer struct {
	Name     string            `yaml:"name"`
	Host     string            `yaml:"host"`
	Timeout  int               `yaml:"timeout"`
	Backends []*stsHTTPBackend `yaml:"backends"`
}

type stsHTTPBackend struct {
	Weight  int    `yaml:"weight"`
	ProxyTo string `yaml:"proxyTo"`
}

func (s *stsLoadBalancer) createHTTPProxy() *loadBalancer {
	res := []*httpProxy{}
	for _, p := range s.HTTP {
		res = append(res, &httpProxy{
			port:    p.Port,
			servers: createServersMap(p.Servers),
		})
	}

	return &loadBalancer{http: res}
}

func createServersMap(sts []*stsHTTPServer) map[string]*httpServer {
	srvMap := make(map[string]*httpServer)
	for _, val := range sts {
		srvMap[val.Host] = &httpServer{
			client: createHTTPClient(val.Timeout),
			name:   val.Name,
			pool:   createBackendHeap(val.Backends),
		}
	}

	return srvMap
}

func createBackendHeap(bs []*stsHTTPBackend) *httpBackendPool {
	heap := new(backendHeap)
	for _, b := range bs {
		heap.Push(&httpBackend{weight: b.Weight, proxyTo: b.ProxyTo})
	}
	return &httpBackendPool{backends: heap}
}

func parseConfig(data []byte) (*stsLoadBalancer, error) {
	sts := stsLoadBalancer{}
	err := yaml.Unmarshal([]byte(data), &sts)

	return &sts, err
}
