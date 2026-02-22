package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"baby-care/internal/server"
	"baby-care/internal/store"
)

func main() {
	port := flag.Int("port", defaultPort(), "HTTP port")
	dbPath := flag.String("db", defaultDBPath(), "SQLite database path")
	flag.Parse()

	st, err := store.Open(*dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("sub static: %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Baby Care Tracker listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.New(st, staticFS)); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func defaultPort() int {
	if p := os.Getenv("PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			return n
		}
	}
	return 8080
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "data.db"
	}
	return filepath.Join(home, ".baby-care", "data.db")
}
