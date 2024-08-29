package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

func main() {
	http.HandleFunc("/cpu", handleCPU)
	http.HandleFunc("/memory", handleMemory)
	http.HandleFunc("/disk", handleDisk)
	http.HandleFunc("/processes", handleProcesses)
	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/version", handleVersion)

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

type response struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

func handleCPU(w http.ResponseWriter, r *http.Request) {
	numCPU := runtime.NumCPU()
	data := map[string]int{"cpu_count": numCPU}
	writeJSONResponse(w, "success", data)
}

func handleMemory(w http.ResponseWriter, r *http.Request) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	data := map[string]uint64{
		"memory_alloc":     memStats.Alloc,
		"memory_sys":       memStats.Sys,
		"memory_heap_alloc": memStats.HeapAlloc,
	}
	writeJSONResponse(w, "success", data)
}

func handleDisk(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("sh", "-c", "df -h")
	output, err := cmd.Output()
	if err != nil {
		writeJSONResponse(w, "failure", "Failed to execute command")
		return
	}
	writeJSONResponse(w, "success", string(output))
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	cmd := exec.Command("sh", "-c", "ps aux | grep " + filter)
	output, err := cmd.Output()
	if err != nil {
		writeJSONResponse(w, "failure", "Failed to execute command")
		return
	}
	writeJSONResponse(w, "success", string(output))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "success", "Program Status: Running")
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "success", "Agent Version: 1.0.0")
}

func writeJSONResponse(w http.ResponseWriter, result string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := response{
		Result: result,
		Data:   data,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}