package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	token := "s.P9iXYsw8m6B6R7gQ0w6eTngG"
	start := time.Now()
	stop := start.Add(time.Minute * 30)
	for {
		if time.Now().After(stop) {
			break
		}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8200/v1/barebones/empty-call", nil)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		req.Header.Set("X-Vault-Token", token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("%s\n", b)
	}
}

