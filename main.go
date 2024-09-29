package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/joho/godotenv"
)

func main() {
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

	// Clone the repository to a temporary directory (non-bare)
	tmpDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		fmt.Println("Error creating temp directory:", err)
		return
	}
	defer os.RemoveAll(tmpDir) // Clean up after

	token := os.Getenv("GITHUB_TOKEN")
	username := os.Getenv("GITHUB_USERNAME")

	auth := &http.BasicAuth{
		Username: username,
		Password: token,
	}

	options := &git.CloneOptions{
		URL:      repoURL + ".git",
		Auth:     auth,
		Progress: os.Stdout,
	}

	repo, err := git.PlainClone(tmpDir, false, options)

	if err != nil {
		fmt.Println("Error cloning repo:", err)
		return
	}

	// Get the working tree from the repository
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting worktree:", err)
		return
	}

	// Checkout the desired branch (if it's not main, modify this line)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		fmt.Println("Error checking out branch:", err)
		return
	}

	// Get the head commit
	headRef, err := repo.Head()
	if err != nil {
		fmt.Println("Error getting head reference:", err)
		return
	}

	// Get the commit object
	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		fmt.Println("Error getting commit object:", err)
		return
	}

	// Get the tree from the commit
	tree, err := commit.Tree()
	if err != nil {
		fmt.Println("Error getting tree:", err)
		return
	}

	// Traverse the tree to find files inside the target folder
	tree.Files().ForEach(func(f *object.File) error {
		// Check if the file is within the target folder
		if strings.HasPrefix(filepath.Clean(f.Name), filepath.Clean(targetFolder)) {
			fmt.Println("Copying:", f.Name)

			// Read the file content
			content, err := f.Contents()
			if err != nil {
				return err
			}

			// Build the destination file path
			relativePath := f.Name[len(targetFolder):] // Remove target folder prefix
			destinationPath := filepath.Join(destinationFolder, relativePath)

			// Ensure the destination folder exists
			err = os.MkdirAll(filepath.Dir(destinationPath), os.ModePerm)
			if err != nil {
				return err
			}

			// Write the file content to the new folder
			err = os.WriteFile(destinationPath, []byte(content), 0644)
			if err != nil {
				return err
			}
		}
		return nil
	})

	fmt.Println("Files copied to:", destinationFolder)
}
