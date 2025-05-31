package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Apa kabar!")
}

func KontakHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		db, err := InitDB()
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var req KontakTeman
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !req.IsValid() {
			http.Error(w, "Invalid kontak", http.StatusBadRequest)
			return
		}

		// TODO : logic database
		if err := req.InsertToDB(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Kontak berhasil ditambahkan",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
	case http.MethodGet:
		db, err := InitDB()
		if err != nil {
			log.Fatal(err)
		}

		var kontak KontakTeman

		queryName := r.URL.Query().Get("nama")
		fmt.Println(queryName)
		if queryName != "" {
			k, err := kontak.GetByName(db, queryName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if k == nil {
				http.Error(w, "Kontak not found", http.StatusNotFound)
				return
			}

			result, err := json.Marshal(k)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(result)
			return
		}

		kontaks, err := kontak.GetAll(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(kontaks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	case http.MethodPut:
		db, err := InitDB()
		if err != nil {
			log.Fatal(err)
		}

		queryName := r.URL.Query().Get("id")
		if queryName == "" {
			http.Error(w, "Parameter 'id' dibutuhkan untuk update", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(queryName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var req KontakTeman
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !req.IsValid() {
			http.Error(w, "Invalid kontak", http.StatusBadRequest)
			return
		}

		k, err := req.GetByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if k == nil {
			http.Error(w, "Kontak tidak ditemukan", http.StatusNotFound)
			return
		}

		req.ID = id
		if err := req.Update(db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Kontak berhasil diubah",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}

		w.Write(result)
		return
	case http.MethodDelete:
		db, err := InitDB()
		if err != nil {
			log.Fatal(err)
		}

		queryName := r.URL.Query().Get("id")
		if queryName == "" {
			http.Error(w, "Parameter 'id' dibutuhkan untuk delete", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(queryName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO : cek apakah id kontak teman nya ada
		var kontak KontakTeman
		k, err := kontak.GetByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO : jika kontak teman tidak ada, kembalikan not found
		if k == nil {
			http.Error(w, "Kontak tidak ditemukan", http.StatusNotFound)
			return
		}

		// TODO : hapus jika kontak teman ada
		err = k.Delete(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseMessage := struct {
			Message string `json:"message"`
		}{
			Message: "Kontak berhasil dihapus",
		}

		result, err := json.Marshal(responseMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
		return

	}
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
