package helper

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var Git _git

type _git struct{}

func (*_git) LogFileLastCommitTime(filePath string) (t time.Time, err error) {
	if !Git.CheckGitExists() {
		return time.Time{}, errors.New("git 命令不存在")
	}
	dir := path.Dir(filePath)
	file := path.Base(filePath)
	command := fmt.Sprintf("cd %s && git log --format='%%ct' -- %s | awk 'NR==1'", dir, file)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return
	}
	timestamp, err := strconv.ParseInt(strings.TrimRight(string(out), "\n"), 10, 64)
	if err != nil {
		return
	}
	t = time.Unix(timestamp, 0)
	return
}

func (*_git) LogFileFirstCommitTime(filePath string) (t time.Time, err error) {
	if !Git.CheckGitExists() {
		return time.Time{}, errors.New("git 命令不存在")
	}
	dir := path.Dir(filePath)
	file := path.Base(filePath)
	command := fmt.Sprintf("cd %s && git log --reverse --format='%%ct' -- %s |awk 'NR == 1'", dir, file)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return
	}
	timestamp, err := strconv.ParseInt(strings.TrimRight(string(out), "\n"), 10, 64)
	if err != nil {
		return
	}
	t = time.Unix(timestamp, 0)
	return
}

func (*_git) CheckGitExists() bool {
	_, err := exec.LookPath("ls")
	return err == nil
}
