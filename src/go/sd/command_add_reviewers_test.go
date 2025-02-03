package main

import (
	"flag"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	sd "stackeddiff"
	ex "stackeddiff/execute"
	"stackeddiff/testinginit"
)

func TestSdAddReviewers_AddReviewers(t *testing.T) {
	assert := assert.New(t)

	testExecutor := testinginit.InitTest(slog.LevelInfo)

	testinginit.AddCommit("first", "")

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"new"})

	allCommits := sd.GetAllCommits()
	testExecutor.SetResponse(
		// There has to be at least 4 checks, each with 3 values: status, conclusion, and state.
		"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n",
		nil, "gh", "pr", "view", ex.MatchAnyRemainingArgs)

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"add-reviewers", "--reviewers=mybestie", allCommits[0].Commit})

	ghExpectedArgs := []string{"pr", "edit", allCommits[0].Branch, "--add-reviewer", "mybestie"}
	expectedResponse := ex.ExecuteResponse{Out: "Ok", Err: nil, ProgramName: "gh", Args: ghExpectedArgs}
	assert.Contains(testExecutor.Responses, expectedResponse)
}

func TestSdAddReviewers_WhenUsingListIndicator_AddReviewers(t *testing.T) {
	assert := assert.New(t)

	testExecutor := testinginit.InitTest(slog.LevelInfo)

	testinginit.AddCommit("first", "")

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"new"})

	allCommits := sd.GetAllCommits()
	testExecutor.SetResponse(
		// There has to be at least 4 checks, each with 3 values: status, conclusion, and state.
		"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n",
		nil, "gh", "pr", "view", ex.MatchAnyRemainingArgs)

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"add-reviewers", "--indicator=list", "--reviewers=mybestie", "1"})

	ghExpectedArgs := []string{"pr", "edit", allCommits[0].Branch, "--add-reviewer", "mybestie"}
	expectedResponse := ex.ExecuteResponse{Out: "Ok", Err: nil, ProgramName: "gh", Args: ghExpectedArgs}
	assert.Contains(testExecutor.Responses, expectedResponse)
}

func TestSdAddReviewers_WhenOmittingCommitIndicator_UsesHead(t *testing.T) {
	assert := assert.New(t)

	testExecutor := testinginit.InitTest(slog.LevelInfo)

	testinginit.AddCommit("first", "")

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"new"})

	allCommits := sd.GetAllCommits()
	testExecutor.SetResponse(
		// There has to be at least 4 checks, each with 3 values: status, conclusion, and state.
		"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n"+
			"SUCCESS\nSUCCESS\nSUCCESS\n",
		nil, "gh", "pr", "view", ex.MatchAnyRemainingArgs)

	ParseArguments(os.Stdout, flag.NewFlagSet("sd", flag.ContinueOnError), []string{"add-reviewers", "--indicator=list", "--reviewers=mybestie"})

	ghExpectedArgs := []string{"pr", "edit", allCommits[0].Branch, "--add-reviewer", "mybestie"}
	expectedResponse := ex.ExecuteResponse{Out: "Ok", Err: nil, ProgramName: "gh", Args: ghExpectedArgs}
	assert.Contains(testExecutor.Responses, expectedResponse)
}
