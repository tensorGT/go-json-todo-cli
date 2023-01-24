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
		}
		fmt.Printf("added #%d: %s\n", nextID, title)

	case "list":
		if len(tasks) == 0 {
			fmt.Println("no tasks yet!")
			return
		}
		for _, t := range tasks {
			mark := "[ ]"
			if t.Done {
				mark = "[x]"
			}
			fmt.Printf("%s #%d %s\n", mark, t.ID, t.Title)
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("give task id")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Done = true
				found = true
				break
			}
		}
		if !found {
			fmt.Println("task not found")
			return
		}
		_ = saveTasks(tasks)
		fmt.Printf("task #%d marked done!\n", id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("give task id")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("invalid id")
			return
		}
		filtered := make([]Task, 0, len(tasks))
		removed := false
		for _, t := range tasks {
			if t.ID == id {
				removed = true
				continue
			}
			filtered = append(filtered, t)
		}
		if !removed {
			fmt.Println("task not found")
			return
		}
		_ = saveTasks(filtered)
		fmt.Printf("deleted #%d\n", id)

	default:
		fmt.Println("unknown command:", cmd)
	}
}
