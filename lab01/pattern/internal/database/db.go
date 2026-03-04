package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)


type Database struct {
	conn         *sql.DB
	instanceID   int
	openCount    int
}

var (
	instance *Database
	once     sync.Once  // sync.Once guarantees the init function runs EXACTLY ONCE
	mu       sync.Mutex
)


func GetInstance() *Database {
	once.Do(func() {
		log.Println("[SINGLETON] Creating the ONE AND ONLY Database instance...")
		db, err := sql.Open("sqlite3", "/data/tasks_pattern.db")
		if err != nil {
			log.Fatalf("Failed to open DB: %v", err)
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)

		instance = &Database{
			conn:       db,
			instanceID: 1,
			openCount:  1,
		}
		log.Println("[SINGLETON] Database instance #1 created. This will NEVER be called again.")
	})

	// times GetInstance was called
	mu.Lock()
	instance.openCount++
	mu.Unlock()

	log.Printf("[SINGLETON] GetInstance() called. Returning the SAME instance #%d (called %d times total)",
		instance.instanceID, instance.openCount)
	return instance
}

func (d *Database) GetConn() *sql.DB {
	return d.conn
}

func (d *Database) GetStats() (int, int, int64) {
	stats := d.conn.Stats()
	return d.instanceID, d.openCount, int64(stats.OpenConnections)
}

func (d *Database) Init() {
	_, err := d.conn.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN DEFAULT FALSE,
		priority TEXT DEFAULT 'medium',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("[SINGLETON] DB initialized via singleton instance")
}
