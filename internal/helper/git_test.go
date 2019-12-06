package helper

import (
	"testing"
)

func Test_git_LogFileLastCommit(t *testing.T) {
	filePath := "/Users/liuning/Documents/github/book/https/letsencrypt.md"
	date, err := Git.LogFileLastCommitTime(filePath)
	if err != nil {
		t.Errorf("获取失败, %s", err.Error())
	}
	t.Log(date)
}

func Test_git_LogFileFirstCommitTime(t *testing.T) {
	filePath := "/Users/liuning/Documents/github/book/https/letsencrypt.md"
	date, err := Git.LogFileFirstCommitTime(filePath)
	if err != nil {
		t.Errorf("获取失败, %s", err.Error())
	}
	t.Log(date)
}
