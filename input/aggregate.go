package input

// AggregatePrompts combines prompt sources with stdin based on role.
func AggregatePrompts(prompts []string, stdin string, role StdinRole) []string {
	switch role {
	case StdinAsPrompt:
		if stdin != "" {
			return []string{stdin}
		}
		return prompts

	case StdinAsPrefixedContent:
		if stdin != "" {
			return append(prompts, stdin)
		}
		return prompts

	case StdinAsFile:
		return prompts

	default:
		return prompts
	}
}

// AggregateFiles combines file sources with stdin based on role.
func AggregateFiles(files []FileData, stdin string, role StdinRole) []FileData {
	if role == StdinAsFile && stdin != "" {
		return append([]FileData{{Path: "input", Content: stdin}}, files...)
	}
	return files
}
