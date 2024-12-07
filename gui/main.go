package main

import (
	pgs "kurs/gui/pages"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	pgs.DbConnectPage(app, pages)
	pgs.MqttClientsPage(app, pages)
	pgs.MonitoringPage(app, pages)

	pages.SwitchToPage("DbConnectPage")

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
