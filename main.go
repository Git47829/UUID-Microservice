package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const schema =`
CREATE TABLE IF NOT EXISTS IDS (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	UUID BLOB NOT NULL UNIQUE
);
`

var DB *sql.DB

func generateUUID() uuid.UUID {
	return uuid.New()
}


func checkUUID(ID uuid.UUID) bool {
	err := DB.QueryRow(`SELECT UUID FROM IDS WHERE UUID = ?;`, ID).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}
	return true
}

func validateUUID(w http.ResponseWriter, r *http.Request) {
	rawId  := r.URL.Query().Get("id")
	w.Header().Set("content-type", "application/json")
	if rawId == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"valid": false,
			"error": "UUID cannot be NULL",
		})
		return
	}
	id, err := uuid.Parse(rawId)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"valid": false,
			"error": err,
		})
		return
	}
	if checkUUID(id) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]bool{
			"valid": false,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"valid": true,
	})
}

func registerUUID(w http.ResponseWriter, r *http.Request) {
	id := generateUUID()
	if checkUUID(id) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "UUID already exists",
		})
		return
	}
	query := `INSERT INTO IDS (UUID) VALUES (?);`
	_, err := DB.Query(query, id)
	if err != nil  {
		fmt.Println("Error in SQL Query: ", err)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]uuid.UUID{
		"id": id,
	})
}

func initDatabase() {
	var err error
	DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatalf("Failed to open Database")
	}
	initSchema()
	fmt.Println("Successfully Initiated Database Schema")	


	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to Ping Database %v", err)
	}

	fmt.Println("Successfully connected to SQLite Database")
}

func initSchema() error {
 _, err := DB.Exec(schema)
 return err
}

func getEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load env variables %v", err)
	}
	return os.Getenv(key)
}

func initWebServer() {
	router := http.NewServeMux()
	router.HandleFunc("POST /generateUUID", registerUUID)
	router.HandleFunc("GET /validateUUID", validateUUID)
	err :=	http.ListenAndServe(":" + getEnv("PORT"), router)
	if err != nil {
		fmt.Println("Failed to Start Webserver. Error: ", err)
	}
	fmt.Println("Listening and Serving on Port " + getEnv("PORT"))
}

func main() {
	initDatabase()
	initWebServer()
}
