package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)


type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

type DBStats struct {
	OpenConnections int    `json:"open_connections"`
	Message         string `json:"message"`
	Pattern         string `json:"pattern"`
}


func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", "/data/tasks_nopattern.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	return db
}

func initDB() {
	db := openDB()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE,
		priority TEXT DEFAULT 'medium',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("[NO-PATTERN] DB initialized (new connection opened and closed)")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	log.Println("[NO-PATTERN] getTasks: opened a new DB connection")

	rows, err := db.Query("SELECT id, title, done, priority, created_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Done, &t.Priority, &t.CreatedAt)
		tasks = append(tasks, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	log.Println("[NO-PATTERN] createTask: opened a new DB connection")

	var t Task
	json.NewDecoder(r.Body).Decode(&t)
	if t.Priority == "" {
		t.Priority = "medium"
	}

	result, err := db.Exec("INSERT INTO tasks (title, priority) VALUES (?, ?)", t.Title, t.Priority)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	id, _ := result.LastInsertId()
	t.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func toggleTask(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	log.Println("[NO-PATTERN] toggleTask: opened a new DB connection")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	db.Exec("UPDATE tasks SET done = NOT done WHERE id = ?", id)
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	log.Println("[NO-PATTERN] deleteTask: opened a new DB connection")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	db.Exec("DELETE FROM tasks WHERE id = ?", id)
	w.WriteHeader(http.StatusOK)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	log.Println("[NO-PATTERN] getStats: opened a new DB connection")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)

	stats := DBStats{
		OpenConnections: -1, // each connection is independent/can't track
		Message:         fmt.Sprintf("Задачи в базе данных: %d.", count),
		Pattern:         "no-pattern",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	initDB()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("/app/static"))
	mux.Handle("/no-pattern/", http.StripPrefix("/no-pattern", fs))

	mux.HandleFunc("/no-pattern/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		switch r.Method {
		case "GET":
			getTasks(w, r)
		case "POST":
			createTask(w, r)
		}
	})
	mux.HandleFunc("/no-pattern/api/tasks/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		toggleTask(w, r)
	})
	mux.HandleFunc("/no-pattern/api/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		deleteTask(w, r)
	})
	mux.HandleFunc("/no-pattern/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		getStats(w, r)
	})

	log.Println("[NO-PATTERN] Server starting on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
