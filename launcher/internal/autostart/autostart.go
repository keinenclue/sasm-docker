package autostart

import (
	"os/exec"
	"time"

	"github.com/keinenclue/sasm-docker/launcher/internal/config"
)

func StartAll() {
	programs := []string{"docker", "xserver"}
	for _, program := range programs {
		if config.Get("autostart." + program + ".enabled").(bool) {
			exe := config.Get("autostart." + program + ".path").(string)
			cmd := exec.Command(exe, "arg")
			cmd.Start()
			time.Sleep(2 * time.Second)
		}
	}
}
