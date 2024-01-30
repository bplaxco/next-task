package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	googlemail "google.golang.org/api/gmail/v1"
	googletasks "google.golang.org/api/tasks/v1"
)

type Config struct {
	Google *Google
	Jira   *Jira
}

func NewConfig(ctx context.Context) *Config {
	return &Config{
		Jira:   newJiraConfig(),
		Google: newGoogleConfig(ctx),
	}
}

type Google struct {
	OAuth2      *oauth2.Config
	OAuth2Token *oauth2.Token
}

func newGoogleConfig(ctx context.Context) *Google {
	credentailsFilePath := filepath.Join(os.Getenv("HOME"), ".config/next-task/google/credentials.json")

	if _, err := os.Stat(credentailsFilePath); err != nil {
		return nil
	}
	oauth2Config := newOAuth2Config(credentailsFilePath)

	return &Google{
		OAuth2:      oauth2Config,
		OAuth2Token: newOAuth2Token(ctx, oauth2Config),
	}
}

func newOAuth2Token(ctx context.Context, oauth2Config *oauth2.Config) *oauth2.Token {
	tokenFilePath := filepath.Join(os.Getenv("HOME"), ".config/next-task/google/token.json")
	token, err := newOAuth2TokenFromFile(tokenFilePath)

	if err != nil {
		token = newOAuth2TokenFromWeb(ctx, oauth2Config)
		saveOAuth2Token(tokenFilePath, token)
	}

	return token
}

func newOAuth2TokenFromFile(path string) (*oauth2.Token, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	return token, err
}

func newOAuth2TokenFromWeb(ctx context.Context, oauth2Config *oauth2.Config) *oauth2.Token {
	authURL := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf(
		"Go to the following link in your browser then type the authorization code: \n%v\nCode: ",
		authURL,
	)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("unable to read authorization code: %v", err)
	}

	token, err := oauth2Config.Exchange(ctx, authCode)

	if err != nil {
		log.Fatalf("unable to retrieve token from web: %v", err)
	}

	return token
}

func newOAuth2Config(credentailsFilePath string) *oauth2.Config {
	credentailsFile, err := os.ReadFile(credentailsFilePath)

	if err != nil {
		log.Fatalf("unable to load client credentials file", err)
	}

	oauth2Config, err := google.ConfigFromJSON(
		credentailsFile,
		googlemail.GmailReadonlyScope,
		googletasks.TasksReadonlyScope,
	)

	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v", err)
	}

	return oauth2Config
}

func saveOAuth2Token(tokenFilePath string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", tokenFilePath)
	f, err := os.OpenFile(tokenFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalf("unable to cache oauth token: %v", err)
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

type Jira struct {
	InstanceURL string
	AccessToken string
	TaskJQL     string
}

func newJiraConfig() *Jira {
	jiraInstanceURL := os.Getenv("NEXT_TASK_JIRA_INSTANCE_URL")

	if len(jiraInstanceURL) > 0 {
		return &Jira{
			InstanceURL: jiraInstanceURL,
			AccessToken: os.Getenv("NEXT_TASK_JIRA_ACCESS_TOKEN"),
			TaskJQL:     os.Getenv("NEXT_TASK_JIRA_TASK_JQL"),
		}
	}

	return nil
}
