package main

// PathParameters provides path input parameters for the CLI arguments
type PathParameters struct {
	path               string // Content source path
	recursive          bool   // Recursive path search
	fromLang           string // Original language code
	toLang             string // Target language code
	destination        string // Target translate content directory
	frontMatterTargets string // Frontmatter target keys
	googleAuthJSON     string // Google translate API credentials JSON file
}

// FrontMatter provides the set of data from the blog Front Matter to target and translate
type FrontMatter struct {
	Name        string   `yaml:"name" toml:"name" json:"name"`
	Description string   `yaml:"description,omitempty" toml:"description,omitempty" json:"description,omitempty"`
	Bio         string   `yaml:"bio,omitempty" toml:"bio,omitempty" json:"bio,omitempty"`
	Tags        []string `yaml:"tags,omitempty" toml:"tags,omitempty" json:"tags,omitempty"`
}
