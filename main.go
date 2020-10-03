package main

import (
	"fmt"
	"goCache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "100",
	"Mark": "120",
	"Geek": "180",
}

func main() {
	goCache.NewGroup("scores", 2<<10, goCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key: ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	handlers := goCache.NewHTTPPool(addr)
	log.Println("Gocache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, handlers))
}
