package main

import (
	"fmt"
	"io/ioutil"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	yaml, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Printf("Error while reading config: %s\n", err)
	}

	sts, err := parseConfig([]byte(yaml))
	if err != nil {
		fmt.Printf("Error while parsing config: %s\n", err)
	}

	h := sts.createHTTPProxy()

	for _, frontend := range h.http {
		frontend.launchServer()
		wg.Add(1)
	}

	wg.Wait()
}
