package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bplaxco/next-task/config"
	"github.com/bplaxco/next-task/google"
	"github.com/bplaxco/next-task/jira"
	"github.com/bplaxco/next-task/tasks"
)

type ParsedArgs struct {
	Reset bool
	Help  bool
}

const Usage string = `
NAME
       next-task - snatch a next task from a source

SYNOPSIS
       next-task [--help|--reset]

OPTIONS
       --help          Show this help text
       --reset         Reset the cached content

`

func ParseArgs(args []string) *ParsedArgs {
	parsedArgs := &ParsedArgs{}

	for _, arg := range args[1:] {
		switch arg {
		case "--reset":
			parsedArgs.Reset = true
		case "--help":
			parsedArgs.Help = true
		default:
			fmt.Printf("%s is not a valid argument\n", arg)
			parsedArgs.Help = true
			break
		}
	}

	return parsedArgs
}

func randomTask(ctx context.Context, cfg *config.Config) *tasks.Task {
	// The order of things pulled here will mater given
	// the capacity system. In the future it might make
	// sense to pull max capacity of each list and shuffle
	// the full list of tasks and then only cache
	// up to the total capacity. It wastes network
	// calls but ensures we're pulling from each
	// source fairly evenly
	if tasks.CachedTaskCount() == 0 {
		if cfg.Google != nil {
			for _, task := range google.FetchTasks(ctx, cfg.Google) {
				task.Cache()
			}
		}

		if cfg.Jira != nil {
			for _, task := range jira.FetchTasks(cfg.Jira) {
				task.Cache()
			}
		}
	}

	task, err := tasks.RandomTask()

	if err != nil {
		log.Fatalf("tasks.RandomTask: %v", err)
	}

	return task
}

func main() {
	parsedArgs := ParseArgs(os.Args)

	if parsedArgs.Help {
		fmt.Print(Usage)
		return
	}

	if parsedArgs.Reset {
		err := os.RemoveAll(tasks.TaskCacheDir())

		if err != nil {
			fmt.Printf("os.RemoveAll: %s", err.Error())
		}

		return
	}

	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	task := randomTask(ctx, cfg)

	if len(task.Description) > 0 {
		fmt.Printf("- %s: %s\n  %s\n\n", task.Kind, task.Title, task.Description)
	} else {
		fmt.Printf("- %s: %s\n\n", task.Kind, task.Title)
	}

	if err := task.ClearCache(); err != nil {
		log.Fatalf("could not clear cache for task: %s", task.Title)
	}
}
