package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/go-chi/chi"
)

type LoadCpuUtilRequest struct {
	Cores   int `json:"cores"`
	Util    int `json:"util"`
	Timeout int `json:"timeout"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := chi.NewMux()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response from backend %s\n", port)
	})

	r.Post("/load/cpu", func(w http.ResponseWriter, r *http.Request) {

		var loadCpuUtilRequest LoadCpuUtilRequest
		err := json.NewDecoder(r.Body).Decode(&loadCpuUtilRequest)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		cmd := exec.Command(
			"stress-ng",
			"--cpu", strconv.Itoa(loadCpuUtilRequest.Cores),
			"--cpu-load", strconv.Itoa(loadCpuUtilRequest.Util),
			"--timeout", strconv.Itoa(loadCpuUtilRequest.Timeout),
		)

		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println(output)

	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting backend on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
