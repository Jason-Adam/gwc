package gwc

// Counts holds the results of a wc operation.
type Counts struct {
	Lines int
	Words int
	Bytes int
	Chars int
}

// CounterOptions specifies which counts to perform.
type CounterOptions struct {
	Lines bool
	Words bool
	Bytes bool
	Chars bool // -m flag
}

// NewCounterOptions creates options based on flags.
// If no specific flags are true, it defaults to Lines, Words, and Bytes.
func NewCounterOptions(l, w, c, m bool) CounterOptions {
	// Default behavior: if no flags set, count lines, words, bytes
	if !l && !w && !c && !m {
		return CounterOptions{Lines: true, Words: true, Bytes: true, Chars: false}
	}
	// If -m is set, it overrides -c for the primary "byte/char" count display
	// We still count bytes internally if needed, but Chars takes precedence for options.
	if m {
		c = false // Prioritize char count if -m is specified
	}
	return CounterOptions{Lines: l, Words: w, Bytes: c, Chars: m}
}
