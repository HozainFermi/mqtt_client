package main

import (
	pgs "kurs/gui/pages"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	pgs.MqttClientsPage(app, pages)
	pgs.DbConnectPage(app, pages)

	pages.SwitchToPage("DbConnectPage")

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
