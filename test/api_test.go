package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestTaskAdd(t *testing.T) {
	addr := "http://192.168.3.4:80"

	task := map[string]interface{}{
		"name":     "become human",
		"complete": true,
	}

	taskJson, err := json.Marshal(&task)

	if err != nil {
		t.Errorf("Failed marshalling task: %v", err)
	}

	http.Post(addr+"/tasks", "application/json", bytes.NewBuffer(taskJson))

	resp, err := http.Get(addr + "/tasks")

	if err != nil {
		t.Errorf("Error while getting tasks: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Error while reading get tasks request body: %v", err)
	}

	var tasks = make([]map[string]interface{}, 10)

	json.Unmarshal(body, &tasks)

	lastTask := tasks[len(tasks)-1]
	if !(lastTask["Name"] == "become human" && lastTask["Complete"] == true) {
		t.Errorf("Last task is not the added one, expected %v, got %v", task, lastTask)
	}
}
