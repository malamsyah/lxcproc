package main

import (
	"encoding/json"
	"log"
	"net/http"

	lxd "github.com/lxc/lxd/client"
)

func getResourceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		response, err := json.Marshal(`{"message":"request not found"}`)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

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

func getContainersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		returnJSON(w, `{"message":"request not found"}`)
		return
	}

	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		panic(err)
	}

	containers, err := c.GetContainers()
	if err != nil {
		panic(err)
	}

	returnJSON(w, containers)
}

func main() {
	http.HandleFunc("/lxd/resource", getResourceHandler)
	http.HandleFunc("/lxd/containers", getContainersHandler)
	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func returnJSON(w http.ResponseWriter, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
