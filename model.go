package main

import "database/sql"

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

type KontakTeman struct {
	ID      int    `json:"id"`
	Name    string `json:"nama"`
	Address string `json:"alamat"`
	Phone   string `json:"telepon"`
	Age     int    `json:"umur"`
}

func (kontak *KontakTeman) IsValid() bool {
	if kontak.Name == "" || kontak.Address == "" || kontak.Phone == "" || kontak.Age < 0 {
		return false
	}
	return true
}

func (kontak *KontakTeman) InsertToDB(db *sql.DB) error {
	query := "INSERT INTO teman (nama, alamat, nomor_hp, umur) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, kontak.Name, kontak.Address, kontak.Phone, kontak.Age)
	return err
}

func (kontak *KontakTeman) GetAll(db *sql.DB) ([]*KontakTeman, error) {
	rows, err := db.Query("SELECT id, nama, alamat, nomor_hp, umur FROM teman")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kontaks []*KontakTeman
	for rows.Next() {
		var kontak KontakTeman

		if err := rows.Scan(&kontak.ID, &kontak.Name, &kontak.Address, &kontak.Phone, &kontak.Age); err != nil {
			return nil, err
		}

		kontaks = append(kontaks, &kontak)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return kontaks, nil
}

func (kontak *KontakTeman) GetByID(db *sql.DB, id int) (*KontakTeman, error) {
	query := "SELECT id, nama, alamat, nomor_hp, umur FROM teman WHERE id = ?"
	row := db.QueryRow(query, id)

	var k KontakTeman
	if err := row.Scan(&k.ID, &k.Name, &k.Address, &k.Phone, &k.Age); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &k, nil
}

func (kontak *KontakTeman) GetByName(db *sql.DB, name string) (*KontakTeman, error) {
	query := "SELECT id, nama, alamat, nomor_hp, umur FROM teman WHERE nama = ?"
	row := db.QueryRow(query, name)

	var k KontakTeman
	if err := row.Scan(&k.ID, &k.Name, &k.Address, &k.Phone, &k.Age); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &k, nil
}

func (kontak *KontakTeman) Update(db *sql.DB) error {
	query := "UPDATE teman SET nama = ?, alamat = ?, nomor_hp = ?, umur = ? WHERE id = ?"
	_, err := db.Exec(query, kontak.Name, kontak.Address, kontak.Phone, kontak.Age, kontak.ID)
	return err
}

func (kontak *KontakTeman) UpdateName(db *sql.DB, name string) error {
	query := "UPDATE teman SET nama = ? WHERE id = ?"
	_, err := db.Exec(query, name, kontak.ID)
	return err
}

func (kontak *KontakTeman) UpdateAge(db *sql.DB, age int) error {
	query := "UPDATE teman SET umur = ? WHERE id = ?"
	_, err := db.Exec(query, age, kontak.ID)
	return err
}

func (kontak *KontakTeman) UpdatePhone(db *sql.DB, phone string) error {
	query := "UPDATE teman SET nomor_hp = ? WHERE id = ?"
	_, err := db.Exec(query, phone, kontak.ID)
	return err
}

func (kontak *KontakTeman) UpdateAddress(db *sql.DB, address string) error {
	query := "UPDATE teman SET alamat = ? WHERE id = ?"
	_, err := db.Exec(query, address, kontak.ID)
	return err
}

func (kontak *KontakTeman) Delete(db *sql.DB) error {
	query := "DELETE FROM teman WHERE id = ?"
	_, err := db.Exec(query, kontak.ID)
	return err
}
