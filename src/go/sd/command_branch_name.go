package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	sd "stackeddiff"
)

func createBranchNameCommand(stdOut io.Writer) Command {
	flagSet := flag.NewFlagSet("branch-name", flag.ContinueOnError)
	indicatorTypeString := addIndicatorFlag(flagSet)
	return Command{
		FlagSet:         flagSet,
		DefaultLogLevel: slog.LevelError,
		Summary:         "Outputs branch name of commit",
		Description: "Outputs the branch name for a given commit indicator.\n" +
			"Useful for your own custom scripting.",
		Usage: "sd " + flagSet.Name() + " [flags] <commitIndicator>",
		OnSelected: func(command Command) {
			if flagSet.NArg() == 0 {
				commandError(flagSet, "missing commitIndicator", command.Usage)
			}
			if flagSet.NArg() > 1 {
				commandError(flagSet, "too many arguments", command.Usage)
			}
			indicatorType := checkIndicatorFlag(command, indicatorTypeString)
			branchName := sd.GetBranchInfo(flagSet.Arg(0), indicatorType).BranchName
			fmt.Fprint(stdOut, branchName)
		}}
}
