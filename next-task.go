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
		// The order of things pulled here will mater given
		// the capacity system. In the future it might make
		// sense to pull max capacity of each list and shuffle
		// the full list of tasks and then only cache
		// up to the total capacity. It wastes network
		// calls but ensures we're pulling from each
		// source fairly evenly
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
