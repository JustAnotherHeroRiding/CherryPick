package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	startTime := time.Now() // Start the timer

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := flag.String("url", "", "GitHub repository URL")
	flag.Parse()

	if *url == "" {
		log.Fatal("Please provide a GitHub repository URL")
	}

	var repoURL, branch, targetFolder string
	re := regexp.MustCompile(`^https:\/\/github\.com\/([^\/]+)\/([^\/]+)(?:\/tree\/([^\/]+))?\/?(.*)`)
	matches := re.FindStringSubmatch(*url)
	fmt.Println("Found Matches: ", matches)

	if len(matches) < 5 {
		log.Fatalf("Invalid URL format: %s", *url)
	}

	repoURL = fmt.Sprintf("https://github.com/%s/%s", matches[1], matches[2])
	fmt.Printf("Attempting to clone: %s\n", repoURL)

	branch = matches[3]
	targetFolder = matches[4]
	destinationFolder := "cherrypicked"

	tmpDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		log.Fatal("Error creating temp directory:", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after

	auth := &http.BasicAuth{
		Username: os.Getenv("GITHUB_USERNAME"),
		Password: os.Getenv("GITHUB_TOKEN"),
	}

	options := &git.CloneOptions{
		URL:        repoURL + ".git",
		Auth:       auth,
		NoCheckout: true,
		Depth:      1,
	}

	repo, err := git.PlainClone(tmpDir, false, options)
	if err != nil {
		log.Fatal("Error cloning repo:", err)
	}

	// Configure sparse checkout
	sparseDir := filepath.Join(tmpDir, ".git", "info")
	if err := os.MkdirAll(sparseDir, os.ModePerm); err != nil {
		log.Fatal("Error creating sparse-checkout directory:", err)
	}

	sparseFilePath := filepath.Join(sparseDir, "sparse-checkout")
	if err := os.WriteFile(sparseFilePath, []byte(targetFolder+"\n"), 0644); err != nil {
		log.Fatal("Error writing sparse-checkout file:", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal("Error getting worktree:", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		log.Fatal("Error checking out branch:", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		log.Fatal("Error getting head reference:", err)
	}

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		log.Fatal("Error getting commit object:", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Fatal("Error getting tree:", err)
	}

	var wg sync.WaitGroup
	tree.Files().ForEach(func(f *object.File) error {
		if strings.HasPrefix(filepath.Clean(f.Name), filepath.Clean(targetFolder)) {
			wg.Add(1)
			go func(file *object.File) {
				defer wg.Done()
				fmt.Println("Copying:", file.Name)
				if err := copyFile(file, destinationFolder, targetFolder); err != nil {
					log.Println("Error copying file:", err)
				}
			}(f)
		}
		return nil
	})

	wg.Wait()
	elapsedTime := time.Since(startTime) // Stop the timer

	fmt.Printf("Files copied to: %s\n", destinationFolder)
	fmt.Printf("Time taken to download the target folder: %s\n", elapsedTime)

}

func copyFile(f *object.File, destinationFolder, targetFolder string) error {
	content, err := f.Contents()
	if err != nil {
		return err
	}

	relativePath := f.Name[len(targetFolder):]
	destinationPath := filepath.Join(destinationFolder, relativePath)

	if err := os.MkdirAll(filepath.Dir(destinationPath), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(destinationPath, []byte(content), 0644)
}
