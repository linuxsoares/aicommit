package configprompt

const SystemPrompt = `You are a senior software engineer. Generate a commit message using semantic commit message format based on the summary below.

Use the format: <type>(<scope>): <subject>

Only return the commit message, no explanations.`

const UserPrompt = `Here is the summary of changes:

main.go | 2 +-
1 file changed, 1 insertion(+), 1 deletion(-)

Generate a concise, clear, and well-written commit message using Semantic Commit Messages.`
