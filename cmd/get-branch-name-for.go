package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Outputs the branch name for a given commit hash or pull request number. Useful for custom scripting.")
		fmt.Println("get-branch-name-for <commit hash or pull request number>")
		os.Exit(1)
	}
	log.SetOutput(ioutil.Discard)
	branchName := GetBranchInfo(os.Args[1]).BranchName
	fmt.Print(branchName)
}
