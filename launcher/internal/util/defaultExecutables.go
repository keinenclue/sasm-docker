package util

import (
	"os/exec"
	"runtime"
)

// GetDockerExecPath returns the default path oft the docker executable
func GetDockerExecPath() string {
	return getExecPath(map[string]string{
		"linux":   "docker",
		"windows": "docker.exe",
		"darwin":  "Docker.app",
	})
}

// GetXserverExecPath returns the default path oft the xserver executable
func GetXserverExecPath() string {
	return getExecPath(map[string]string{
		"linux":   "Xorg",
		"windows": "xlaunch.exe",
		"darwin":  "XQuartz.app",
	})
}

func getExecPath(execNames map[string]string) string {
	path, err := exec.LookPath(execNames[runtime.GOOS])
	if err != nil {
		return ""
	}
	return path
}
