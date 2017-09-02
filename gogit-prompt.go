package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "bytes"
)

func readCommand(cmdName string, cmdArgs string) string {
    split := strings.Split(cmdArgs, " ")
    cmd := exec.Command(cmdName, split...)
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

func checkMaster(gitBranchList string) bool {
    lines := strings.Split(gitBranchList, "\n")
    for i := 0; i < len(lines); i++ {
		if (lines[i] != "" && lines[i][2:] == "master") {
            return true
        }
    }

    return false
}

func countStringsWithPrefixInList(lines []string, prefix string) int64 {
    count := int64(0)
    for _, line := range lines {
        if (strings.HasPrefix(line, prefix)) {
            count += 1
        }
    }

    return count
}

func countCommitDiff(branch string, againstBranch string) (int64, int64) {
    gitCommits := readCommand("git", "rev-list --left-right " + branch + "..." + againstBranch)
    commits := strings.Split(gitCommits, "\n")
    return countStringsWithPrefixInList(commits, "<"), countStringsWithPrefixInList(commits, ">")
}

func main() {
    text := " "
    gitStatus := readCommand("git", "status --porcelain -b")
    if (strings.HasPrefix(gitStatus, "fatal")) {
        os.Exit(0)
    }

    gitStatusLines := strings.Split(gitStatus, "\n")
    branches := strings.Split(gitStatusLines[0][3:], "...")

    currentBranch := branches[0]
    text += currentBranch
    if (len(gitStatusLines) > 2) {
        text += "*"
    }

    remoteBranch := ""
    if (len(branches) > 1) {
        remoteData := strings.Split(branches[1], " ")
        remoteBranch = remoteData[0]
    }

    gitBranchList := readCommand("git", "branch --list master")
    hasMaster := checkMaster(gitBranchList)

    gitCommitsBehindMaster := int64(0)
    gitCommitsAheadMaster := int64(0)
    gitCommitsBehindOrigin := int64(0)
    gitCommitsAheadOrigin := int64(0)

    if (remoteBranch != "") {
        gitCommitsBehindOrigin, gitCommitsAheadOrigin = countCommitDiff(remoteBranch, currentBranch)
        text += fmt.Sprintf(" R[-%s|+%s]",
                            strconv.FormatInt(gitCommitsBehindOrigin, 10),
                            strconv.FormatInt(gitCommitsAheadOrigin, 10))
    }

    if (hasMaster && remoteBranch != "origin/master") {
        gitCommitsBehindMaster, gitCommitsAheadMaster = countCommitDiff("origin/master", currentBranch)
        text += fmt.Sprintf(" M[-%s|+%s]",
                            strconv.FormatInt(gitCommitsBehindMaster, 10),
                            strconv.FormatInt(gitCommitsAheadMaster, 10))
    }

    fmt.Println(text)
}
