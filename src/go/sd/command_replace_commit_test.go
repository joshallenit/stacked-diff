package main

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	sd "stackeddiff"
	ex "stackeddiff/execute"
	"stackeddiff/testinginit"
)

func TestSdReplaceCommit_WithMultipleCommits_ReplacesCommitWithBranch(t *testing.T) {
	assert := assert.New(t)
	testinginit.InitTest(slog.LevelInfo)

	testinginit.AddCommit("first", "1")
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "push", "origin", sd.GetMainBranchOrDie())
	testinginit.AddCommit("second", "will-be-replaced")
	testParseArguments("new")
	testinginit.AddCommit("fifth", "5")

	allCommits := sd.GetAllCommits()

	testParseArguments("checkout", allCommits[1].Commit)

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "reset", "--hard", allCommits[2].Commit)
	testinginit.AddCommit("on-second-branch-only", "2")
	testinginit.AddCommit("on-second-branch-only", "3")
	testinginit.AddCommit("on-second-branch-only", "4")

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "checkout", sd.GetMainBranchOrDie())

	testParseArguments("replace-commit", allCommits[1].Commit)

	allCommits = sd.GetAllCommits()

	assert.Equal(4, len(allCommits))
	assert.Equal("fifth", allCommits[0].Subject)
	assert.Equal("second", allCommits[1].Subject)
	assert.Equal("first", allCommits[2].Subject)
	assert.Equal(testinginit.InitialCommitSubject, allCommits[3].Subject)

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		panic("Could not read dir: " + err.Error())
	}
	assert.Equal(6, len(dirEntries))
	assert.Equal(".git", dirEntries[0].Name())
	assert.Equal("1", dirEntries[1].Name())
	assert.Equal("2", dirEntries[2].Name())
	assert.Equal("3", dirEntries[3].Name())
	assert.Equal("4", dirEntries[4].Name())
	assert.Equal("5", dirEntries[5].Name())
}

func TestSdReplaceCommit_WhenPrPushed_ReplacesCommitWithBranch(t *testing.T) {
	if true {
		return
	}
	assert := assert.New(t)
	testinginit.InitTest(slog.LevelInfo)

	testinginit.AddCommit("first", "1")
	testinginit.AddCommit("second", "will-be-replaced")
	testParseArguments("new")
	testinginit.AddCommit("fifth", "5")

	allCommits := sd.GetAllCommits()

	testParseArguments("checkout", allCommits[1].Commit)

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "reset", "--hard", allCommits[2].Commit)
	testinginit.AddCommit("on-second-branch-only", "2")
	testinginit.AddCommit("on-second-branch-only", "3")
	testinginit.AddCommit("on-second-branch-only", "4")

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "push", "origin", allCommits[1].Branch)

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "checkout", sd.GetMainBranchOrDie())

	testParseArguments("replace-commit", allCommits[1].Commit)

	allCommits = sd.GetAllCommits()

	assert.Equal(3, len(allCommits))
	assert.Equal("fifth", allCommits[0].Subject)
	assert.Equal("second", allCommits[1].Subject)
	assert.Equal("first", allCommits[2].Subject)

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		panic("Could not read dir: " + err.Error())
	}
	assert.Equal(6, len(dirEntries))
	assert.Equal(".git", dirEntries[0].Name())
	assert.Equal("1", dirEntries[1].Name())
	assert.Equal("2", dirEntries[2].Name())
	assert.Equal("3", dirEntries[3].Name())
	assert.Equal("4", dirEntries[4].Name())
	assert.Equal("5", dirEntries[5].Name())
}
