package config

var defaultConfig = ConfigData{
	StdinAsFile:    false,
	Protocol:       ProtocolOpenAI,
	URL:            "https://api.ppq.ai/chat/completions",
	Model:          "gpt-4o-mini",
	FallbackModels: []string{"gpt-4.1-mini"},
	Quiet:          false,
	Verbose:        false,
}
