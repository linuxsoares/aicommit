package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	configprompt "github.com/linuxsoares/aicommand/configPrompt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/tcnksm/go-gitconfig"
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
	for file := range status {
		changedFiles = append(changedFiles, file)
	}

	// Generate a text about the changes
	changeText, err := generateChangeText(repo, changedFiles)
	if err != nil {
		log.Fatalf("Could not generate change text: %v", err)
	}

	// Generate a commit message using OpenAI
	commitMessage, err := generateCommitMessageWithOpenAI(changeText)
	if err != nil {
		log.Fatalf("Could not generate commit message: %v", err)
	}
	fmt.Println("Generated commit message:")
	fmt.Println(commitMessage)

	// Ask the user if the commit message is okay
	var userResponse string
	fmt.Println("Is this commit message okay? (yes/no)")
	fmt.Scanln(&userResponse)

	if strings.ToLower(userResponse) == "yes" {
		// Commit the changes
		commitChanges(w, commitMessage)
	} else {
		fmt.Println("Commit aborted by user.")
	}
}

func generateChangeText(repo *git.Repository, files []string) (string, error) {
	var changeText strings.Builder
	changeText.WriteString("The following changes have been made:\n")

	for _, file := range files {
		changeText.WriteString(fmt.Sprintf("File: %s\n", file))

		// Get the diff for the file
		patch, err := getDiff(repo, file)
		if err != nil {
			return "", err
		}

		changeText.WriteString(patch)
		changeText.WriteString("\n")
	}

	return changeText.String(), nil
}

func getDiff(repo *git.Repository, file string) (string, error) {
	// Get the HEAD reference
	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}

	// Get the commit object for the HEAD reference
	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return "", err
	}

	// Get the file content at HEAD
	headFile, err := headCommit.File(file)
	var headContent string
	if err == nil {
		headContent, err = headFile.Contents()
		if err != nil {
			return "", err
		}
	}

	// Get the current file content
	currentFile, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	currentFileContent, err := currentFile.Filesystem.Open(file)
	if err != nil {
		return "", err
	}
	defer currentFileContent.Close()

	currentContent, err := io.ReadAll(currentFileContent)
	if err != nil {
		return "", err
	}

	// Generate the diff
	dmp := diffmatchpatch.New()
	var diffs []diffmatchpatch.Diff
	if headContent == "" {
		diffs = dmp.DiffMain("", string(currentContent), false)
	} else {
		diffs = dmp.DiffMain(headContent, string(currentContent), false)
	}
	patch := dmp.DiffPrettyText(diffs)

	if len(patch) > MAXDIFFSIZE {
		patch = patch[:MAXDIFFSIZE] + "\n...diff truncated...\n"
	}
	return patch, nil
}

func generateCommitMessageWithOpenAI(changeText string) (string, error) {
	client := openai.NewClient(AICOMMAND_OPEN_AI_TOKEN)
	ctx := context.Background()

	prompt := fmt.Sprintf(configprompt.UserPrompt, changeText)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: configprompt.SystemPrompt,
			},
		},
		MaxTokens: 1000,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

func commitChanges(w *git.Worktree, message string) {
	// Get author information from git config
	authorName, err := gitconfig.Username()
	if err != nil {
		log.Fatalf("Could not get author name from git config: %v", err)
	}

	authorEmail, err := gitconfig.Email()
	if err != nil {
		log.Fatalf("Could not get author email from git config: %v", err)
	}

	// Add all changes
	_, err = w.Add(".")
	if err != nil {
		log.Fatalf("Could not add changes: %v", err)
	}

	// Commit the changes
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Could not commit changes: %v", err)
	}

	fmt.Printf("Committed changes with hash: %s\n", commit.String())
}
