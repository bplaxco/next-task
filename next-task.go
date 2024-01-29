package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bplaxco/next-task/google"
	"github.com/bplaxco/next-task/tasks"
)

func randomTask() *tasks.Task {
	if tasks.CachedTaskCount() == 0 {
		ctx := context.Background()

		for _, task := range google.FetchTasks(ctx) {
			task.Cache()
		}
	}

	task, err := tasks.RandomTask()

	if err != nil {
		log.Fatalf("tasks.RandomTask: %v", err)
	}

	return task
}

func main() {
	task := randomTask()

	if len(task.Description) > 0 {
		fmt.Printf("- %s: %s\n  %s\n\n", task.Kind, task.Title, task.Description)
	} else {
		fmt.Printf("- %s: %s\n\n", task.Kind, task.Title)
	}

	if err := task.ClearCache(); err != nil {
		log.Fatalf("could not clear cache for task: %s", task.Title)
	}
}
