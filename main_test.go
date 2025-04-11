package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGenerateChangeText tests the generateChangeText function
func TestGenerateChangeText(t *testing.T) {
	// Initialize a new in-memory repository
	repo, _ := git.Init(memory.NewStorage(), memfs.New())

	// Create a new file and commit it
	w, _ := repo.Worktree()
	fs := w.Filesystem
	file, _ := fs.Create("testfile.txt")
	file.Write([]byte("Hello, World!"))
	w.Add("testfile.txt")
	w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})

	// Modify the file
	file, _ = fs.OpenFile("testfile.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	file.Write([]byte("Hello, Go!"))
	w.Add("testfile.txt")

	// Generate change text
	changedFiles := []string{"testfile.txt"}
	changeText, err := generateChangeText(repo, changedFiles)

	// Assert no error and check the change text
	assert.NoError(t, err)
	assert.Contains(t, changeText, "File: testfile.txt")
	assert.Contains(t, changeText, "The following changes have been made")
}

// TestGetDiff tests the getDiff function
func TestGetDiff(t *testing.T) {
	// Initialize a new in-memory repository
	repo, _ := git.Init(memory.NewStorage(), memfs.New())

	// Create a new file and commit it
	w, _ := repo.Worktree()
	fs := w.Filesystem
	file, _ := fs.Create("testfile.txt")
	file.Write([]byte("Hello, World!"))
	w.Add("testfile.txt")
	w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})

	// Modify the file
	file, _ = fs.OpenFile("testfile.txt", os.O_WRONLY|os.O_TRUNC, 0644)
	file.Write([]byte("Hello, Go!"))
	w.Add("testfile.txt")

	// Get the diff
	_, err := getDiff(repo, "testfile.txt")

	// Assert no error and check the patch
	assert.NoError(t, err)
}

// TestGenerateCommitMessageWithOpenAI tests the generateCommitMessageWithOpenAI function
// func TestGenerateCommitMessageWithOpenAI(t *testing.T) {
// 	// Mock the OpenAI client
// 	client := &MockOpenAIClient{}
// 	client.On("CreateChatCompletion", mock.Anything, mock.Anything).Return(&openai.ChatCompletionResponse{
// 		Choices: []openai.ChatCompletionChoice{
// 			{
// 				Message: openai.ChatCompletionMessage{
// 					Content: "Test commit message",
// 				},
// 			},
// 		},
// 	}, nil)

// 	// Generate commit message
// 	changeText := "The following changes have been made:\nFile: testfile.txt\n"
// 	commitMessage, err := generateCommitMessageWithOpenAI(changeText)

// 	// Assert no error and check the commit message
// 	assert.NoError(t, err)
// 	assert.Contains(t, commitMessage, "fix(main):")
// }

// MockOpenAIClient is a mock implementation of the OpenAI client
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (*openai.ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*openai.ChatCompletionResponse), args.Error(1)
}
