package swan

// Config determines how the extractor functions
type Config struct {
	Error struct {
		// Return an error when a language could not be detected rather than
		// defaulting to English
		OnNoLanguage bool
	}
}
