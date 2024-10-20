package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/semaphore"
)

type FileResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "dir"
	URL  string `json:"download_url"`
}

func main() {
	startTime := time.Now()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the destination folder from the environment variable
	destinationFolder := os.Getenv("CHERRYPICK_DOWNLOAD_DIR")
	if destinationFolder == "" {
		destinationFolder = "cherrypicked" // Set default folder name
	}

	// Ensure the destination folder exists
	err = os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create destination folder: %v", err)
	}

	// Get the URLs from the command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide GitHub URLs as arguments")
	}

	// Split the URLs by comma
	urls := strings.Split(os.Args[1], ",")

	for _, url := range urls {
		println("Getting", url)
		// Declare variables to hold user, repo, branch, and targetFolder
		var user, repo, branch, targetFolder string

		// Check if the URL points to a file or a directory
		if strings.Contains(url, "/blob/") {
			fetchFileFromURL(url)
		} else {
			var err error
			user, repo, branch, targetFolder, err = parseGitHubURL(url) // Assign values
			if err != nil {
				log.Fatalf("Error parsing URL '%s': %v", url, err)
			}

			fetchDirectoryContents(user, repo, branch, targetFolder, destinationFolder)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Time taken to download the directories: %s\n", elapsedTime)
}

func parseGitHubURL(url string) (user, repo, branch, targetFolder string, err error) {
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return "", "", "", "", fmt.Errorf("invalid URL format. Expected format: https://github.com/user/repo")
	}

	user = parts[3]
	repo = parts[4]

	// Handle optional branch and targetFolder paths
	if len(parts) >= 7 && parts[5] == "tree" {
		branch = parts[6]
		targetFolder = strings.Join(parts[7:], "/")
	} else {
		branch = "main" // Default to main
		targetFolder = strings.Join(parts[5:], "/")
	}

	return user, repo, branch, targetFolder, nil
}

func fetchDirectoryContents(user, repo, branch, targetFolder, destinationFolder string) {
	client := &http.Client{}

	if err := os.MkdirAll(destinationFolder, os.ModePerm); err != nil {
		log.Fatalf("Failed to create destination folder: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", user, repo, targetFolder, branch)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_TOKEN"))

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch files: %s. \nIf this is a private repository enter a token that has access to it.", response.Status)
	}

	var files []FileResponse
	if err := json.NewDecoder(response.Body).Decode(&files); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(50) // Limit concurrent downloads to 50

	for _, file := range files {
		wg.Add(1)
		go func(file FileResponse) {
			defer wg.Done()
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Printf("Failed to acquire semaphore: %v", err)
				return
			}
			defer sem.Release(1)
			if file.Type == "file" {
				fetchFile(user, repo, branch, file.Path, destinationFolder)
			} else if file.Type == "dir" {
				subDirDestination := filepath.Join(destinationFolder, file.Name)
				fetchDirectoryContents(user, repo, branch, file.Path, subDirDestination)
			}
		}(file)
	}

	wg.Wait()
}

func fetchFile(user, repo, branch, path, destinationFolder string) {
	client := &http.Client{}
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", user, repo, branch, path)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(os.Getenv("GITHUB_USERNAME"), os.Getenv("GITHUB_TOKEN"))

	response, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch file %s: %v", path, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch file %s: %s", path, response.Status)
		return
	}

	// Extract the filename from the path
	fileName := filepath.Base(path)
	destinationPath := filepath.Join(destinationFolder, fileName) // Use only the filename for destination

	if err := os.MkdirAll(destinationFolder, os.ModePerm); err != nil {
		log.Printf("Failed to create directory for file %s: %v", path, err)
		return
	}

	file, err := os.Create(destinationPath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", destinationPath, err)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		log.Printf("Failed to write file %s: %v", destinationPath, err)
		return
	}

	log.Printf("Downloaded: %s", destinationPath)
}

func fetchFileFromURL(fileURL string) {
	parts := strings.Split(fileURL, "/")
	if len(parts) < 5 {
		log.Fatal("Invalid file URL format.")
	}
	user := parts[3]
	repo := parts[4]
	branch := "main" // Default to main if not specified
	filePath := strings.Join(parts[7:], "/")

	destinationFolder := "cherrypicked" // Specify destination folder
	fetchFile(user, repo, branch, filePath, destinationFolder)
}
