/*
For use as a sequence editor for an interactive git rebase.
Drop all commits.

usage: sequence_editor_drop_all rebaseFilename
*/
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {
	slog.Debug(fmt.Sprint("Got args ", os.Args))
	if len(os.Args) != 2 {
		fmt.Printf("usage: sequence_editor_drop_all rebaseFilename")
		os.Exit(1)
	}
	rebaseFilename := os.Args[1]

	data, err := os.ReadFile(rebaseFilename)

	if err != nil {
		panic(fmt.Sprint("Could not open ", rebaseFilename, err))
	}

	originalText := string(data)
	var newText strings.Builder

	i := 0
	lines := strings.Split(strings.TrimSuffix(originalText, "\n"), "\n")
	for _, line := range lines {
		if isDropLine(line) {
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

func isDropLine(line string) bool {
	return strings.HasPrefix(line, "pick ")
}
