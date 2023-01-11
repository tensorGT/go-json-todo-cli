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
