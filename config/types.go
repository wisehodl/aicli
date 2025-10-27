package config

type APIProtocol int

const (
	ProtocolOpenAI APIProtocol = iota
	ProtocolOllama
)

type ConfigData struct {
	// Input
	FilePaths   []string
	PromptFlags []string
	PromptPaths []string
	StdinAsFile bool

	// System
	SystemPrompt string

	// API
	Protocol APIProtocol
	URL      string
	APIKey   string

	// Models
	Model          string
	FallbackModels []string

	// Output
	Output  string
	Quiet   bool
	Verbose bool
}

type flagValues struct {
	files      []string
	prompts    []string
	promptFile string
	system     string
	systemFile string
	key        string
	keyFile    string
	protocol   string
	url        string
	model      string
	fallback   string
	output     string
	config     string
	stdinFile  bool
	quiet      bool
	verbose    bool
	version    bool
}

type envValues struct {
	protocol string
	url      string
	key      string
	model    string
	fallback string
	system   string
}

type fileValues struct {
	protocol   string
	url        string
	keyFile    string
	model      string
	fallback   string
	systemFile string
}
