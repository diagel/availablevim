package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	// Сначала регистрируем API-роуты
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		os.WriteFile("saved.go", body, 0644)
		fmt.Fprintf(w, "Saved")
	})

	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		data, err := os.ReadFile("saved.go")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(data)
	})

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		code, _ := io.ReadAll(r.Body)

		tmpfile, err := os.CreateTemp("", "go-code-*.go")
		if err != nil {
			http.Error(w, "Could not create temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.Write(code)
		if err != nil {
			http.Error(w, "Could not write to file", http.StatusInternalServerError)
			return
		}

		cmd := exec.Command("go", "run", tmpfile.Name())
		output, err := cmd.CombinedOutput()

		if err != nil {
			fmt.Fprintf(w, "Error:\n%s\nOutput:\n%s", err, output)
			return
		}

		w.Write(output)
	})

	// Статический сервер (после API)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	fmt.Println("Сервер запущен на http://localhost:4490")
	log.Fatal(http.ListenAndServe(":4490", nil))
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
