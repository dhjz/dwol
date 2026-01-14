package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func sendJSON(w http.ResponseWriter, v interface{}) error {
	setJSONHeader(w)
	return json.NewEncoder(w).Encode(v)
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func HandleMachines(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		machines := GetMachines()
		sendJSON(w, machines)
	case http.MethodPost:
		var m Machine
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			sendError(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if m.ID == "" || m.MAC == "" {
			sendError(w, "ID and MAC are required", http.StatusBadRequest)
			return
		}
		if m.Port == 0 {
			m.Port = 9
		}
		if err := AddMachine(m); err != nil {
			sendError(w, "Failed to add machine", http.StatusInternalServerError)
			return
		}
		sendJSON(w, m)
	case http.MethodPut:
		var m Machine
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			sendError(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if m.ID == "" || m.MAC == "" {
			sendError(w, "ID and MAC are required", http.StatusBadRequest)
			return
		}
		if m.Port == 0 {
			m.Port = 9
		}
		if err := UpdateMachineInfo(m); err != nil {
			sendError(w, "Failed to update machine", http.StatusInternalServerError)
			return
		}
		sendJSON(w, m)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			sendError(w, "ID is required", http.StatusBadRequest)
			return
		}
		if err := DeleteMachine(id); err != nil {
			sendError(w, "Failed to delete machine", http.StatusInternalServerError)
			return
		}
		sendJSON(w, map[string]string{"status": "deleted"})
	default:
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandleWake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	machines := GetMachines()
	var target *Machine
	for i := range machines {
		if machines[i].ID == req.ID {
			target = &machines[i]
			break
		}
	}

	if target == nil {
		sendError(w, "Machine not found", http.StatusNotFound)
		return
	}

	broadcast := "255.255.255.255"
	if target.Host != "" {
		broadcast = target.Host
	}

	// log.Printf("Sending WOL packet to %s (MAC: %s, Port: %d, Broadcast: %s)", target.Host, target.MAC, target.Port, broadcast)

	if err := SendWOL(target.MAC, broadcast, target.Port); err != nil {
		log.Printf("Failed to send WOL: %v", err)
		sendError(w, "Failed to send WOL packet", http.StatusInternalServerError)
		return
	}

	log.Printf("WOL packet sent successfully to %s (MAC: %s)", target.Host, target.MAC)

	sendJSON(w, map[string]string{"status": "success"})
}

func HandleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newMachines []Machine
	if err := json.NewDecoder(r.Body).Decode(&newMachines); err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := SetMachines(newMachines); err != nil {
		sendError(w, "Failed to import data", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]string{"status": "imported"})
}

func HandleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	machines := GetMachines()
	setJSONHeader(w)
	w.Header().Set("Content-Disposition", "attachment; filename=machines.json")
	json.NewEncoder(w).Encode(machines)
}

func HandleAuthStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sendJSON(w, map[string]bool{"required": IsPasswordRequired()})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if CheckPassword(req.Password) {
		sendJSON(w, map[string]string{"status": "success"})
	} else {
		sendError(w, "Invalid password", http.StatusUnauthorized)
	}
}
