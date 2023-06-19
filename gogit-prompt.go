package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func readCommand(cmdName string) string {
	split := strings.Split(cmdName, " ")
	cmd := exec.Command(split[0], split[1:]...)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return stderr.String()
	}

	return out.String()
}

func getMainBarnch() string {
	return readCommand("git symbolic-ref refs/remotes/origin/HEAD")
}

func countStringsWithPrefixInList(lines []string, prefix string) (count int) {
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			count += 1
		}
	}

	return
}

func countCommitDiff(branch string, againstBranch string) (int, int) {
	gitCommits := readCommand("git rev-list --left-right " + branch + "..." + againstBranch)
	commits := strings.Split(gitCommits, "\n")
	return countStringsWithPrefixInList(commits, "<"), countStringsWithPrefixInList(commits, ">")
}

func main() {
	var (
		text          string = ""
		currentBranch string = ""
		remoteBranch  string = ""
		mainBranch    string = ""

		gitCommitsBehindMain   int
		gitCommitsAheadMain    int
		gitCommitsBehindOrigin int
		gitCommitsAheadOrigin  int
	)

	gitStatus := readCommand("git status --porcelain -b")
	if strings.HasPrefix(gitStatus, "fatal") {
		os.Exit(1)
	}

	gitStatusLines := strings.Split(gitStatus, "\n")
	branches := strings.Split(gitStatusLines[0][3:], "...")

	currentBranch = branches[0]
	text += currentBranch
	if len(gitStatusLines) > 2 {
		text += "*"
	}

	if len(branches) > 1 {
		remoteBranch = strings.Split(branches[1], " ")[0]
	}

	mainBranch = getMainBarnch()

	if remoteBranch != "" {
		gitCommitsBehindOrigin, gitCommitsAheadOrigin = countCommitDiff(remoteBranch, currentBranch)
		text += fmt.Sprintf(" R[-%v|+%v]", gitCommitsBehindOrigin, gitCommitsAheadOrigin)

		if !strings.Contains(mainBranch, remoteBranch) {
			gitCommitsBehindMain, gitCommitsAheadMain = countCommitDiff(mainBranch, currentBranch)
			text += fmt.Sprintf(" M[-%v|+%v]", gitCommitsBehindMain, gitCommitsAheadMain)

		}
	}

	fmt.Print(text)
}
