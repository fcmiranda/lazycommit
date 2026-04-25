package provider

import (
	"fmt"
	"strings"

	"github.com/m7medvision/lazycommit/internal/config"
)

// GetCommitMessagePrompt returns the standardized prompt for generating commit messages
func GetCommitMessagePrompt(diff string) string {
	basePrompt := config.GetCommitMessagePromptFromConfig(diff)
	opts := GetRuntimeCommitPromptOptions()

	var extraRules []string
	if opts.Generate > 0 {
		extraRules = append(extraRules, fmt.Sprintf("IMPORTANT: Generate exactly %d commit messages. Put each message on its own line.", opts.Generate))
	}
	if strings.TrimSpace(opts.Language) != "" {
		extraRules = append(extraRules, fmt.Sprintf("IMPORTANT: Generate all content in %s.", strings.TrimSpace(opts.Language)))
	}

	if len(extraRules) == 0 {
		return basePrompt
	}

	return fmt.Sprintf("%s\n\n%s", basePrompt, strings.Join(extraRules, "\n"))
}

// GetPRTitlePrompt returns the standardized prompt for generating pull request titles
func GetPRTitlePrompt(diff string) string {
	return config.GetPRTitlePromptFromConfig(diff)
}

// GetSystemMessage returns the standardized system message for commit message generation
func GetSystemMessage() string {
	return config.GetSystemMessageFromConfig()
}
