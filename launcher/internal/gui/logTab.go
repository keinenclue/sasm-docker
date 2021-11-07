package gui

import (
	"fmt"
	"math"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

var logContent *widget.TextGrid = nil
var logSessionRunning = false
var noLogs = true

func newLogTab(w fyne.Window) fyne.CanvasObject {
	copyButton := widget.NewButton("Copy to clipboard", func() {
		clipboard.Write(clipboard.FmtText, []byte(logContent.Text()))
	})
	reportIssueButton := widget.NewButton("Report issue", func() {
		openURL("https://github.com/keinenclue/sasm-docker/issues/new")
	})
	logContent = widget.NewTextGridFromString("No logs so far.")
	scroll := container.NewScroll(logContent)

	bottomRow := container.NewHBox(
		copyButton,
		reportIssueButton,
	)

	return container.New(
		layout.NewBorderLayout(nil, bottomRow, nil, nil),
		bottomRow,
		scroll,
	)
}

func newLogSession() {
	logSessionRunning = true
	appendLog("", "---- start ----\n")
}

func endLogsSession() {
	appendLog("", "\n---- end ----\n")
	logSessionRunning = false
}

func appendLog(level string, message string) {
	if !logSessionRunning {
		return
	}

	if noLogs {
		logContent.SetText("")
		noLogs = false
	}

	line := fmt.Sprintf("[%s] %s", level, message)
	if level == "" {
		line = message
	}

	cells := make([]widget.TextGridCell, 0, len(line))
	for _, r := range line {
		cells = append(cells, widget.TextGridCell{Rune: r})
		if r == '\t' {
			col := len(cells)
			tabStop, _ := math.Modf(float64(col-1+logContent.TabWidth) / float64(logContent.TabWidth))
			next := logContent.TabWidth * int(tabStop)
			for i := col; i < next; i++ {
				cells = append(cells, widget.TextGridCell{Rune: ' '})
			}
		}
	}

	logContent.SetRow(
		len(logContent.Rows),
		widget.TextGridRow{Cells: cells},
	)
	logContent.SetText(logContent.Text() + "")
}

// openURL opens the specified URL in the default browser of the user.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
