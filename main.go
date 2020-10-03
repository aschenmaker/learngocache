package main

import (
	"flag"
	"fmt"
	"goCache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":           "100",
	"Mark":          "120",
	"Geek":          "180",
	"ddododododood": "0000",
}

func createGroup() *goCache.Group {
	return goCache.NewGroup("scores", 2<<10, goCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key: ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheSever(addr string, addrs []string, group *goCache.Group) {
	peers := goCache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("goCache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPISever(apiAddr string, group *goCache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/octet-stream")
		_, _ = w.Write(view.ByteSlice())
	}))
	log.Println("end server is running at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool

	flag.IntVar(&port, "port", 8801, "goCache sever port")
	flag.BoolVar(&api, "api", false, "start api sever")
	flag.Parse()

	apiAddr := "http://localhost:9999"

	addrMap := map[int]string{
		8801: "http://localhost:8801",
		8802: "http://localhost:8802",
		8803: "http://localhost:8803",
	}

	var addrs []string

	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	group := createGroup()
	if api {
		go startAPISever(apiAddr, group)
	}

	startCacheSever(addrMap[port], []string(addrs), group)
}
