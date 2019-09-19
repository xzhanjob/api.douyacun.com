package initialize

import (
	"dyc/internal/helper"
	"github.com/sevlyar/go-daemon"
	"os"
	"syscall"
)

func ChildAlreadyRunning(filePath string) (pid int, running bool) {
	if !helper.FileExists(filePath) {
		return pid, false
	}
	pid, err := daemon.ReadPidFile(filePath)
	if err != nil {
		return pid, false
	}

	process, err := os.FindProcess(int(pid))
	if err != nil {
		return pid, false
	}

	return pid, process.Signal(syscall.Signal(0)) == nil
}
