# AI-Powered Git Commit Message Generator

This is a CLI tool that automatically generates meaningful Git commit messages using OpenAI's GPT-4 Turbo model. It analyzes the changes in your repository, creates a summary of the differences, and prompts OpenAI to generate a descriptive commit message. You can then approve or reject the message before the commit is performed.

## Features

- Detects modified files in the current Git repository.
- Computes human-readable diffs for each changed file.
- Uses OpenAI to generate a commit message based on file changes.
- Asks for user confirmation before committing.
- Commits changes with Git author info from global config.

## Dependencies

This project is written in Go and uses the following packages:
- [go-git/go-git](https://github.com/go-git/go-git) — to interact with the Git repository.
- [go-gitconfig](https://github.com/tcnksm/go-gitconfig) — to retrieve Git author name and email.
- [go-diff/diffmatchpatch](https://github.com/sergi/go-diff) — to compute diffs between file versions.
- [sashabaranov/go-openai](https://github.com/sashabaranov/go-openai) — to interact with the OpenAI API.
- [`linuxsoares/aicommand/configPrompt`](https://github.com/linuxsoares/aicommand) — contains custom prompts for OpenAI.


## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/linuxsoares/aicommit.git
   cd yourproject
   ```
2. **Install Go modules**:

   ```bash
   go mod tidy
    ```
   
3. **Set your OpenAI API key**:
   Export your OpenAI token as an environment variable:

   ```bash
export AICOMMAND_OPEN_AI_TOKEN=your-openai-api-key
    ```
   
## Usage
Simply run the Go program from the root of a Git repository:
```bash 
go run main.go
```

You will be prompted with a commit message suggestion. If you approve it by typing yes, the tool will stage and commit the changes.

## Example

Generated commit message:

```bash
Refactor API endpoint to improve error handling and add logging
Is this commit message okay? (yes/no)
yes
Committed changes with hash: 1a2b3c4d5e...
```

## Notes

Make sure you have uncommitted changes in your Git repository before running this tool.
Diffs are truncated if they exceed a certain length to avoid hitting OpenAI token limits.
You can customize the prompt logic in the configPrompt package.

## License
[MIT License](https://rem.mit-license.org/)
