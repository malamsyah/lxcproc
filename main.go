package main

import (
	"encoding/json"
	"log"
	"net/http"

	lxd "github.com/lxc/lxd/client"
)

func handler(w http.ResponseWriter, r *http.Request) {
	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		panic(err)
	}

	resource, err := c.GetServerResources()
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(resource)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {

	http.HandleFunc("/lxd/status", handler)
	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
