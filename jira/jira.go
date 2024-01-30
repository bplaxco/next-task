package jira

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/bplaxco/next-task/config"
	"github.com/bplaxco/next-task/tasks"
)

func FetchTasks(cfg *config.Jira) []*tasks.Task {
	var fetchedTasks []*tasks.Task
	client := &Jira{cfg: cfg}

	for _, issue := range client.Search() {
		if tasks.Capacity() == 0 {
			break
		}

		fetchedTasks = append(fetchedTasks, tasks.NewTask("Jira", issue.Key, issue.Fields.Summary, issue.Fields.Description))
		tasks.DecrementCapacity()
	}

	return fetchedTasks
}

type Jira struct {
	cfg *config.Jira
}

type JiraSearchResults struct {
	Issues []*JiraIssue `json:"issues"`
}

type JiraIssue struct {
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

type JiraIssueFields struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func (j *Jira) searchURL() string {
	return fmt.Sprintf(
		"%s/rest/api/latest/search?jql=%s",
		j.cfg.InstanceURL,
		url.QueryEscape(j.cfg.TaskJQL),
	)
}

func (j *Jira) Search() []*JiraIssue {
	var results JiraSearchResults
	req, err := http.NewRequest("GET", j.searchURL(), nil)

	if err != nil {
		log.Printf("could not look up jira tasks: %v", err)
		return results.Issues
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", j.cfg.AccessToken))

	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Printf("could not look up jira tasks: %v", err)
		return results.Issues
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&results)

	if err != nil {
		log.Printf("could not look up jira tasks: %v", err)
	}

	return results.Issues
}
