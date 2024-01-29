package tasks

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

var taskCacheDir string

type Task struct {
	Kind        string `json:"kind"`
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewTask(kind, id, title, description string) *Task {
	return &Task{
		Kind:        kind,
		Id:          id,
		Title:       title,
		Description: description,
	}
}

func (t *Task) CacheKey() string {
	h := sha256.New()
	h.Write([]byte(t.Kind))
	h.Write([]byte(":"))
	h.Write([]byte(t.Id))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (t *Task) CachePath() string {
	return filepath.Join(TaskCacheDir(), t.CacheKey())
}

func (t *Task) IsCached() bool {
	_, err := os.Stat(t.CachePath())
	return err == nil
}

func (t *Task) Cache() (string, error) {
	path := t.CachePath()

	data, err := json.Marshal(t)
	if err != nil {
		return path, err
	}

	return path, os.WriteFile(path, data, 0600)
}

func (t *Task) ClearCache() error {
	return os.Remove(t.CachePath())
}

func TaskCacheDir() string {
	if len(taskCacheDir) == 0 {
		taskCacheDir = filepath.Join(os.Getenv("HOME"), ".cache/next-task/tasks")

		if _, err := os.Stat(taskCacheDir); err != nil {
			os.MkdirAll(taskCacheDir, 0700)
		}
	}

	return taskCacheDir
}

func CachedTaskCount() int {
	entries, err := cachedTaskPaths()

	if err != nil {
		return 0
	}

	return len(entries)
}

func RandomTask() (*Task, error) {
	var task Task

	paths, err := cachedTaskPaths()
	if err != nil {
		return &task, fmt.Errorf("cachedTaskPaths: %w", err)
	}

	// Heh, crypto/rand is overkill but I'm working on an idea where
	// I wanted things to be statistically spread out across the set
	// of tasks
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(paths))))

	if err != nil {
		return &task, fmt.Errorf("rand.Int: %w", err)
	}

	path := paths[i.Int64()]

	file, err := os.Open(path)
	if err != nil {
		return &task, fmt.Errorf("os.Open: %w", err)
	}

	err = json.NewDecoder(file).Decode(&task)
	return &task, err
}

func TaskAlreadyExists(currentTasks []*Task, newTaskTitle string) bool {
	for _, task := range currentTasks {
		if task.Title == newTaskTitle {
			return true
		}
	}

	return false
}

func cachedTaskPaths() ([]string, error) {
	var paths []string
	cacheDirPath := TaskCacheDir()

	entries, err := os.ReadDir(cacheDirPath)

	if err != nil {
		return paths, fmt.Errorf("os.ReadDir: %w", err)
	}

	for _, entry := range entries {
		paths = append(paths, filepath.Join(cacheDirPath, entry.Name()))
	}

	return paths, err
}
