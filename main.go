package main

import (
	"fmt"
	"log"
	"net/http"
)

// HTTP METHOD
// - GET (Mengambil data dari server)
// - POST (Mengirim data ke server, biasanya digunakan untuk menambahkan data)
// - PUT (Mengubah data yang sudah ada di server)
// - PATCH (Mengubah data yang sudah ada di server)
// - DELETE (Menghapus data yang sudah ada di server)

var data = []Biodata{
	{Name: "Della", Age: 20, Address: "Jakarta"},
	{Name: "John", Age: 25, Address: "New York"},
	{Name: "Jane", Age: 30, Address: "London"},
	{Name: "Vika", Age: 27, Address: "Pekalongan"},
}

func main() {
	PORT := ":3000"

	http.HandleFunc("/sapa", Greeting)
	http.HandleFunc("/biodata", BiodataFunc)
	http.HandleFunc("/kontak", KontakHandler)

	fmt.Println("Server running on port", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}
