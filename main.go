package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func getLXCInfo(name string) {
	cmd := exec.Command("lxc", "info", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func getLXCList() ([]LXCInfoJSON, error) {
	cmd := exec.Command("lxc", "list", "--format", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lxcs, err := parseFile(string(out))
	if err != nil {
		return nil, err
	}
	for _, lxc := range lxcs {
		getLXCInfo(lxc.Name)
	}
	return lxcs, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, err := getLXCList()
	if err != nil {
		panic(err)
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dataJSON)
}

func parseFile(raw string) ([]LXCInfoJSON, error) {
	test_input := []byte(raw)
	var result []LXCInfoJSON
	err := json.Unmarshal(test_input, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	http.HandleFunc("/_status", handler)
	err := http.ListenAndServe(":4041", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type LXCInfoJSON struct {
	Architecture string       `json:"architecture"`
	Ephemeral    bool         `json:"ephemeral"`
	Stateful     bool         `json:"stateful"`
	Description  string       `json:"description"`
	CreatedAt    time.Time    `json:"created_at"`
	Name         string       `json:"name"`
	Status       string       `json:"status"`
	StatusCode   int          `json:"status_code"`
	LastUsedAt   time.Time    `json:"last_used_at"`
	State        LXCStateJSON `json:"state"`
}

type LXCStateJSON struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Disk       struct {
		Root struct {
			Usage int `json:"usage"`
		} `json:"root"`
	} `json:"disk"`
	Memory struct {
		Usage         int `json:"usage"`
		UsagePeak     int `json:"usage_peak"`
		SwapUsage     int `json:"swap_usage"`
		SwapUsagePeak int `json:"swap_usage_peak"`
	} `json:"memory"`
	PID       int `json:"pid"`
	Processes int `json:"processes"`
	CPU       struct {
		Usage int `json:"usage"`
	} `json:"cpu"`
}
