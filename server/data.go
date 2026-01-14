package server

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	dataFile string
	mu       sync.RWMutex
	machines []Machine
)

type Machine struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
	MAC  string `json:"mac"`
	Port int    `json:"port"`
}

func InitData(file string) error {
	dataFile = file
	return loadData()
}

func loadData() error {
	mu.Lock()
	defer mu.Unlock()

	content, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			machines = []Machine{}
			return saveData()
		}
		return err
	}

	return json.Unmarshal(content, &machines)
}

func saveData() error {
	content, err := json.MarshalIndent(machines, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFile, content, 0644)
}

func GetMachines() []Machine {
	mu.RLock()
	defer mu.RUnlock()
	return machines
}

func AddMachine(m Machine) error {
	mu.Lock()
	defer mu.Unlock()
	machines = append(machines, m)
	return saveData()
}

func DeleteMachine(id string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, m := range machines {
		if m.ID == id {
			machines = append(machines[:i], machines[i+1:]...)
			return saveData()
		}
	}
	return nil
}

func UpdateMachineInfo(m Machine) error {
	mu.Lock()
	defer mu.Unlock()

	for i := range machines {
		if machines[i].ID == m.ID {
			machines[i].Name = m.Name
			machines[i].Host = m.Host
			machines[i].MAC = m.MAC
			machines[i].Port = m.Port
			return saveData()
		}
	}
	return nil
}

func SetMachines(newMachines []Machine) error {
	mu.Lock()
	defer mu.Unlock()
	machines = newMachines
	return saveData()
}
