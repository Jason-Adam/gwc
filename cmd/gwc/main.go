package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	gwc "github.com/Jason-Adam/gwc"
)

func main() {
	// Define flags
	countLines := flag.Bool("l", false, "count lines")
	countWords := flag.Bool("w", false, "count words")
	countBytes := flag.Bool("c", false, "count bytes")
	countChars := flag.Bool("m", false, "count characters (overrides -c)")

	flag.Parse()

	// Determine options based on flags
	opts := gwc.NewCounterOptions(*countLines, *countWords, *countBytes, *countChars)

	// Get file arguments
	files := flag.Args()
	if len(files) == 0 {
		readFromStdin(opts)
	} else {
		readFromFiles(files, opts)
	}

	os.Exit(0)
}

func readFromStdin(opts gwc.CounterOptions) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin:", scanner.Err())
		os.Exit(1)
	}
}

func readFromFiles(files []string, opts gwc.CounterOptions) {
	for _, file := range files {
		file, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if scanner.Err() != nil {
			fmt.Fprintln(os.Stderr, "error reading from file:", scanner.Err())
			os.Exit(1)
		}
	}
}
