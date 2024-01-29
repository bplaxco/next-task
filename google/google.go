package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	googletasks "google.golang.org/api/tasks/v1"

	"github.com/bplaxco/next-task/tasks"
)

func FetchTasks(ctx context.Context) []*tasks.Task {
	var fetchedTasks []*tasks.Task

	fetchedTasks = append(fetchedTasks, tasksFromGmail(ctx)...)
	fetchedTasks = append(fetchedTasks, tasksFromTasks(ctx)...)

	return fetchedTasks
}

func tasksFromGmail(ctx context.Context) []*tasks.Task {
	var fetchedTasks []*tasks.Task
	user := "me"

	svc, err := gmail.NewService(ctx, option.WithHTTPClient(getClient(ctx)))
	if err != nil {
		log.Fatalf("unable to retrieve Gmail client: %v", err)
	}

	messagesList, err := svc.Users.Messages.List(user).Do()
	if err != nil {
		log.Fatalf("unable to retrieve messages: %v", err)
	}

	for _, m := range messagesList.Messages {
		m, err = svc.Users.Messages.Get(user, m.Id).Do()

		if err != nil {
			log.Fatalf("Unable to retrieve labels: %v", err)
		}

		for _, h := range m.Payload.Headers {
			if h.Name == "Subject" {
				title := fmt.Sprintf("Process: %s", h.Value)

				if tasks.TaskAlreadyExists(fetchedTasks, title) {
					break
				}

				fetchedTasks = append(fetchedTasks, tasks.NewTask("GoogleMail", m.Id, title, m.Snippet))
				break
			}
		}
	}

	return fetchedTasks
}

func tasksFromTasks(ctx context.Context) []*tasks.Task {
	var fetchedTasks []*tasks.Task

	svc, err := googletasks.NewService(ctx, option.WithHTTPClient(getClient(ctx)))
	if err != nil {
		log.Fatalf("unable to retrieve tasks client %v", err)
	}

	taskLists, err := svc.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("unable to retrieve task lists. %v", err)
	}

	for _, taskList := range taskLists.Items {
		tasksList, err := svc.Tasks.List(taskList.Id).ShowCompleted(false).ShowDeleted(false).Do()

		if err != nil {
			log.Fatalf("unable to retrieve tasks list. %v", err)
		}

		for _, task := range tasksList.Items {
			if tasks.TaskAlreadyExists(fetchedTasks, task.Title) {
				continue
			}

			fetchedTasks = append(fetchedTasks, tasks.NewTask("GoogleTask", task.Id, task.Title, ""))
		}
	}

	return fetchedTasks
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf(
		"Go to the following link in your browser then type the authorization code: \n%v\nCode: ",
		authURL,
	)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(ctx, authCode)

	if err != nil {
		log.Fatalf("unable to retrieve token from web: %v", err)
	}

	return token
}

func loadGoogleConfig() *oauth2.Config {
	// This credentials.json should be for an OAuth2 app that can read the resources outlined in the scopes requested below
	credsFile, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config/next-task/google/credentials.json"))

	if err != nil {
		log.Fatalf("unable to load client credentials file", err)
	}

	config, err := google.ConfigFromJSON(
		credsFile,
		gmail.GmailReadonlyScope,
		googletasks.TasksReadonlyScope,
	)

	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v", err)
	}

	return config
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("unable to cache oauth token: %v", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func googleOAuth2Token(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	tokenFilePath := filepath.Join(os.Getenv("HOME"), ".config/next-task/google/token.json")
	token, err := tokenFromFile(tokenFilePath)

	if err != nil {
		token = tokenFromWeb(ctx, config)
		saveToken(tokenFilePath, token)
	}

	return token, nil
}

func getClient(ctx context.Context) *http.Client {
	config := loadGoogleConfig()

	token, err := googleOAuth2Token(ctx, config)

	if err != nil {
		log.Fatalf("unable to load Google OAuth2 token: %w", err)
	}

	return config.Client(ctx, token)
}
