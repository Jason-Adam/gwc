package main

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
	"unicode/utf8"
)

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

// Count performs the counting operation based on options.
// It reads from the provided reader and returns the counts.
func Count(r io.Reader, opts CounterOptions) (Counts, error) {
	var counts Counts
	wasSpace := true    // Assume starting with whitespace to count the first word correctly
	var byteCount int64 // Use int64 for potentially large files

	// Use bufio.Reader for efficient reading
	// We need to read byte-by-byte to correctly identify words and lines simultaneously.
	// Reading runes directly complicates byte counting if needed alongside char counting.
	// So, we'll read bytes and decode runes only if necessary for char counting.
	br := bufio.NewReader(r)
	buf := make([]byte, 32*1024) // Read in chunks for efficiency

	for {
		n, err := br.Read(buf)
		if n > 0 {
			byteCount += int64(n)
			chunk := buf[:n]

			// Count lines by finding '\n'
			if opts.Lines {
				counts.Lines += bytes.Count(chunk, []byte{'\n'})
			}

			// Count words and optionally characters
			if opts.Words || opts.Chars {
				byteIndex := 0
				for byteIndex < n {
					// If counting chars (-m), decode runes
					if opts.Chars {
						runeValue, runeSize := utf8.DecodeRune(chunk[byteIndex:])
						counts.Chars++ // Increment char count for every rune decoded

						// Word counting logic (based on runes if -m)
						if opts.Words {
							isSpace := unicode.IsSpace(runeValue)
							if !isSpace && wasSpace {
								counts.Words++
							}
							wasSpace = isSpace
						}
						byteIndex += runeSize // Move index by rune size
					} else {
						// If not counting chars, process byte by byte for word counting
						b := chunk[byteIndex]
						if opts.Words {
							// Treat byte as rune for IsSpace check (safe for ASCII whitespace)
							// For full Unicode correctness without -m, this might be slightly off
							// but matches simple `wc` behavior more closely.
							isSpace := unicode.IsSpace(rune(b))
							if !isSpace && wasSpace {
								counts.Words++
							}
							wasSpace = isSpace
						}
						byteIndex++ // Move index by 1 byte
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return Counts{}, err
		}
	}

	// Always assign the total byte count read
	counts.Bytes = int(byteCount) // Note: potential overflow if > MaxInt, but standard wc does this

	return counts, nil
}

// Add accumulates counts from another Counts struct.
func (c *Counts) Add(other Counts) {
	c.Lines += other.Lines
	c.Words += other.Words
	c.Bytes += other.Bytes
	c.Chars += other.Chars
}
