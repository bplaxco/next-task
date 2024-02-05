package google

import (
	"context"
	"log"
	"net/http"
	"time"

	googlemail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	googletasks "google.golang.org/api/tasks/v1"

	"github.com/bplaxco/next-task/config"
	"github.com/bplaxco/next-task/tasks"
)

func FetchTasks(ctx context.Context, cfg *config.Google) []*tasks.Task {
	var fetchedTasks []*tasks.Task

	// The order here matters given that these burn through
	// capacity. Email should be top since those items can be
	// turned into other tasks if not dealt with quickly.
	fetchedTasks = append(fetchedTasks, tasksFromGmail(ctx, cfg)...)
	fetchedTasks = append(fetchedTasks, tasksFromTasks(ctx, cfg)...)

	return fetchedTasks
}

func tasksFromGmail(ctx context.Context, cfg *config.Google) []*tasks.Task {
	var fetchedTasks []*tasks.Task
	user := "me"

	svc, err := googlemail.NewService(ctx, option.WithHTTPClient(getClient(ctx, cfg)))
	if err != nil {
		log.Fatalf("unable to retrieve Gmail client: %v", err)
	}

	messagesList, err := svc.Users.Messages.List(user).LabelIds("INBOX").MaxResults(tasks.Capacity()).Do()
	if err != nil {
		log.Fatalf("unable to retrieve messages: %v", err)
	}

	for _, m := range messagesList.Messages {
		if tasks.Capacity() == 0 {
			log.Println("skipping Gmail messages because capacity has been reached")
			break
		}

		m, err = svc.Users.Messages.Get(user, m.Id).Do()

		if err != nil {
			log.Fatalf("Unable to retrieve message: %v", err)
		}

		for _, h := range m.Payload.Headers {
			if h.Name == "Subject" {
				if tasks.TaskAlreadyExists(fetchedTasks, h.Value) {
					break
				}

				fetchedTasks = append(fetchedTasks, tasks.NewTask("GoogleMail", m.Id, h.Value, m.Snippet))
				tasks.DecrementCapacity()
				break
			}
		}
	}

	return fetchedTasks
}

func tasksFromTasks(ctx context.Context, cfg *config.Google) []*tasks.Task {
	var fetchedTasks []*tasks.Task

	svc, err := googletasks.NewService(ctx, option.WithHTTPClient(getClient(ctx, cfg)))
	if err != nil {
		log.Fatalf("unable to retrieve Tasks client %v", err)
	}

	taskLists, err := svc.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("unable to retrieve task lists. %v", err)
	}

	for _, taskList := range taskLists.Items {
		if tasks.Capacity() == 0 {
			log.Println("skipping Google Tasks because capacity has been reached")
			break
		}

		tomorrow := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		tasksList, err := svc.Tasks.List(taskList.Id).MaxResults(tasks.Capacity()).DueMax(tomorrow).ShowCompleted(false).ShowDeleted(false).Do()

		if err != nil {
			log.Fatalf("unable to retrieve tasks list. %v", err)
		}

		for _, task := range tasksList.Items {
			if tasks.TaskAlreadyExists(fetchedTasks, task.Title) {
				continue
			}

			fetchedTasks = append(fetchedTasks, tasks.NewTask("GoogleTask", task.Id, task.Title, ""))
			tasks.DecrementCapacity()
		}
	}

	return fetchedTasks
}

func getClient(ctx context.Context, cfg *config.Google) *http.Client {
	return cfg.OAuth2.Client(ctx, cfg.OAuth2Token)
}
