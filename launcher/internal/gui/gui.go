package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	c "github.com/keinenclue/sasm-docker/launcher/internal/container"
)

var currentTabIndex = 0
var tabs *container.AppTabs = nil

// New creates the gui
func New(w fyne.Window) fyne.CanvasObject {
	c.New()

	tabs = container.NewAppTabs(
		container.NewTabItem("Launch", newLaunchTab(w)),
		container.NewTabItem("Logs", newLogTab(w)),
		container.NewTabItem("Settings", newSettingsTab(w)),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.Resize(w.Canvas().Size())

	return tabs
}
