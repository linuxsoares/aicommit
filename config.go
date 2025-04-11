package main

import (
	"os"
)

var AICOMMAND_OPEN_AI_TOKEN string

func init() {
	AICOMMAND_OPEN_AI_TOKEN = os.Getenv("OPENAI_API_KEY")
}
