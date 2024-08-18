package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

func main() {
	jar, _ := cookiejar.New(nil)
	http.DefaultClient.Jar = jar
	r := bytes.NewBuffer(nil)
	json.NewEncoder(r).Encode(map[string]any{
		"email":    "w.matt.pavin@gmail.com",
		"password": "abcd1234",
	})
	res, err := http.Post(buildPath("/api/users/signin"), "application/json", r)
	if err != nil {
		fmt.Println("error auth:", err)
	}
	var authR map[string]any
	json.NewDecoder(res.Body).Decode(&authR)
	wg := new(sync.WaitGroup)
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createOrder()
		}()
	}
	wg.Wait()
}

func buildPath(path string) string {
	return "http://ticketing.dev" + path
}

func createOrder() {
	r := bytes.NewBuffer(nil)
	json.NewEncoder(r).Encode(map[string]any{
		"title": "abc",
		"price": 5,
	})
	res, _ := http.Post(buildPath("/api/tickets"), "application/json", r)
	var ticket map[string]any
	json.NewDecoder(res.Body).Decode(&ticket)
	ticketId := ticket["id"].(string)

	r = bytes.NewBuffer(nil)
	json.NewEncoder(r).Encode(map[string]any{
		"title": ticket["title"],
		"price": 10,
	})
	req, _ := http.NewRequest(http.MethodPut, buildPath("/api/tickets/"+ticketId), r)
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)

	r = bytes.NewBuffer(nil)
	json.NewEncoder(r).Encode(map[string]any{
		"title": ticket["title"],
		"price": 15,
	})
	req, _ = http.NewRequest(http.MethodPut, buildPath("/api/tickets/"+ticketId), r)
	req.Header.Set("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}
