package gui

import (
	"fmt"
	"math"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/keinenclue/sasm-docker/launcher/internal/autostart"
	"github.com/keinenclue/sasm-docker/launcher/internal/config"
	c "github.com/keinenclue/sasm-docker/launcher/internal/container"
)

var statusLabel *widget.Label = nil
var statusLabelSameMessageCount = 0

func newLaunchTab(a fyne.App, w fyne.Window) fyne.CanvasObject {
	var hellos string

	if runtime.GOOS == "darwin" {
		hellos = `XQuartz`
	} else {
		hellos = `XServer`
	}

	hello := widget.NewLabel(fmt.Sprintf(`Welcome to the sasm docker launcher!
	Before you start, make sure to have these tools installed:
	- Docker
	- %s
	
	Once you have them installed, just click Launch :D`, hellos))

	statusLabel = widget.NewLabel("")
	imageSelector := widget.NewSelect(c.AvailableImages(), nil)
	imageSelector.SetSelected("test")
	imageSelector.PlaceHolder = "(select an image)"

	startButton := widget.NewButton("Launch!", nil)
	imageSelector.OnChanged = handleContainerSelectionChanged(startButton)

	if imageSelector.SelectedIndex() == -1 {
		startButton.Disable()
	}

	vBox := container.NewVBox(
		hello,
		layout.NewSpacer(),
		statusLabel,
		imageSelector,
		startButton,
	)

	startButton.OnTapped = launchImage(a, startButton, imageSelector, vBox)
	return vBox
}

func launchAppendLog(level string, message string) {
	if level != "CONTAINER" {
		if statusLabel.Text == message || strings.HasSuffix(statusLabel.Text, message) {
			statusLabelSameMessageCount++
			message = fmt.Sprintf("(%d) %s", statusLabelSameMessageCount, message)
		} else {
			statusLabelSameMessageCount = 0
		}
		statusLabel.SetText(message)
	}
	appendLog(level, message)
}

func handleContainerSelectionChanged(startButton *widget.Button) func(string) {
	return func(selectedImage string) {
		fmt.Println(selectedImage)
		startButton.Enable()
	}
}

func handleContainerEvent(app fyne.App, layerProgress map[string]*widget.ProgressBar, vBox *fyne.Container, startButton *widget.Button, imageSelector *widget.Select) func(event c.Event) {
	return func(event c.Event) {
		switch event.Type {
		case c.ImagePullStatusChanged:
			s := event.Data.(c.ImagePullStatus)
			//fmt.Println(s)
			if s.ID == "" || s.ProgressDetail == nil {
				launchAppendLog("INFO", s.Status)
				return
			}

			if layerProgress[s.ID] == nil {
				pb := widget.NewProgressBar()
				layerProgress[s.ID] = pb
				pb.SetValue(0)
				vBox.AddObject(pb)
				vBox.Refresh()
			}
			//fmt.Printf("%s %s", s.Status, s.Progress)
			//
			pb := layerProgress[s.ID]
			currentMB := float64(s.ProgressDetail.Current) / 1000000
			totalMB := float64(s.ProgressDetail.Total) / 1000000
			value := currentMB / totalMB
			if math.IsNaN(value) {
				value = 1
			}
			pb.SetValue(value)
			pb.TextFormatter = func() string {
				if value == 1 {
					return fmt.Sprintf("%s: %s", s.ID, s.Status)
				}
				return fmt.Sprintf("%s: %s (%.2fMB/%.2fMB)", s.ID, s.Status, currentMB, totalMB)
			}

			pb.Refresh()

		case c.StateChanged:
			state := event.Data.(c.State)

			launchAppendLog("INFO", fmt.Sprintf("Container state is now: %s", event.Self.StateString()))

			switch state {
			case c.OfflineState:
				handleEndOfSession(startButton, imageSelector)
			case c.RunningState:
				if config.Get("closeAfterLaunch").(bool) {
					app.Quit()
				}
			default:
				startButton.Hide()
			}
			if state != c.PullingState && len(layerProgress) > 0 {
				for _, bar := range layerProgress {
					bar.Hide()
				}
				layerProgress = map[string]*widget.ProgressBar{}
			}

		case c.ConsoleOutput:
			launchAppendLog("CONTAINER", event.Data.(string))
		case c.LogMessage:
			launchAppendLog("LOG", event.Data.(string))
		case c.ErrorMessage:
			launchAppendLog("ERROR", event.Data.(string))
		}
	}
}

func launchImage(app fyne.App, startButton *widget.Button, imageSelector *widget.Select, vBox *fyne.Container) func() {
	return func() {
		go func() {
			startButton.Hide()
			imageSelector.Hide()
			newLogSession()
			statusLabel.SetText("")

			launchAppendLog("INFO", "Starting autostart programs if configured ...")
			autostart.StartAll()

			cont, err := c.NewSasmContainer(imageSelector.Selected)

			if err != nil {
				launchAppendLog("ERROR", err.Error())
				tabs.SelectTabIndex(1)
				handleEndOfSession(startButton, imageSelector)
				return
			}

			layerProgress := make(map[string]*widget.ProgressBar)
			cont.OnContainerEvent(handleContainerEvent(app, layerProgress, vBox, startButton, imageSelector))
			err = cont.Launch()

			if err != nil {
				launchAppendLog("ERROR", err.Error())
				tabs.SelectTabIndex(1)
				handleEndOfSession(startButton, imageSelector)
			}
		}()
	}
}

func handleEndOfSession(startButton *widget.Button, imageSelector *widget.Select) {
	endLogsSession()
	statusLabel.SetText("Sasm exited")
	startButton.Show()
	imageSelector.Show()
}
