package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const todoFile = "todos.json"

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(todoFile)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(todoFile, data, 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cli-todo-list <add|list|done|delete> [args]")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("error loading tasks:", err)
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("please give task title")
			return
		}
		title := strings.Join(os.Args[2:], " ")
		nextID := 1
		for _, t := range tasks {
			if t.ID >= nextID {
				nextID = t.ID + 1
			}
		}
		tasks = append(tasks, Task{ID: nextID, Title: title, Done: false})
		if err := saveTasks(tasks); err != nil {
			fmt.Println("save error:", err)
			return
