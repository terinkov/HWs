package models

import (
	"log"
	"time"
)

type Task struct {
	UUID    string `json:"-"`  
	File    string `json:"image"`
	Filter  string `json:"filter"`
	Status  string `json:"-"`
	Content string `json:"-"`
}

func (task *Task) DoTask() {
	time.Sleep(15 * time.Second)
	task.Status = "ready"
	task.Content = "RESULT"
	log.Printf("task %s complited\n", task.UUID)
}
