package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID      	int
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func main() {
	//provide arguments
	if len(os.Args) < 2 {
		fmt.Println("usage: task-cli <command> [arguments]")
		return
	}

	command := os.Args[1]
	switch command {
	    case "add":
            if len(os.Args) < 3 {
			fmt.Println("Usage: task-cli add <description>")
			return
		    }
		 	description := os.Args[2]
		 	addTask(description)
        case "list":
			listTasks()
		
		default:
			fmt.Println("Unknown command:",command)
		 

	}


}

func readTasks() []Task {
	file,err := os.OpenFile("tasks.json", os.O_RDWR| os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:",err)
		return nil
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error decoding JSON:",err)
		return nil
	}
	return tasks
}

func saveTasks(tasks []Task) {
	file,err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error creating file:",err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}

}




func addTask(description string) {
	tasks := readTasks()
	newTask := Task{
		ID:     len(tasks) + 1,
		Description: description,
		Status: "todo",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Println("Task added successfully!")

}

func listTasks() {
	tasks := readTasks()
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range tasks {
		fmt.Printf("ID: %d, Description: %s, Status: %s\n",task.ID, task.Description, task.Status)
		return
	}
}
