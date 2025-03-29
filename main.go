package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Define flags
	countLines := flag.Bool("l", false, "count lines")
	countWords := flag.Bool("w", false, "count words")
	countBytes := flag.Bool("c", false, "count bytes")
	countChars := flag.Bool("m", false, "count characters (overrides -c)")

	flag.Parse()

	// Determine options based on flags
	opts := NewCounterOptions(*countLines, *countWords, *countBytes, *countChars)

	// Get file arguments
	files := flag.Args()

	var totalCounts Counts
	exitCode := 0

	if len(files) == 0 {
		// Read from stdin
		counts, err := processReader(os.Stdin, "<stdin>", opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
			exitCode = 1
		} else {
			printCounts(os.Stdout, counts, opts, "") // No filename for stdin
			totalCounts.Add(counts)                  // Keep track for potential future extensions
		}
	} else {
		// Process each file
		for _, filename := range files {
			file, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %v\n", filepath.Base(os.Args[0]), filename, err)
				exitCode = 1 // Record error, but continue processing other files
				continue     // Skip to the next file
			}

			counts, err := processReader(file, filename, opts)
			file.Close() // Close the file promptly

			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: processing error: %v\n", filepath.Base(os.Args[0]), filename, err)
				exitCode = 1
			} else {
				printCounts(os.Stdout, counts, opts, filename)
				totalCounts.Add(counts)
			}
		}

		// Print total if more than one file was processed successfully (or attempted)
		if len(files) > 1 {
			printCounts(os.Stdout, totalCounts, opts, "total")
		}
	}

	os.Exit(exitCode)
}

// processReader handles the reading and counting for a single source.
func processReader(r io.Reader, sourceName string, opts CounterOptions) (Counts, error) {
	// Pass the direct boolean flags needed by the Count function
	counts, err := Count(r, opts)
	if err != nil {
		return Counts{}, fmt.Errorf("error counting %s: %w", sourceName, err)
	}
	return counts, nil
}

// printCounts formats and prints the counts to the writer.
func printCounts(w io.Writer, counts Counts, opts CounterOptions, label string) {
	// Determine which counts to print based on the options used
	// Use standard formatting similar to wc (right-aligned, fixed width)
	if opts.Lines {
		fmt.Fprintf(w, "%7d", counts.Lines)
	}
	if opts.Words {
		fmt.Fprintf(w, "%8d", counts.Words) // Note: adjusted spacing slightly
	}
	if opts.Chars { // -m takes precedence
		fmt.Fprintf(w, "%8d", counts.Chars)
	} else if opts.Bytes { // Print bytes if -c or default
		fmt.Fprintf(w, "%8d", counts.Bytes)
	}

	// Print label (filename or "total") if provided
	if label != "" {
		fmt.Fprintf(w, " %s", label)
	}
	fmt.Fprintln(w) // Add newline
}
