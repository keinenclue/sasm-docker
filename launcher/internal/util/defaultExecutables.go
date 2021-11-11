package util

import (
	"os/exec"
	"runtime"
	"strings"
)

// GetDockerExecPath returns the default path oft the docker executable
func GetDockerExecPath() string {
	return getExecPath(map[string]string{
		"linux":   "docker",
		"windows": "docker.exe",
		"darwin":  "/Applications/Docker.app",
	})
}

// GetXserverExecPath returns the default path oft the xserver executable
func GetXserverExecPath() string {
	return getExecPath(map[string]string{
		"linux":   "Xorg",
		"windows": "xlaunch.exe",
		"darwin":  "/opt/X11/bin/xquartz",
	})
}

func getExecPath(execNames map[string]string) string {
	execName := execNames[runtime.GOOS]
	if runtime.GOOS == "darwin" && strings.HasSuffix(execName, ".app") {
		return execName
	}

	path, err := exec.LookPath(execName)
	if err != nil {
		return ""
	}
	return path
}
