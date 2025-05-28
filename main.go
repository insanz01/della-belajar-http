package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	switch r.Method {
	case http.MethodGet:
		queryName := r.URL.Query().Get("name")
		fmt.Println(queryName)
		if queryName != "" {
			for _, bio := range data {
				if bio.Name == queryName {
					result, err := json.Marshal(bio)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Write(result)
					return
				}
			}
			http.Error(w, "Biodata not found", http.StatusNotFound)
			return
		}

		result, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)

	case http.MethodPost:
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

	case http.MethodPut:
		// Untuk PUT dan PATCH, kita perlu tahu biodata mana yang akan diubah.
		// Asumsi kita akan menggunakan parameter `name` di URL sebagai ID.
		// Contoh: PUT /biodata?name=Della
		queryName := r.URL.Query().Get("name")
		if queryName == "" {
			http.Error(w, "Parameter 'name' dibutuhkan untuk update", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var updatedBio Biodata
		err = json.Unmarshal(body, &updatedBio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		found := false
		for i, bio := range data {
			if bio.Name == queryName {
				// PUT mengganti seluruh data
				if !updatedBio.IsValid() {
					http.Error(w, "Invalid biodata for update", http.StatusBadRequest)
					return
				}
				data[i] = updatedBio
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Biodata tidak ditemukan", http.StatusNotFound)
			return
		}

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Biodata berhasil diperbarui (PUT)",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)

	case http.MethodPatch:
		// Untuk PATCH, kita perlu tahu biodata mana yang akan diubah.
		// Asumsi kita akan menggunakan parameter `name` di URL sebagai ID.
		// Contoh: PATCH /biodata?name=Della
		queryName := r.URL.Query().Get("name")
		if queryName == "" {
			http.Error(w, "Parameter 'name' dibutuhkan untuk update", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var partialBio Biodata
		err = json.Unmarshal(body, &partialBio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		found := false
		for i, bio := range data {
			if bio.Name == queryName {
				// PATCH hanya mengubah field yang diberikan
				if partialBio.Name != "" {
					data[i].Name = partialBio.Name
				}
				if partialBio.Age != 0 {
					data[i].Age = partialBio.Age
				}
				if partialBio.Address != "" {
					data[i].Address = partialBio.Address
				}
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Biodata tidak ditemukan", http.StatusNotFound)
			return
		}

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Biodata berhasil diperbarui (PATCH)",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)

	case http.MethodDelete:
		// Untuk DELETE, kita perlu tahu biodata mana yang akan dihapus.
		// Asumsi kita akan menggunakan parameter `name` di URL sebagai ID.
		// Contoh: DELETE /biodata?name=Della
		queryName := r.URL.Query().Get("name")
		if queryName == "" {
			http.Error(w, "Parameter 'name' dibutuhkan untuk delete", http.StatusBadRequest)
			return
		}

		foundIndex := -1
		for i, bio := range data {
			if bio.Name == queryName {
				foundIndex = i
				break
			}
		}

		if foundIndex == -1 {
			http.Error(w, "Biodata tidak ditemukan", http.StatusNotFound)
			return
		}

		// Menghapus elemen dari slice
		data = append(data[:foundIndex], data[foundIndex+1:]...)

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Biodata berhasil dihapus",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func main() {
	PORT := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Della!")
	})
	http.HandleFunc("/sapa", Greeting)
	http.HandleFunc("/biodata", BiodataFunc)

	fmt.Println("Server running on port", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
