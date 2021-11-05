package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/keinenclue/sasm-docker/launcher/internal/config"
)

func newSettingsTab(w fyne.Window) fyne.CanvasObject {
	vBox := container.NewVBox()

	vBox.AddObject(widget.NewButton("Select storage folder", func() {
		d := dialog.NewFolderOpen(func(uc fyne.ListableURI, e error) {
			if e == nil && uc != nil {
				config.Set("dataPath", uc.Path())
			}
		}, w)

		dir, _ := storage.ListerForURI(storage.NewFileURI(config.Get("dataPath").(string)))

		d.SetLocation(dir)
		d.Resize(fyne.NewSize(w.Canvas().Size().Width*0.9, w.Canvas().Size().Height*0.9))
		d.Show()
	}))

	return vBox
}
