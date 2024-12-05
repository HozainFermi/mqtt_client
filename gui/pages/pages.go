package pages

import (
	client "kurs/client"

	"github.com/rivo/tview"
)

var username string = "postgres"
var password string = "1357902479"
var dbname string = "Devices"

var clientID string = ""
var clientusername string = ""
var clientpassword string = ""

func DbConnectPage(app *tview.Application, pages *tview.Pages) {

	form := tview.NewForm().
		AddInputField("User name", username, 20, nil, func(text string) { username = text }).
		AddPasswordField("Password", password, 20, '*', func(text string) { password = text }).
		AddInputField("Database name", dbname, 20, nil, func(text string) { dbname = text }).
		AddButton("Save", func() { saveHandler(app, pages) }).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(false), 0, 2, false).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 2, false)

	pages.AddPage("DbConnectPage", flex, true, false)

}

func saveHandler(app *tview.Application, pages *tview.Pages) {
	go func() {
		client.Init(username, password, dbname)
	}()
	pages.SwitchToPage("MqttClientsPage")

}
func connectHandler(app *tview.Application, pages *tview.Pages) {

	go func() {
		client.ConnectMqtt(clientID, clientusername, clientpassword)
	}()

	go app.QueueUpdate(func() { pages.SwitchToPage("Monitoring") })
}

//

func MqttClientsPage(app *tview.Application, pages *tview.Pages) {

	ConnectAsForm := tview.NewForm().
		AddInputField("Client ID", clientID, 20, nil, func(text string) { clientID = text }).
		AddInputField("User name", clientusername, 20, nil, func(text string) { clientusername = text }).
		AddInputField("Password", clientpassword, 20, nil, func(text string) { clientpassword = text }).
		AddButton("Connect", func() { connectHandler(app, pages) })

	ConnectAsForm.SetBorder(true).SetTitle("Connect as:").SetTitleAlign(tview.AlignLeft)

	AddNewUsers := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextArea().SetPlaceholder("sername,password,clientid,topic,pubsub,accses").SetBorder(true).SetTitle("Add new users"), 0, 30, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewButton("Add"), 0, 1, false).
			AddItem(tview.NewButton("Save csv"), 0, 1, false), 0, 2, false)

	leftmenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ConnectAsForm, 0, 1, false).
		AddItem(AddNewUsers, 0, 1, false)

	leftmenu.SetTitle("Menu:")

	grid := tview.NewGrid().
		SetRows(100).
		SetColumns(40, 50).
		SetBorders(true).
		AddItem(leftmenu, 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewBox().SetBorder(true), 0, 1, 1, 3, 0, 0, false)

	pages.AddPage("MqttClientsPage", grid, true, false)

}
