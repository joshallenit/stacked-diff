package commands

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"fmt"

	"github.com/fatih/color"
	ex "github.com/joshallenit/stacked-diff/v2/execute"
	"github.com/joshallenit/stacked-diff/v2/templates"
	"github.com/joshallenit/stacked-diff/v2/util"

	"time"
)

func createUpdateCommand() Command {
	flagSet := flag.NewFlagSet("update", flag.ContinueOnError)
	indicatorTypeString := addIndicatorFlag(flagSet)
	reviewers, silent, minChecks := addReviewersFlags(flagSet, "")
	return Command{
		FlagSet: flagSet,
		Summary: "Add commits from " + util.GetMainBranchForHelp() + " to an existing PR",
		Description: "Add commits from local " + util.GetMainBranchForHelp() + " branch to an existing PR.\n" +
			"\n" +
			"Can also add reviewers once PR checks have passed, see \"--reviewers\" flag.",
		Usage: "sd " + flagSet.Name() + " [flags] <PR commitIndicator> [fixup commitIndicator (defaults to head commit) [fixup commitIndicator...]]",
		OnSelected: func(command Command) {
			if flagSet.NArg() == 0 {
				commandError(flagSet, "missing commitIndicator", command.Usage)
			}
			indicatorType := checkIndicatorFlag(command, indicatorTypeString)
			var otherCommits []string
			if len(flagSet.Args()) > 1 {
				otherCommits = flagSet.Args()[1:]
			}
			destCommit := templates.GetBranchInfo(flagSet.Arg(0), indicatorType)
			updatePr(destCommit, otherCommits, indicatorType)
			if *reviewers != "" {
				addReviewersToPr([]string{destCommit.Commit}, templates.IndicatorTypeCommit, true, *silent, *minChecks, *reviewers, 30*time.Second)
			}
		}}
}

// Add commits from main to an existing PR.
func updatePr(destCommit templates.GitLog, otherCommits []string, indicatorType templates.IndicatorType) {
	util.RequireMainBranch()
	templates.RequireCommitOnMain(destCommit.Commit)
	var commitsToCherryPick []string
	if len(otherCommits) > 0 {
		if indicatorType == templates.IndicatorTypeGuess || indicatorType == templates.IndicatorTypeList {
			commitsToCherryPick = util.MapSlice(otherCommits, func(commit string) string {
				return templates.GetBranchInfo(commit, indicatorType).Commit
			})
		} else {
			commitsToCherryPick = otherCommits
		}
	} else {
		commitsToCherryPick = make([]string, 1)
		commitsToCherryPick[0] = strings.TrimSpace(ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "rev-parse", "--short", "HEAD"))
	}
	shouldPopStash := false
	stashResult := strings.TrimSpace(ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "stash", "save", "-u", "before update-pr "+destCommit.Commit))
	if strings.HasPrefix(stashResult, "Saved working") {
		slog.Info(stashResult)
		shouldPopStash = true
	}
	slog.Info(fmt.Sprint("Switching to branch ", destCommit.Branch))
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "switch", destCommit.Branch)
	slog.Info("Fast forwarding in case there were any commits made via github web interface")
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "fetch", "origin", destCommit.Branch)
	forcePush := false
	if _, err := ex.Execute(ex.ExecuteOptions{}, "git", "merge", "--ff-only", "origin/"+destCommit.Branch); err != nil {
		slog.Info(fmt.Sprint("Could not fast forward to match origin. Rebasing instead. ", err))
		ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "rebase", "origin", destCommit.Branch)
		// As we rebased, a force push may be required.
		forcePush = true
	}

	slog.Info(fmt.Sprint("Cherry picking ", commitsToCherryPick))
	cherryPickArgs := make([]string, 1+len(commitsToCherryPick))
	cherryPickArgs[0] = "cherry-pick"
	for i, commit := range commitsToCherryPick {
		cherryPickArgs[i+1] = commit
	}
	_, cherryPickError := ex.Execute(ex.ExecuteOptions{}, "git", cherryPickArgs...)
	if cherryPickError != nil {
		slog.Info("First attempt at cherry-pick failed")
		ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "cherry-pick", "--abort")
		rebaseCommit := util.FirstOriginMainCommit(util.GetMainBranchOrDie())
		slog.Info(fmt.Sprint("Rebasing with the base commit on "+util.GetMainBranchOrDie()+" branch, ", rebaseCommit,
			", in case the local "+util.GetMainBranchOrDie()+" was rebased with origin/"+util.GetMainBranchOrDie()))
		rebaseOutput, rebaseError := ex.Execute(ex.ExecuteOptions{}, "git", "rebase", rebaseCommit)
		if rebaseError != nil {
			slog.Info(fmt.Sprint(color.RedString("Could not rebase, aborting... "), rebaseOutput))
			ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "rebase", "--abort")
			ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "switch", util.GetMainBranchOrDie())
			util.PopStash(shouldPopStash)
			os.Exit(1)
		}
		slog.Info(fmt.Sprint("Cherry picking again ", commitsToCherryPick))
		var cherryPickOutput string
		cherryPickOutput, cherryPickError = ex.Execute(ex.ExecuteOptions{}, "git", cherryPickArgs...)
		if cherryPickError != nil {
			slog.Info(fmt.Sprint(color.RedString("Could not cherry-pick, aborting... "), cherryPickArgs, " ", cherryPickOutput, " ", cherryPickError))
			ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "cherry-pick", "--abort")
			ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "switch", util.GetMainBranchOrDie())
			util.PopStash(shouldPopStash)
			os.Exit(1)
		}
		forcePush = true
	}
	slog.Info("Pushing to remote")
	if forcePush {
		if _, err := ex.Execute(ex.ExecuteOptions{}, "git", "push", "origin", destCommit.Branch); err != nil {
			slog.Info("Regular push failed, force pushing instead.")
			ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "push", "-f", "origin", destCommit.Branch)
		}
	} else {
		ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "push", "origin", destCommit.Branch)
	}
	slog.Info("Switching back to " + util.GetMainBranchOrDie())
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "switch", util.GetMainBranchOrDie())
	slog.Info(fmt.Sprint("Rebasing, marking as fixup ", commitsToCherryPick, " for target ", destCommit.Commit))
	environmentVariables := []string{
		"GIT_SEQUENCE_EDITOR=sequence_editor_mark_as_fixup " +
			destCommit.Commit + " " +
			strings.Join(commitsToCherryPick, " "),
	}
	slog.Debug(fmt.Sprint("Using sequence editor ", environmentVariables))
	options := ex.ExecuteOptions{EnvironmentVariables: environmentVariables, Output: ex.NewStandardOutput()}
	rootCommit := strings.TrimSpace(ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "log", "--max-parents=0", "--format=%h", "HEAD"))
	if rootCommit == destCommit.Commit {
		slog.Info("Rebasing root commit")
		ex.ExecuteOrDie(options, "git", "rebase", "-i", "--root")
	} else {
		ex.ExecuteOrDie(options, "git", "rebase", "-i", destCommit.Commit+"^")
	}
	util.PopStash(shouldPopStash)
}
