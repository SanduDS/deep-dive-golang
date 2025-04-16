package main

import (
	"log"
	"time"
)

// Worker represents a worker that processes tasks
type Worker struct {
	id int
}

// Task represents a unit of work to be processed
type Task struct {
	id        int
	processed bool
}

func main() {
	// Create a channel to send tasks to workers
	// The channel can hold up to 5 tasks at a time
	taskChannel := make(chan Task, 5)

	// Create a channel to receive results from workers
	resultChannel := make(chan Task, 5)

	// Number of workers we want to create
	numWorkers := 3

	// Start the workers
	for i := 1; i <= numWorkers; i++ {
		worker := Worker{id: i}
		// Launch a goroutine for each worker
		go worker.processTasks(taskChannel, resultChannel)
	}

	// Create some tasks
	tasks := []Task{
		{id: 1},
		{id: 2},
		{id: 3},
		{id: 4},
		{id: 5},
	}

	// Send tasks to the workers
	log.Println("📦 Sending tasks to workers...")
	for _, task := range tasks {
		log.Printf("📤 Sending task %d to workers", task.id)
		taskChannel <- task
		// Simulate some delay between sending tasks
		time.Sleep(500 * time.Millisecond)
	}

	// Close the task channel to signal that no more tasks are coming
	close(taskChannel)

	// Collect results from workers
	log.Println("🔄 Collecting results from workers...")
	for i := 0; i < len(tasks); i++ {
		result := <-resultChannel
		log.Printf("✅ Task %d has been processed by a worker", result.id)
	}

	log.Println("🎉 All tasks completed!")
}

// processTasks is a method that processes tasks from the taskChannel
// and sends the results to the resultChannel
func (w Worker) processTasks(taskChannel <-chan Task, resultChannel chan<- Task) {
	log.Printf("👷 Worker %d started and waiting for tasks", w.id)

	// Loop through tasks in the channel
	for task := range taskChannel {
		log.Printf("🔄 Worker %d is processing task %d", w.id, task.id)

		// Simulate some work being done
		time.Sleep(1 * time.Second)

		// Mark the task as processed
		task.processed = true

		// Send the processed task back through the result channel
		log.Printf("📤 Worker %d finished processing task %d", w.id, task.id)
		resultChannel <- task
	}

	log.Printf("🏁 Worker %d has finished all tasks and is shutting down", w.id)
}
