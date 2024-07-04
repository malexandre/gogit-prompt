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
		currentBranch     string = ""
		changesInProgress string = ""
		remoteBranch      string = ""
		mainBranch        string = ""

		gitCommitsBehindMain   int
		gitCommitsAheadMain    int
		gitCommitsBehindOrigin int
		gitCommitsAheadOrigin  int
	)

	gitStatus := readCommand("git status --porcelain -b")
	// No git in current folder
	if strings.HasPrefix(gitStatus, "fatal") {
		os.Exit(1)
	}

	gitStatusLines := strings.Split(gitStatus, "\n")
	branches := strings.Split(gitStatusLines[0][3:], "...")

	currentBranch = branches[0]
	if len(gitStatusLines) > 2 {
		changesInProgress += "*"
	}

	if len(branches) > 1 {
		remoteBranch = strings.Split(branches[1], " ")[0]
	}

	mainBranch = getMainBarnch()

	prompt := "\ue725 " + currentBranch + changesInProgress
	if remoteBranch != "" {
		gitCommitsBehindOrigin, gitCommitsAheadOrigin = countCommitDiff(remoteBranch, currentBranch)
		prompt += fmt.Sprintf(" \ue726[\uf175%v \uf176%v]", gitCommitsBehindOrigin, gitCommitsAheadOrigin)

		if !strings.Contains(mainBranch, remoteBranch) {
			gitCommitsBehindMain, gitCommitsAheadMain = countCommitDiff(mainBranch, currentBranch)
			prompt += fmt.Sprintf(" \uf09b[\uf175%v \uf176%v]", gitCommitsBehindMain, gitCommitsAheadMain)

		}
	} else {
		gitCommits := readCommand("git log --oneline")
		prompt += fmt.Sprintf(" [\uf1750 \uf176%v]", len(strings.Split(gitCommits, "\n")))
	}

	fmt.Print(prompt)
}
