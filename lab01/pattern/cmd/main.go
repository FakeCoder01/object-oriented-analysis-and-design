package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cslab/pattern/internal/database"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

type DBStats struct {
	InstanceID      int    `json:"instance_id"`
	GetInstanceCalls int   `json:"get_instance_calls"`
	OpenConnections int64  `json:"open_connections"`
	Message         string `json:"message"`
	Pattern         string `json:"pattern"`
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance().GetConn()
	log.Println("[SINGLETON] getTasks: reusing the singleton DB instance")

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
	db := database.GetInstance().GetConn()
	log.Println("[SINGLETON] createTask: reusing the singleton DB instance")

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
	db := database.GetInstance().GetConn()
	log.Println("[SINGLETON] toggleTask: reusing the singleton DB instance")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	db.Exec("UPDATE tasks SET done = NOT done WHERE id = ?", id)
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	db := database.GetInstance().GetConn()
	log.Println("[SINGLETON] deleteTask: reusing the singleton DB instance")

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	db.Exec("DELETE FROM tasks WHERE id = ?", id)
	w.WriteHeader(http.StatusOK)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	dbInstance := database.GetInstance()
	instanceID, callCount, openConns := dbInstance.GetStats()
	log.Println("[SINGLETON] getStats: reusing the singleton DB instance")

	var taskCount int
	dbInstance.GetConn().QueryRow("SELECT COUNT(*) FROM tasks").Scan(&taskCount)

	stats := DBStats{
		InstanceID:       instanceID,
		GetInstanceCalls: callCount,
		OpenConnections:  openConns,
		Message: fmt.Sprintf("(#%d) GetInstance() был назван %d раз. Задачи: %d", instanceID, callCount, taskCount),
		Pattern: "singleton",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	database.GetInstance().Init()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("/app/static"))
	mux.Handle("/pattern/", http.StripPrefix("/pattern", fs))

	mux.HandleFunc("/pattern/api/tasks", func(w http.ResponseWriter, r *http.Request) {
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
	mux.HandleFunc("/pattern/api/tasks/toggle", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		toggleTask(w, r)
	})
	mux.HandleFunc("/pattern/api/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		deleteTask(w, r)
	})
	mux.HandleFunc("/pattern/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		getStats(w, r)
	})

	log.Println("[SINGLETON] Server starting on :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
