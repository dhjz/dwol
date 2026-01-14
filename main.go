package main

import (
	"dwol/server"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed public
var publicFS embed.FS

func main() {
	port := flag.String("p", "999", "HTTP server port")
	pwd := flag.String("pwd", "adminer", "Login password (empty means no password required)")
	flag.Parse()

	server.SetPassword(*pwd)

	dataFile := "data.json"
	if err := server.InitData(dataFile); err != nil {
		log.Printf("Warning: Failed to init data file: %v", err)
	}

	subFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	http.Handle("/", http.FileServer(http.FS(subFS)))
	http.HandleFunc("/api/machines", server.HandleMachines)
	http.HandleFunc("/api/wake", server.HandleWake)
	http.HandleFunc("/api/import", server.HandleImport)
	http.HandleFunc("/api/export", server.HandleExport)
	http.HandleFunc("/api/auth/status", server.HandleAuthStatus)
	http.HandleFunc("/api/login", server.HandleLogin)
	http.HandleFunc("/api/status", server.HandleStatus)
	http.HandleFunc("/api/test-port", server.HandleTestPort)

	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
