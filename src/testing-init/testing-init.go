package testing_init

import (
	"log"
	"os"
	"path"
	"runtime"
	ex "stacked-diff-workflow/src/execute"
)

var TestWorkingDir string
var thisFile string

func init() {
	_, file, _, ok := runtime.Caller(0)
	thisFile = file
	if !ok {
		panic("No caller information")
	}
	TestWorkingDir = path.Join(path.Dir(file), "/../../../.test-stacked-diff-workflow")
}

func CdTestDir() {
	var functionName string
	for i := 0; i < 10; i++ {
		pc, file, _, ok := runtime.Caller(i)
		if !ok {
			panic("No caller information")
		}
		if file != thisFile {
			functionName = runtime.FuncForPC(pc).Name()
			break
		}
	}
	if functionName == "" {
		panic("Could not find caller outside of " + thisFile)
	}
	individualTestDir := TestWorkingDir + "/" + functionName
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "rm", "-rf", individualTestDir)
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "mkdir", "-p", individualTestDir)
	os.Chdir(individualTestDir)
	log.Println("Changed to test directory: " + individualTestDir)
}

func CdTestRepo() {
	CdTestDir()
	// Create a git repository with a local remote
	remoteDir := "remote-repo"
	repositoryDir := "local-repo"

	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "init", "--bare", remoteDir)
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "clone", remoteDir, repositoryDir)

	os.Chdir(repositoryDir)
	log.Println("Changed to repository directory: " + repositoryDir)
}

func AddCommit(commitMessage string, fileName string) {
	if fileName == "" {
		fileName = commitMessage
	}
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "touch", commitMessage)
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "add", ".")
	ex.ExecuteOrDie(ex.ExecuteOptions{}, "git", "commit", "-m", commitMessage)
}
