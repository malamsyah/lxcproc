package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	lxd "github.com/lxc/lxd/client"
)

func getResourceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		returnJSON(w, map[string]string{"message": "request not found"}, http.StatusBadRequest)
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

	returnJSON(w, resource, http.StatusOK)
}

func getContainersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		returnJSON(w, map[string]string{"message": "request not found"}, http.StatusBadRequest)
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

	returnJSON(w, containers, http.StatusOK)
}

func getContainerStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		returnJSON(w, map[string]string{"message": "request not found"}, http.StatusBadRequest)
		return
	}

	names, ok := r.URL.Query()["name"]

	if !ok || len(names[0]) < 1 {
		returnJSON(w, map[string]string{"message": "url param 'name' is missing"}, http.StatusBadRequest)
		return
	}

	name := names[0]

	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		panic(err)
	}

	container, _, err := c.GetContainerState(name)
	if err != nil {
		returnJSON(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	returnJSON(w, container, http.StatusOK)
}

func getMetricsHandler(w http.ResponseWriter, r *http.Request) {

	c, err := lxd.ConnectLXDUnix("", nil)
	if err != nil {
		panic(err)
	}

	resource, err := c.GetServerResources()
	if err != nil {
		panic(err)
	}

	var data string

	data += "lxd_resource_memory_total " + fmt.Sprint(resource.Memory.Total) + "\n"
	data += "lxd_resource_memory_used " + fmt.Sprint(resource.Memory.Used) + "\n"

	returnText(w, data, http.StatusOK)
}

func main() {
	http.HandleFunc("/metrics", getMetricsHandler)
	http.HandleFunc("/lxd/resource", getResourceHandler)
	http.HandleFunc("/lxd/containers", getContainersHandler)
	http.HandleFunc("/lxd/state", getContainerStateHandler)
	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func returnJSON(w http.ResponseWriter, data interface{}, status int) {
	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func returnText(w http.ResponseWriter, data string, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(data))
}
