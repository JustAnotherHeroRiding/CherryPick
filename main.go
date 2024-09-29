package main

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	// Define the repository URL and directory to clone
	repoURL := "https://github.com/JustAnotherHeroRiding/CherryPick.git"
	directory := "./cloned-repo"

	// Clone the repository into the specified directory
	fmt.Println("Cloning the repository...")
	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println("Error during clone:", err)
		return
	}

	// Optionally checkout a specific branch or tag
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting worktree:", err)
		return
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"), // Replace "main" with the desired branch
	})

	if err != nil {
		fmt.Println("Error during checkout:", err)
		return
	}

	fmt.Println("Repository successfully cloned and checked out!")
}
