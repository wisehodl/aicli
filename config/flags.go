package config

import "flag"

type stringSlice []string

func (s *stringSlice) String() string {
	if s == nil {
		return ""
	}
	return ""
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func parseFlags(args []string) (flagValues, error) {
	fv := flagValues{}

	fs := flag.NewFlagSet("aicli", flag.ContinueOnError)
	fs.Usage = printUsage

	var files stringSlice
	var prompts stringSlice

	// Input flags
	fs.Var(&files, "f", "")
	fs.Var(&files, "file", "")
	fs.Var(&prompts, "p", "")
	fs.Var(&prompts, "prompt", "")
	fs.StringVar(&fv.promptFile, "pf", "", "")
	fs.StringVar(&fv.promptFile, "prompt-file", "", "")

	// System flags
	fs.StringVar(&fv.system, "s", "", "")
	fs.StringVar(&fv.system, "system", "", "")
	fs.StringVar(&fv.systemFile, "sf", "", "")
	fs.StringVar(&fv.systemFile, "system-file", "", "")

	// API flags
	fs.StringVar(&fv.key, "k", "", "")
	fs.StringVar(&fv.key, "key", "", "")
	fs.StringVar(&fv.keyFile, "kf", "", "")
	fs.StringVar(&fv.keyFile, "key-file", "", "")
	fs.StringVar(&fv.protocol, "l", "", "")
	fs.StringVar(&fv.protocol, "protocol", "", "")
	fs.StringVar(&fv.url, "u", "", "")
	fs.StringVar(&fv.url, "url", "", "")

	// Model flags
	fs.StringVar(&fv.model, "m", "", "")
	fs.StringVar(&fv.model, "model", "", "")
	fs.StringVar(&fv.fallback, "b", "", "")
	fs.StringVar(&fv.fallback, "fallback", "", "")

	// Output flags
	fs.StringVar(&fv.output, "o", "", "")
	fs.StringVar(&fv.output, "output", "", "")
	fs.StringVar(&fv.config, "c", "", "")
	fs.StringVar(&fv.config, "config", "", "")

	// Boolean flags
	fs.BoolVar(&fv.stdinFile, "F", false, "")
	fs.BoolVar(&fv.stdinFile, "stdin-file", false, "")
	fs.BoolVar(&fv.quiet, "q", false, "")
	fs.BoolVar(&fv.quiet, "quiet", false, "")
	fs.BoolVar(&fv.verbose, "v", false, "")
	fs.BoolVar(&fv.verbose, "verbose", false, "")
	fs.BoolVar(&fv.version, "version", false, "")

	if err := fs.Parse(args); err != nil {
		return flagValues{}, err
	}

	fv.files = files
	fv.prompts = prompts

	return fv, nil
}
