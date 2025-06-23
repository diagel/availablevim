package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		body, _ := ioutil.ReadAll(r.Body)
		os.WriteFile("saved.txt", body, 0644)
		fmt.Fprintf(w, "Saved")
	})

	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("saved.txt")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(data)
	})

	fmt.Println("Сервер запущен на http://localhost:4490")
	http.ListenAndServe(":4490", nil)
}
