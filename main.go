package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTP METHOD
// - GET (Mengambil data dari server)
// - POST (Mengirim data ke server, biasanya digunakan untuk menambahkan data)
// - PUT (Mengubah data yang sudah ada di server)
// - PATCH (Mengubah data yang sudah ada di server)
// - DELETE (Menghapus data yang sudah ada di server)

type Biodata struct {
	Name    string `json:"nama"`
	Age     int    `json:"umur"`
	Address string `json:"alamat"`
}

func (bio *Biodata) IsValid() bool {
	if bio.Name == "" || bio.Age == 0 || bio.Address == "" {
		return false
	}
	return true
}

var data = []Biodata{
	{Name: "Della", Age: 20, Address: "Jakarta"},
	{Name: "John", Age: 25, Address: "New York"},
	{Name: "Jane", Age: 30, Address: "London"},
	{Name: "Vika", Age: 27, Address: "Pekalongan"},
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func Greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Apa kabar!")
}

func BiodataFunc(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var req Biodata
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !req.IsValid() {
			http.Error(w, "Invalid biodata", http.StatusBadRequest)
			return
		}

		data = append(data, req)

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Biodata berhasil ditambahkan",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Della!")
	})
	http.HandleFunc("/sapa", Greeting)
	http.HandleFunc("/biodata", BiodataFunc)

	fmt.Println("Server run on port 3000")
	http.ListenAndServe(":3000", nil)
}
