package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"jenkins-plugins/cache"
	"jenkins-plugins/parse"
)

func main() {
	port := 3000
	mux := http.NewServeMux()
	cache := cache.New()
	handler := NewHandler(&cache)

	mux.HandleFunc("/jenkins/updates/current/update-center.json", handler.DummyUpdateCenterFile)
	mux.HandleFunc("/ParseAwayJSONP", parse.WithoutJSONP)
	mux.HandleFunc("/Update", handler.UpdateCache)
	mux.HandleFunc("/Plugins", handler.Plugins)
	mux.HandleFunc("/", handler.Index)

	go func() {
		time.Sleep(time.Second * 2)
		log.Println("Performing ache update.")
		handler.PerformUpdateCache()
	}()
	log.Println("Server to listen on a port: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))

}
