package pages

import (
	"kurs/client"
	"kurs/systeminfo"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var username string = "postgres"
var password string = "1357902479"
var dbname string = "Devices"

var clientID string = ""
var clientusername string = ""
var clientpassword string = ""

var users []client.User

//var info *tview.TextArea

var str string = ""
var counter int = 0
var disconnectflag bool = false

const refreshInterval = 500 * time.Millisecond

/*--DbConnectPage--*/
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
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false).
		AddItem(form, 0, 2, true).
		AddItem(tview.NewBox().SetBorder(false), 0, 1, false)

	pages.AddPage("DbConnectPage", flex, true, false)

}

func saveHandler(app *tview.Application, pages *tview.Pages) {

	client.Init(username, password, dbname)
	users = client.GetUsers()
	MqttClientsPage(app, pages)
	pages.SwitchToPage("MqttClientsPage")

}

/*--MqttClientsPage--*/
func MqttClientsPage(app *tview.Application, pages *tview.Pages) {

	ConnectAsForm := tview.NewForm().
		AddInputField("Client ID", clientID, 20, nil, func(text string) { clientID = text }).
		AddInputField("User name", clientusername, 20, nil, func(text string) { clientusername = text }).
		AddInputField("Password", clientpassword, 20, nil, func(text string) { clientpassword = text }).
		AddButton("Connect", func() { connectHandler(app, pages) })
	ConnectAsForm.SetBorder(true).SetTitle("Connect as:").SetTitleAlign(tview.AlignLeft)

	textarea := tview.NewTextArea().SetPlaceholder("sername,password,clientid,topic,pubsub,accses")
	textarea.SetBorder(true)
	textarea.SetBackgroundColor(tcell.ColorLightSeaGreen)

	AddNewUsers := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textarea, 0, 30, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewButton("Add"), 0, 1, false).
			AddItem(tview.NewButton("Save csv"), 0, 1, false), 0, 2, false)

	leftmenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ConnectAsForm, 0, 1, false).
		AddItem(AddNewUsers, 0, 1, false)
	leftmenu.SetTitle("Menu:")

	rightmenu := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)

	rightmenu.SetCell(0, 0, &tview.TableCell{Text: "Username", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 1, &tview.TableCell{Text: "Password", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 2, &tview.TableCell{Text: "ClientId", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 3, &tview.TableCell{Text: "Topic", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 4, &tview.TableCell{Text: "Action", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 5, &tview.TableCell{Text: "Access", Align: tview.AlignCenter, Color: tcell.ColorYellow})

	for row := 0; row < len(users); row++ {
		color := tcell.ColorWhite
		align := tview.AlignLeft
		exectrow := row + 1
		rightmenu.SetCell(exectrow, 0, &tview.TableCell{Text: users[row].Username, Align: align, Color: color}).
			SetCell(exectrow, 1, &tview.TableCell{Text: users[row].Password, Align: align, Color: color}).
			SetCell(exectrow, 2, &tview.TableCell{Text: users[row].ClientId, Align: align, Color: color}).
			SetCell(exectrow, 3, &tview.TableCell{Text: users[row].Topic, Align: align, Color: color}).
			SetCell(exectrow, 4, &tview.TableCell{Text: users[row].Action, Align: align, Color: color}).
			SetCell(exectrow, 5, &tview.TableCell{Text: users[row].Access, Align: align, Color: color})
	}
	rightmenu.SetTitle("Users")
	rightmenu.SetBorders(true)

	grid := tview.NewGrid().
		SetRows(100).
		SetColumns(40, 50).
		SetBorders(true).
		AddItem(leftmenu, 0, 0, 1, 1, 0, 0, false).
		AddItem(rightmenu, 0, 1, 51, 10, 0, 0, false)
		//setborder -> box
	pages.AddPage("MqttClientsPage", grid, true, false)

}
func connectHandler(app *tview.Application, pages *tview.Pages) {
	//var str string

	client.ConnectMqtt(clientID, clientusername, clientpassword)
	MonitoringPage(app, pages)
	pages.SwitchToPage("MonitoringPage")
	//go flaghandler(app)
}

/*--MonitoringPage--*/
func MonitoringPage(app *tview.Application, pages *tview.Pages) {

	info := tview.NewTextArea().SetPlaceholder("Here will be a system data")
	info.SetBorder(true)

	btns := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewButton("Start monitoring").SetSelectedFunc(func() { client.Flag = true; StartInfoTextView(app, info) }), 0, 1, false).
		AddItem(tview.NewButton("Stop monitoring").SetSelectedFunc(func() { client.Flag = false }), 0, 1, false).
		AddItem(tview.NewButton("Save").SetSelectedFunc(saveInfoHandler), 0, 1, false).
		AddItem(tview.NewButton("Disconnect").SetSelectedFunc(func() { disconnectHandler(pages) }), 0, 1, false)
	btns.SetBorder(true)
	btns.SetTitle("Menu")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(info, 0, 10, false).
		AddItem(btns, 0, 1, false)
	flex.SetBorder(true)
	pages.AddPage("MonitoringPage", flex, true, false)
}

func StartInfoTextView(app *tview.Application, info *tview.TextArea) {

	go func() {
		for {
			if client.Flag == true {
				counter += 1
				for {
					if client.Flag && counter < 2 {
						time.Sleep(refreshInterval)
						str = "\n" + client.GetCurrentUser() + "\n"
						str += systeminfo.GetMemUsage()
						str += systeminfo.GetPercent()
						str += systeminfo.GetPercentEvery()
						app.QueueUpdateDraw(func() {
							info.SetText(str, false)
						})

						client.Publish("topic/info", str)
					} else {
						counter -= 1
						break
					}
				}
			} else if disconnectflag {
				break
			}
		}
	}()
}

func flaghandler(app *tview.Application) {
	// for {
	// 	if client.Flag == true {
	// 		StartInfoTextView(app)
	// 	}
	// }

}

func saveInfoHandler() {

}

func disconnectHandler(pages *tview.Pages) {
	client.Flag = false
	disconnectflag = true
	client.DisconnectMQTT()
	pages.SwitchToPage("MqttClientsPage")
}
