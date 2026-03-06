package server

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronScheduler *cron.Cron
	cronTasks     []CronTask
	muCron        sync.RWMutex
	cronFile      string
)

type CronTask struct {
	ID        string `json:"id"`
	MachineID string `json:"machineId"`
	CronExpr  string `json:"cronExpr"`
	Enabled   bool   `json:"enabled"`
	NextRun   string `json:"nextRun,omitempty"`
	LastRun   string `json:"lastRun,omitempty"`
}

func InitCron(file string) {
	cronFile = file
	cronScheduler = cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
	cronScheduler.Start()
	loadCronTasks()
}

func loadCronTasks() error {
	muCron.Lock()
	defer muCron.Unlock()

	content, err := os.ReadFile(cronFile)
	if err != nil {
		if os.IsNotExist(err) {
			cronTasks = []CronTask{}
			return saveCronTasks()
		}
		return err
	}

	return json.Unmarshal(content, &cronTasks)
}

func saveCronTasks() error {
	content, err := json.MarshalIndent(cronTasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cronFile, content, 0644)
}

func GetCronTasks() []CronTask {
	muCron.RLock()
	defer muCron.RUnlock()
	return cronTasks
}

func AddCronTask(task CronTask) error {
	muCron.Lock()
	defer muCron.Unlock()

	for _, t := range cronTasks {
		if t.ID == task.ID {
			return nil
		}
	}

	cronTasks = append(cronTasks, task)
	if task.Enabled {
		scheduleCronTask(task)
	}
	return saveCronTasks()
}

func UpdateCronTask(task CronTask) error {
	muCron.Lock()
	defer muCron.Unlock()

	for i, t := range cronTasks {
		if t.ID == task.ID {
			cronTasks[i] = task
			if task.Enabled {
				scheduleCronTask(task)
			}
			return saveCronTasks()
		}
	}
	return nil
}

func DeleteCronTask(id string) error {
	muCron.Lock()
	defer muCron.Unlock()

	for i, t := range cronTasks {
		if t.ID == id {
			cronTasks = append(cronTasks[:i], cronTasks[i+1:]...)
			return saveCronTasks()
		}
	}
	return nil
}

func scheduleCronTask(task CronTask) {
	entryID, err := cronScheduler.AddFunc(task.CronExpr, func() {
		executeCronTask(task)
	})
	if err != nil {
		log.Printf("Failed to schedule cron task %s: %v", task.ID, err)
		return
	}
	log.Printf("Scheduled cron task %s with entry ID %d, machineID: %s, cronExpr: %s", task.ID, entryID, task.MachineID, task.CronExpr)
}

func executeCronTask(task CronTask) {
	log.Printf("Executing cron task %s for machine %s", task.ID, task.MachineID)

	machines := GetMachines()
	var target *Machine
	for i := range machines {
		if machines[i].ID == task.MachineID {
			target = &machines[i]
			break
		}
	}

	if target == nil {
		log.Printf("Machine %s not found for cron task %s", task.MachineID, task.ID)
		return
	}

	broadcast := "255.255.255.255"
	if target.Host != "" {
		if !isPrivateIP(target.Host) {
			broadcast = target.Host
		}
	}

	if err := SendWOL(target.MAC, broadcast, target.Port); err != nil {
		log.Printf("Failed to send WOL for cron task %s: %v", task.ID, err)
	} else {
		log.Printf("WOL sent successfully for cron task %s, cronExpr: %s, machineID: %s", task.ID, task.CronExpr, task.MachineID)
		muCron.Lock()
		for i, t := range cronTasks {
			if t.ID == task.ID {
				cronTasks[i].LastRun = time.Now().Format("2006-01-02 15:04:05")
				break
			}
		}
		muCron.Unlock()
		saveCronTasks()
	}
}

func GetNextRun(cronExpr string) (string, error) {
	schedule, err := cron.ParseStandard(cronExpr)
	if err != nil {
		return "", err
	}
	nextTime := schedule.Next(time.Now())
	return nextTime.Format("2006-01-02 15:04:05"), nil
}

func ValidateCronExpr(expr string) bool {
	_, err := cron.ParseStandard(expr)
	return err == nil
}
