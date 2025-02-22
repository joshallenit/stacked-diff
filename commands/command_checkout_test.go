package commands

import (
	"log/slog"
	"testing"

	"github.com/joshallenit/stacked-diff/templates"
	"github.com/joshallenit/stacked-diff/testutil"
	"github.com/joshallenit/stacked-diff/util"
	"github.com/stretchr/testify/assert"
)

func TestSdCheckout_ChecksOutBranch(t *testing.T) {
	assert := assert.New(t)

	testutil.InitTest(slog.LevelInfo)

	testutil.AddCommit("first", "")

	allCommits := templates.GetAllCommits()

	testParseArguments("new")

	testParseArguments("checkout", allCommits[0].Commit)

	assert.Equal(allCommits[0].Branch, util.GetCurrentBranchName())
}
