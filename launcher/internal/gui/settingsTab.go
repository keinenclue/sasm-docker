package gui

import (
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/keinenclue/sasm-docker/launcher/internal/config"
)

func newSettingsTab(w fyne.Window) fyne.CanvasObject {

	vBox := container.NewVBox()
	vBox.AddObject(widget.NewLabel("Here you can change some settings :)\nThe settings are stored in: " + config.PathUsed() + "/config.yml"))

	c := newConfigSection(vBox)
	newPathConfig("dataPath", "Data path", w, c, false)
	newSwitchConfig("closeAfterLaunch", "Close after successfull launch", c)
	newSwitchConfig("offlineMode", "Offline mode", c)

	c = newConfigSection(vBox)
	newSwitchConfig("autostart.docker.enabled", "Enable Docker autostart", c)
	newPathConfig("autostart.docker.path", "Docker executable path", w, c, true)

	c = newConfigSection(vBox)
	newSwitchConfig("autostart.xserver.enabled", "Enable Xserver autostart", c)
	newPathConfig("autostart.xserver.path", "Xserver executable path", w, c, true)

	return vBox
}

func newConfigSection(parent *fyne.Container) *fyne.Container {
	parent.AddObject(widget.NewSeparator())
	c := container.New(layout.NewFormLayout())
	parent.AddObject(c)
	return c
}

func newPathConfig(configKey string, label string, w fyne.Window, container *fyne.Container, file bool) {
	var buttonWidget *widget.Button = nil
	buttonWidget = widget.NewButton(config.Get(configKey).(string), func() {
		var fileDialog *dialog.FileDialog = nil
		if !file {
			fileDialog = dialog.NewFolderOpen(func(uc fyne.ListableURI, e error) {
				if e == nil && uc != nil {
					config.Set(configKey, uc.Path())
					buttonWidget.SetText(uc.Path())
				}
			}, w)
		} else {
			fileDialog = dialog.NewFileOpen(func(uc fyne.URIReadCloser, e error) {
				if e == nil && uc != nil {
					config.Set(configKey, uc.URI().Path())
					buttonWidget.SetText(uc.URI().Path())
				}
			}, w)
		}

		var currentPath string = config.Get(configKey).(string)
		if file {
			currentPath = path.Dir(currentPath)
		}
		dir, _ := storage.ListerForURI(storage.NewFileURI(currentPath))

		fileDialog.SetLocation(dir)
		fileDialog.Resize(fyne.NewSize(w.Canvas().Size().Width*0.9, w.Canvas().Size().Height*0.9))
		fileDialog.Show()
	})
	newSettingsRow(label, buttonWidget, container)
}

func newSwitchConfig(configKey string, label string, container *fyne.Container) {
	checkWidget := widget.NewCheck("", func(b bool) {
		config.Set(configKey, b)
	})
	checkWidget.SetChecked(config.Get(configKey).(bool))
	newSettingsRow(label, checkWidget, container)
}

func newSettingsRow(label string, controllerWidget fyne.CanvasObject, container *fyne.Container) {
	labelWidget := widget.NewLabel(label)

	container.AddObject(labelWidget)
	container.AddObject(controllerWidget)
}
