package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/dialog"

	data "github.com/sarumaj/water-maker/pkg/data"
	handler "github.com/sarumaj/water-maker/pkg/handler"
)

func main() {
	h := handler.NewHandler()

	myApp := app.NewWithID("com.github.sarumaj.domestic.water-maker")
	icon, err := data.Fs.ReadFile("images/icon.png")
	if err != nil {
		panic(err)
	}
	app.SetMetadata(
		fyne.AppMetadata{
			Name:    "Watermaker",
			Build:   1,
			Version: "1.0.0",
			ID:      myApp.Metadata().ID,
			Icon:    fyne.NewStaticResource("icon", icon),
		},
	)

	myWindow := myApp.NewWindow(myApp.Metadata().Name)
	myWindow.Resize(fyne.Size{Width: 800, Height: 600})
	myWindow.SetFixedSize(true)

	resolver := func() {
		if err := recover(); err != nil {
			if val, ok := err.(error); ok {
				dialog.ShowError(val, myWindow)
			} else {
				dialog.ShowError(fmt.Errorf("%v", err), myWindow)
			}
		}
	}
	defer resolver()

	progress := widget.NewProgressBar()
	progress.Resize(fyne.Size{Width: 390, Height: 60})
	progress.Move(fyne.Position{X: 390, Y: 10})

	list := widget.NewListWithData(
		h.Messages,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	list.Resize(fyne.Size{Width: 780, Height: 510})
	list.Move(fyne.Position{X: 10, Y: 80})

	SendMessage := func(msg string) {
		h.Log(msg)
		list.ScrollToBottom()
	}

	browseBtn := widget.NewButton("browse", func() {
		defer resolver()
		dialog.ShowFolderOpen(h.Browse, myWindow)
		list.ScrollToBottom()
	})
	browseBtn.Resize(fyne.Size{Width: 180, Height: 60})
	browseBtn.Move(fyne.Position{X: 10, Y: 10})

	startBtn := widget.NewButton("start", func() {
		defer resolver()
		if length := len(h.Files); length > 0 {
			wg := &sync.WaitGroup{}
			var cnt uint32
			for i, length := 0, len(h.Files); i < length; i++ {
				wg.Add(1)
				go func(i int, wg *sync.WaitGroup) {
					file := h.Files[i]
					SendMessage(fmt.Sprintf("Setting watermark on %s\n", file.Base))
					file.SetWatermark()
					SendMessage(fmt.Sprintf("Set watermark on %s\n", file.Base))
					progress.SetValue(float64(atomic.LoadUint32(&cnt)+1) / float64(length))
					atomic.AddUint32(&cnt, 1)
					wg.Done()
				}(i, wg)
			}
			wg.Wait()
			h.Clear()
			dialog.ShowInformation("Completed", "Done", myWindow)
		} else {
			dialog.ShowInformation("Nothing to do", "Select files first", myWindow)
		}
	})
	startBtn.Resize(fyne.Size{Width: 180, Height: 60})
	startBtn.Move(fyne.Position{X: 200, Y: 10})

	myWindow.SetContent(
		container.NewWithoutLayout(
			browseBtn,
			startBtn,
			progress,
			list,
		),
	)

	myWindow.ShowAndRun()
}
