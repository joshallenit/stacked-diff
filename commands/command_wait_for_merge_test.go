package commands

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	sd "stackeddiff"

	ex "github.com/joshallenit/stacked-diff/execute"
)

func TestSdWaitForMerge_WaitsForMerge(t *testing.T) {
	assert := assert.New(t)

	testExecutor := testutil.InitTest(slog.LevelInfo)

	testutil.AddCommit("first", "")
	allCommits := sd.GetAllCommits()
	testExecutor.SetResponse("2025-01-01", nil, "gh", "pr", "view", ex.MatchAnyRemainingArgs)

	out := testParseArguments("wait-for-merge", allCommits[0].Commit)

	assert.Contains(out, "Merged!")
}
