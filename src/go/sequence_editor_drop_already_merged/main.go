package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

/*
Drop any commits specified in the input parameters, keep the others.
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: sequence_editor_drop_already_merged dropCommit1 [dropCommit2...] rebaseFilename")
		os.Exit(1)
	}
	slog.Debug(fmt.Sprint("Got args ", os.Args))
	dropCommits := os.Args[1 : len(os.Args)-1]
	rebaseFilename := os.Args[len(os.Args)-1]

	data, err := os.ReadFile(rebaseFilename)

	if err != nil {
		panic(fmt.Sprint("Could not open ", rebaseFilename, err))
	}

	originalText := string(data)
	var newText strings.Builder

	i := 0
	lines := strings.Split(strings.TrimSuffix(originalText, "\n"), "\n")
	for _, line := range lines {
		if isDropLine(line, dropCommits) {
			dropLine := strings.Replace(line, "pick", "drop", 1)
			newText.WriteString(dropLine)
			newText.WriteString("\n")
			i++
		} else {
			newText.WriteString(line)
			newText.WriteString("\n")
		}
	}

	err = os.WriteFile(rebaseFilename, []byte(newText.String()), 0)
	if err != nil {
		panic(err)
	}
}

func isDropLine(line string, dropCommits []string) bool {
	if !strings.HasPrefix(line, "pick ") {
		return false
	}
	for _, commit := range dropCommits {
		if strings.Index(line, commit) != -1 {
			return true
		}
	}
	return false
}
