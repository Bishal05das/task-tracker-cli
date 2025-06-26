package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func main() {
	//provide arguments
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go <command> [arguments]")
		return
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go add <description>")
			return
		}
		description := os.Args[2]
		addTask(description)
	case "list":
		listTasks()
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go update <id> <new_description>")
			return
		}
		id := ParseID(os.Args[2])
		newDescription := os.Args[3]
		updateTask(id, newDescription)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("usage: go run main.go delete <id>")
			return
		}
		id := ParseID(os.Args[2])
		deleteTask(id)

	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run main.go mark-in-progress <id>")
			return
		}
		id := ParseID(os.Args[2])
		updateStatus(id, "in-progress")

	case "list-task-by-status":
		if len(os.Args) == 3 {
			status := os.Args[2]
			if status == "todo" || status == "in-progress" || status == "done" {
				listTasksByStatus(status)

			} else {
				fmt.Println("Invalid status. Use 'todo', 'in-progress', or 'done'.")

			}
		} 

	default:
		fmt.Println("Unknown command:", command)

	}

}

func ParseID(idStr string) int {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID. ID must be a number.")
		os.Exit(1)
	}
	return id
}

func readTasks() []Task {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}
	
	return tasks
}

func saveTasks(tasks []Task) {
	file, err := os.Create("tasks.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
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
		ID:          len(tasks) + 1,
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
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
		fmt.Printf("ID: %d, Description: %s, Status: %s\n", task.ID, task.Description, task.Status)
		
	}
}

//update task

func updateTask(id int, newDescription string) {
	tasks := readTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			saveTasks(tasks)
			fmt.Println("Task updated successfully!")
			return
		}
	}
	fmt.Println("Task not found.")
}

//delete task

func deleteTask(id int) {
	tasks := readTasks()
	newtasks := []Task{}

	for _, task := range tasks {
		if id != task.ID {
			newtasks = append(newtasks, task)
		}
	}
	saveTasks(newtasks)
	fmt.Println("Task deleted successfully!")
}

//update status

func updateStatus(id int, status string) {
	tasks := readTasks()

	for _, task := range tasks {
		if id == task.ID {
			tasks[id].Status = status
			tasks[id].UpdatedAt = time.Now().Format(time.RFC3339)
			saveTasks(tasks)
			fmt.Println("Task status updated successfully!")
			return
		}
	}
	fmt.Println("Task not found.")
}

//filter by status

func listTasksByStatus(status string) {
	tasks := readTasks()
	filterTasks := []Task{}
	for _, task := range tasks {
		if task.Status == status {
			filterTasks = append(filterTasks, task)

		}
	}
	if len(filterTasks) == 0 {
		fmt.Println("No task found with status:", status)
		return
	}
	for _, task := range filterTasks {
		fmt.Printf("ID: %d | Descrption: %s | Status: %s\n", task.ID, task.Description, task.Status)
		
	}
}
