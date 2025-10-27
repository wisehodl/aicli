package input

// StdinRole determines how stdin content participates in the query
type StdinRole int

const (
	// StdinAsPrompt: stdin becomes the entire prompt (replaces other prompts)
	StdinAsPrompt StdinRole = iota

	// StdinAsPrefixedContent: stdin appends after explicit prompts
	StdinAsPrefixedContent

	// StdinAsFile: stdin becomes first file in files array
	StdinAsFile
)

// FileData represents a single input file
type FileData struct {
	Path    string
	Content string
}

// InputData holds all resolved input streams after aggregation
type InputData struct {
	Prompts []string
	Files   []FileData
}
