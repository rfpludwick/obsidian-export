package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	flagObsidianVaultDirectory string
	flagOutputFile             string
	flagStartDate              string
	flagEndDate                string
	flagAppend                 bool
)

func initFlags() {
	flag.StringVar(&flagObsidianVaultDirectory, "vault-directory", ".", "Obsidian vault directory to scan")
	flag.StringVar(&flagOutputFile, "output-file", "output.md", "File to write the extracted content")
	flag.StringVar(&flagStartDate, "start-date", "", "Start date for matching level 1 headers (format: YYYY-MM-DD)")
	flag.StringVar(&flagEndDate, "end-date", "", "End date for matching level 1 headers (format: YYYY-MM-DD)")
	flag.BoolVar(&flagAppend, "append", false, "Append to the output file instead of overwriting it")
}

func main() {
	initFlags()
	flag.Parse()

	// Default start and end dates to the current date if not provided
	currentDate := time.Now().Format("2006-01-02")
	if flagStartDate == "" {
		flagStartDate = currentDate
	}
	if flagEndDate == "" {
		flagEndDate = currentDate
	}

	// If not appending, truncate the output file before processing the directory
	if !flagAppend {
		err := os.WriteFile(flagOutputFile, []byte{}, 0644)
		if err != nil {
			fmt.Printf("Error truncating output file: %v\n", err)
			return
		}
	}

	fmt.Printf("Scanning for Markdown files with level 1 headers matching date range: %s to %s\n", flagStartDate, flagEndDate)

	// Traverse the vault directory
	err := filepath.Walk(flagObsidianVaultDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the output file to avoid scanning it
		if filepath.Clean(path) == filepath.Clean(flagOutputFile) {
			return nil
		}

		// Check if the file is a Markdown file
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			processMarkdownFile(path, flagStartDate, flagEndDate)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error while scanning directory: %v\n", err)
	}
}

func processMarkdownFile(filePath, startDate, endDate string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var isMatchingHeader bool
	var outputContent strings.Builder

	for _, line := range lines {
		// Check for level 1 headers
		if strings.HasPrefix(line, "# ") {
			// If a matching header is found, start capturing content
			if isDateInRange(line, startDate, endDate) {
				fmt.Printf("Match found in file %s: %s\n", filePath, line)
				outputContent.WriteString(fmt.Sprintf("# %s\n", strings.TrimSuffix(filePath, ".md")))
				isMatchingHeader = true
			} else {
				// Stop capturing content when a new level 1 header is found
				isMatchingHeader = false
			}
		} else if isMatchingHeader {
			// Append content under the matching header
			outputContent.WriteString(line + "\n")
		}
	}

	// Write the captured content to the output file
	if outputContent.Len() > 0 {
		flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND

		f, err := os.OpenFile(flagOutputFile, flags, 0644)
		if err != nil {
			fmt.Printf("Error opening output file: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.WriteString(outputContent.String() + "\n"); err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
		}
	}
}

func isDateInRange(headerLine, startDate, endDate string) bool {
	// Extract the date from the header line
	headerDate := extractDateFromHeader(headerLine)
	if headerDate == "" {
		return false
	}

	// Compare the header date with the start and end dates
	return headerDate >= startDate && headerDate <= endDate
}

func extractDateFromHeader(headerLine string) string {
	// Define possible date formats
	dateFormats := []string{
		"2006-01-02",              // YYYY-MM-DD
		"January 2, 2006",         // Full month name, day, year
		"Jan 2, 2006",             // Abbreviated month name, day, year
		"Jan 2 2006",              // Abbreviated month name, day, year (no comma)
		"Monday, January 2, 2006", // Full weekday, full month name, day, year
		"02 Jan 2006",             // Day, abbreviated month name, year
		"2 Jan 2006",              // Day (no leading zero), abbreviated month name, year
	}

	// Remove leading pound signs and any extra whitespace
	headerLine = strings.TrimSpace(strings.TrimLeft(headerLine, "#"))

	// Remove day ordinals (e.g., 'st', 'nd', 'rd', 'th') using a regular expression
	ordinalRegex := regexp.MustCompile(`\b(\d+)(st|nd|rd|th)\b`)
	headerLine = ordinalRegex.ReplaceAllString(headerLine, `$1`)

	// Try to parse the entire header line as a date
	for _, format := range dateFormats {
		if parsedDate, err := time.Parse(format, headerLine); err == nil {
			// Return the parsed date in a standard format (YYYY-MM-DD)
			return parsedDate.Format("2006-01-02")
		}
	}

	// If no date is found, return an empty string
	return ""
}
