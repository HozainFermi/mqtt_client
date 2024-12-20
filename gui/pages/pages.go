package pages

import (
	"fmt"
	"kurs/client"
	"kurs/systeminfo"
	"os"
	"strings"
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
var clientsgrid *tview.Grid

//var info *tview.TextArea

var str string = ""
var counter int = 0
var disconnectflag bool = false
var viewmonitoringflag bool = false

const refreshInterval = 1500 * time.Millisecond

/*--DbConnectPage--*/
func DbConnectPage(app *tview.Application, pages *tview.Pages) {

	form := tview.NewForm().
		AddInputField("User name", username, 20, nil, func(text string) { username = text }).
		AddPasswordField("Password", password, 20, '*', func(text string) { password = text }).
		AddInputField("Database name", dbname, 20, nil, func(text string) { dbname = text }).
		AddInputField("Mqtt broker", client.Broker, 20, nil, func(text string) { client.Broker = text }).
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

	var rightmenu *tview.Table
	rightmenu = tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)

	ConnectAsForm := tview.NewForm().
		AddInputField("Client ID", clientID, 20, nil, func(text string) { clientID = text }).
		AddInputField("User name", clientusername, 20, nil, func(text string) { clientusername = text }).
		AddInputField("Password", clientpassword, 20, nil, func(text string) { clientpassword = text }).
		AddButton("Connect", func() { connectHandler(app, pages) })
	ConnectAsForm.SetBorder(true).SetTitle("Connect as:").SetTitleAlign(tview.AlignLeft)

	textarea := tview.NewTextArea().SetPlaceholder("username,password,clientid,topic,pubsub,allow/deny(newline)")
	textarea.SetBorder(true)
	textarea.SetBackgroundColor(tcell.ColorLightSeaGreen)

	AddNewUsers := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textarea, 0, 50, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewButton("Add").SetSelectedFunc(func() {
				result := AddUsersHandler(textarea.GetText(), rightmenu)
				text := textarea.GetText() + result
				textarea.SetText(text, true)
			}), 0, 1, false).
			AddItem(tview.NewButton("Save csv"), 0, 1, false), 0, 2, false)

	leftmenu := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ConnectAsForm, 0, 1, false).
		AddItem(AddNewUsers, 0, 2, false)
	leftmenu.SetTitle("Menu:")

	//rightmenu := tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)

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

	rightmenu.SetSelectedFunc(func(row, column int) {
		// Get the values of the selected row.
		var selectedRow []string
		for col := 0; col < 6; col++ {
			cell := rightmenu.GetCell(row, col)
			if cell != nil {
				selectedRow = append(selectedRow, cell.Text)

			}
		}
		clientID = selectedRow[2]
		clientusername = selectedRow[0]
		clientpassword = selectedRow[1]

		item := ConnectAsForm.GetFormItem(0)
		inputid := item.(*tview.InputField)
		inputid.SetText(clientID)

		item = ConnectAsForm.GetFormItem(1)
		inputid = item.(*tview.InputField)
		inputid.SetText(clientusername)

		item = ConnectAsForm.GetFormItem(2)
		inputid = item.(*tview.InputField)
		inputid.SetText(clientpassword)

	})

	clientsgrid = tview.NewGrid().
		SetRows(100).
		SetColumns(-5, -1).
		SetBorders(true).
		AddItem(leftmenu, 0, 0, 1, 1, 0, 0, false).
		AddItem(rightmenu, 0, 1, 51, 10, 0, 0, false)
		//setborder -> box
	pages.AddPage("MqttClientsPage", clientsgrid, true, false)

}

func AddUsersHandler(text string, rightmenu *tview.Table) string {

	var parsedusers []string = strings.Split(text, "\n")
	var userarray []client.User
	var user []string

	for index, element := range parsedusers {
		user = strings.Split(element, ",")
		if len(user) != 6 {
			return fmt.Sprintf("Not enought arguments in the %d string", index)
		}
		userarray = append(userarray, client.User{Username: user[0], Password: user[1], ClientId: user[2], Topic: user[3], Action: user[4], Access: user[5]})
	}
	for row := len(users); row < len(users)+len(userarray); row++ {
		color := tcell.ColorWhite
		align := tview.AlignLeft
		exectrow := row + 1
		rightmenu.SetCell(exectrow, 0, &tview.TableCell{Text: userarray[row].Username, Align: align, Color: color}).
			SetCell(exectrow, 1, &tview.TableCell{Text: userarray[row].Password, Align: align, Color: color}).
			SetCell(exectrow, 2, &tview.TableCell{Text: userarray[row].ClientId, Align: align, Color: color}).
			SetCell(exectrow, 3, &tview.TableCell{Text: userarray[row].Topic, Align: align, Color: color}).
			SetCell(exectrow, 4, &tview.TableCell{Text: userarray[row].Action, Align: align, Color: color}).
			SetCell(exectrow, 5, &tview.TableCell{Text: userarray[row].Access, Align: align, Color: color})
	}

	err := client.AddNewUsers(userarray)
	if err == nil {
		//!!как перерисовать таблицу то бл :(
		return "\nAdded successfully"

	} else {
		return err.Error()
	}

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
		AddItem(tview.NewButton("View monitoring").SetSelectedFunc(func() {
			viewmonitoringflag = true
			client.Flag = false
			client.SubToInfoTopic()
			StartViewMonitoring(app, info)
		}), 0, 1, false).
		AddItem(tview.NewButton("Stop monitoring").SetSelectedFunc(func() { client.Flag = false; viewmonitoringflag = false }), 0, 1, false).
		AddItem(tview.NewButton("Save").SetSelectedFunc(saveInfoHandler), 0, 1, false).
		AddItem(tview.NewButton("Clear").SetSelectedFunc(func() { info.SetText("", false) }), 0, 1, false).
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
						str = client.GetCurrentUser() + "\n" + systeminfo.GetMemUsage() + systeminfo.GetPercent() + systeminfo.GetPercentEvery()
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

func StartViewMonitoring(app *tview.Application, info *tview.TextArea) {
	go func() {
		for {
			if viewmonitoringflag {
				counter += 1
				for {
					if viewmonitoringflag && counter < 2 && !client.Flag {
						time.Sleep(refreshInterval)
						str = info.GetText()
						str += client.DisplayInfo
						// str += client.Payload[0] + "\n"
						// str += client.DisplayInfo
						app.QueueUpdateDraw(func() {
							info.SetText(str, true)
						})
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

	fo, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}

	if _, err := fo.Write([]byte(str)); err != nil {

	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
}

func disconnectHandler(pages *tview.Pages) {
	client.Flag = false
	disconnectflag = true
	client.DisconnectMQTT()
	pages.SwitchToPage("MqttClientsPage")
}
