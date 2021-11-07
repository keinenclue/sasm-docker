package main

import (
	"os/user"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/keinenclue/sasm-docker/launcher/internal/config"
	"github.com/keinenclue/sasm-docker/launcher/internal/gui"
)

func main() {
	u, _ := user.Current()
	config.Setup(u.HomeDir + "/sasm-data")
	a := app.New()
	w := a.NewWindow("Sasm-docker launcher")
	w.Resize(fyne.NewSize(500, 300))
	w.SetContent(gui.New(w))

	w.ShowAndRun()
}
