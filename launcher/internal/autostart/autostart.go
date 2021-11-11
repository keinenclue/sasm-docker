package autostart

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/keinenclue/sasm-docker/launcher/internal/config"
)

// StartAll executes all configured autostarts
func StartAll() {
	programs := []string{"docker", "xserver"}
	for _, program := range programs {
		if config.Get("autostart." + program + ".enabled").(bool) {
			exe := config.Get("autostart." + program + ".path").(string)
			var cmd *exec.Cmd
			if runtime.GOOS == "darwin" && strings.HasSuffix(exe, ".app") {
				cmd = exec.Command("/usr/bin/open", exe)
			} else {
				cmd = exec.Command(exe)
			}
			e := cmd.Start()
			fmt.Println(e)
			time.Sleep(2 * time.Second)
		}
	}
}
