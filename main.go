package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func main() {
	basePath := "/api/v1"
	mux := http.NewServeMux()
	baseURLMux := http.StripPrefix(basePath, mux)
	mux.HandleFunc("/users", listUserHandler)
	mux.HandleFunc("POST /user", createUserHandler)
	http.ListenAndServe(":8080", baseURLMux)
}

func listUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, &u)
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding JSON", http.StatusBadRequest)
		return
	}

	if _, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", u.Name, u.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
