package main

import (
	"dwol/server"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed public
var publicFS embed.FS

func main() {
	port := flag.String("p", "999", "HTTP server port")
	pwd := flag.String("pwd", "adminer", "Login password (empty means no password required)")
	flag.Parse()

	server.SetPassword(*pwd)

	execPath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Failed to get executable path: %v", err)
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	dataFile := filepath.Join(execDir, "data.json")
	if err := server.InitData(dataFile); err != nil {
		log.Printf("Warning: Failed to init data file: %v", err)
	}

	cronFile := filepath.Join(execDir, "cron.json")
	server.InitCron(cronFile)

	subFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	http.Handle("/", http.FileServer(http.FS(subFS)))
	http.HandleFunc("/api/machines", server.AuthMiddleware(server.HandleMachines))
	http.HandleFunc("/api/wake", server.AuthMiddleware(server.HandleWake))
	http.HandleFunc("/api/import", server.AuthMiddleware(server.HandleImport))
	http.HandleFunc("/api/export", server.AuthMiddleware(server.HandleExport))
	http.HandleFunc("/api/auth/status", server.HandleAuthStatus)
	http.HandleFunc("/api/login", server.HandleLogin)
	http.HandleFunc("/api/status", server.AuthMiddleware(server.HandleStatus))
	http.HandleFunc("/api/test-port", server.AuthMiddleware(server.HandleTestPort))
	http.HandleFunc("/api/cron-tasks", server.AuthMiddleware(server.HandleCronTasks))
	http.HandleFunc("/api/validate-cron", server.AuthMiddleware(server.HandleValidateCron))

	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("Server starting on http://localhost%s\n", addr)
	fmt.Printf("Data directory: %s\n", execDir)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
