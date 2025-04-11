package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	// Open the current repository
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Could not open repository: %v", err)
	}

	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Could not get worktree: %v", err)
	}

	// Get the status of the working directory
	status, err := w.Status()
	if err != nil {
		log.Fatalf("Could not get status: %v", err)
	}

	// Generate a list of changed files
	var changedFiles []string
	for file, _ := range status {
		changedFiles = append(changedFiles, file)
	}

	// Generate a text about the changes
	changeText := generateChangeText(changedFiles)
	fmt.Println("Changes detected:")
	fmt.Println(changeText)

	// Create a simple commit message
	commitMessage := generateCommitMessage(changedFiles)
	fmt.Println("Commit message:")
	fmt.Println(commitMessage)

	// Commit the changes
	commitChanges(w, commitMessage)
}

func generateChangeText(files []string) string {
	return fmt.Sprintf("The following files have been changed:\n%s", strings.Join(files, "\n"))
}

func generateCommitMessage(files []string) string {
	return fmt.Sprintf("Updated %d files", len(files))
}

func commitChanges(w *git.Worktree, message string) {
	// Add all changes
	_, err := w.Add(".")
	if err != nil {
		log.Fatalf("Could not add changes: %v", err)
	}

	// Commit the changes
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your.email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Could not commit changes: %v", err)
	}

	fmt.Printf("Committed changes with hash: %s\n", commit.String())
}
