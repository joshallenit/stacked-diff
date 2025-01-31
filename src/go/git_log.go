package stackeddiff

import (
	"fmt"
	"io"
	"slices"
	"strings"

	ex "stackeddiff/execute"
)

func PrintGitLog(out io.Writer) {
	if GetCurrentBranchName() != ex.GetMainBranch() {
		gitArgs := []string{"--no-pager", "log", "--pretty=oneline", "--abbrev-commit"}
		if RemoteHasBranch(ex.GetMainBranch()) {
			gitArgs = append(gitArgs, "origin/"+ex.GetMainBranch()+"..HEAD")
		}
		gitArgs = append(gitArgs, "--color=always")
		ex.ExecuteOrDie(ex.ExecuteOptions{Output: &ex.ExecutionOutput{Stdout: out, Stderr: out}}, "git", gitArgs...)
		return
	}
	logs := GetNewCommits(ex.GetMainBranch(), "HEAD")
	gitBranchArgs := make([]string, 0, len(logs)+2)
	gitBranchArgs = append(gitBranchArgs, "branch", "-l")
	for _, log := range logs {
		gitBranchArgs = append(gitBranchArgs, log.Branch)
	}
	checkedBranchesRaw := ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", gitBranchArgs...)
	checkedBranches := strings.Split(strings.TrimSpace(checkedBranchesRaw), "\n")
	for i, log := range logs {
		numberPrefix := getNumberPrefix(i, len(logs))
		if slices.Contains(checkedBranches, log.Branch) {
			fmt.Fprint(out, numberPrefix+"✅ ")
		} else {
			fmt.Fprint(out, numberPrefix+"   ")
		}
		fmt.Fprintln(out, ex.Yellow+log.Commit+ex.Reset+" "+log.Subject)
		// find first commit that is not in main branch
		if slices.Contains(checkedBranches, log.Branch) {
			branchCommits := GetNewCommits(ex.GetMainBranch(), log.Branch)
			if len(branchCommits) > 1 {
				for _, branchCommit := range branchCommits {
					padding := strings.Repeat(" ", len(numberPrefix))
					fmt.Fprintln(out, padding+"   - "+branchCommit.Subject)
				}
			}
		}
	}
}

func getNumberPrefix(i int, numLogs int) string {
	maxIndex := fmt.Sprint(numLogs)
	currentIndex := fmt.Sprint(i + 1)
	padding := strings.Repeat(" ", len(maxIndex)-len(currentIndex))
	return padding + currentIndex + ". "
}
