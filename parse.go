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
			client:   createHTTPClient(val.Timeout),
			name:     val.Name,
			backends: createBackendHeap(val.Backends),
		}
	}

	return srvMap
}

// TODO: use real heap with push and pop
func createBackendHeap(bs []*stsHTTPBackend) []*httpBackend {
	res := []*httpBackend{}
	for _, val := range bs {
		res = append(res, &httpBackend{weight: val.Weight, proxyTo: val.ProxyTo})
	}

	return res
}

func parseConfig(data []byte) (*stsLoadBalancer, error) {
	sts := stsLoadBalancer{}
	err := yaml.Unmarshal([]byte(data), &sts)

	return &sts, err
}
