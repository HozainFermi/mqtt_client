package pages

import (
	"fmt"
	"kurs/client"
	"kurs/systeminfo"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var username string
var password string
var dbname string

var clientID string = ""
var clientusername string = ""
var clientpassword string = ""

var users []client.User
var clientsgrid *tview.Grid

var counter int = 0

// FLAGS
var disablelogflag bool = false
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
	var index = 0

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
		rightmenu.SetCell(exectrow, 0, &tview.TableCell{Text: userarray[index].Username, Align: align, Color: color}).
			SetCell(exectrow, 1, &tview.TableCell{Text: userarray[index].Password, Align: align, Color: color}).
			SetCell(exectrow, 2, &tview.TableCell{Text: userarray[index].ClientId, Align: align, Color: color}).
			SetCell(exectrow, 3, &tview.TableCell{Text: userarray[index].Topic, Align: align, Color: color}).
			SetCell(exectrow, 4, &tview.TableCell{Text: userarray[index].Action, Align: align, Color: color}).
			SetCell(exectrow, 5, &tview.TableCell{Text: userarray[index].Access, Align: align, Color: color})

		index++
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

	if clientID == "" {
		return
	}
	client.ConnectMqtt(clientID, clientusername, clientpassword)
	MonitoringPage(app, pages)
	pages.SwitchToPage("MonitoringPage")
	//go flaghandler(app)
}

/*--MonitoringPage--*/
func MonitoringPage(app *tview.Application, pages *tview.Pages) {
	// Основной flex контейнер
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	// Заголовок с информацией о клиенте
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	header.SetBorder(false)
	updateHeader(header)

	// Таблица для отображения данных с группировкой по client ID
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	table.SetBorder(true).SetTitle(" Данные клиентов ").SetBorderPadding(0, 0, 1, 1)

	// Устанавливаем заголовки таблицы
	setTableHeaders(table)

	// Область для лога (если нужно показывать лог)
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	logView.SetBorder(true).SetTitle(" Лог событий ")
	logView.SetChangedFunc(func() {
		app.Draw()
		logView.ScrollToEnd()
	})

	// Статус бар
	statusBar := tview.NewTextView().
		SetDynamicColors(true)
	statusBar.SetBorder(false)
	updateStatusBar(statusBar, "Готов к работе")

	// Кнопки управления
	buttons := createButtons(app, pages, table, logView, statusBar)

	// Инициализируем map для хранения данных клиентов
	clientData = make(map[string]*ClientData)

	// Сборка интерфейса
	mainFlex.
		AddItem(header, 1, 0, false).
		AddItem(table, 0, 8, false).   // 80% для таблицы
		AddItem(logView, 0, 2, false). // 20% для лога
		AddItem(statusBar, 1, 0, false).
		AddItem(buttons, 3, 0, true)

	pages.AddPage("MonitoringPage", mainFlex, true, false)
}

// Структура для хранения данных клиента
type ClientData struct {
	Metrics  map[string]string
	LastSeen time.Time
	RowIndex int
}

var clientData map[string]*ClientData
var currentRow = 1 // начинаем с 1, т.к. 0 - заголовки

func setTableHeaders(table *tview.Table) {
	headers := []string{"Client ID", "CPU %", "Memory", "Disk", "Network", "Last Update"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}
}

func updateClientInTable(table *tview.Table, clientID string, metrics map[string]string) {
	// Ищем клиента в данных
	data, exists := clientData[clientID]
	if !exists {
		// Новый клиент - добавляем строку
		data = &ClientData{
			Metrics:  make(map[string]string),
			LastSeen: time.Now(),
			RowIndex: currentRow,
		}
		clientData[clientID] = data
		currentRow++
	}

	// Обновляем метрики
	for k, v := range metrics {
		data.Metrics[k] = v
	}
	data.LastSeen = time.Now()

	// Обновляем ячейки таблицы
	table.SetCell(data.RowIndex, 0,
		tview.NewTableCell(clientID).
			SetAlign(tview.AlignLeft))

	table.SetCell(data.RowIndex, 1,
		tview.NewTableCell(data.Metrics["cpu"]).
			SetAlign(tview.AlignCenter).
			SetTextColor(getValueColor(data.Metrics["cpu"])))

	table.SetCell(data.RowIndex, 2,
		tview.NewTableCell(data.Metrics["memory"]).
			SetAlign(tview.AlignCenter).
			SetTextColor(getValueColor(data.Metrics["memory"])))

	table.SetCell(data.RowIndex, 3,
		tview.NewTableCell(data.Metrics["disk"]).
			SetAlign(tview.AlignCenter).
			SetTextColor(getValueColor(data.Metrics["disk"])))

	table.SetCell(data.RowIndex, 4,
		tview.NewTableCell(data.Metrics["network"]).
			SetAlign(tview.AlignCenter).
			SetTextColor(getValueColor(data.Metrics["network"])))

	table.SetCell(data.RowIndex, 5,
		tview.NewTableCell(data.LastSeen.Format("15:04:05")).
			SetAlign(tview.AlignCenter))
}

func getValueColor(value string) tcell.Color {
	if strings.Contains(value, "%") {
		// Анализируем процентное значение
		if perc := extractPercentage(value); perc > 80 {
			return tcell.ColorRed
		} else if perc > 60 {
			return tcell.ColorYellow
		}
	}
	return tcell.ColorGreen
}

func extractPercentage(text string) int {
	// Простая функция для извлечения числа из строки типа "45%"
	if idx := strings.Index(text, "%"); idx != -1 {
		if num, err := strconv.Atoi(strings.TrimSpace(text[:idx])); err == nil {
			return num
		}
	}
	return 0
}

func createButtons(app *tview.Application, pages *tview.Pages, table *tview.Table, logView *tview.TextView, statusBar *tview.TextView) *tview.Flex {
	buttons := tview.NewFlex().SetDirection(tview.FlexColumn)

	startBtn := tview.NewButton("Старт")
	startBtn.SetSelectedFunc(func() {
		client.Flag = true
		viewmonitoringflag = false
		StartInfoTextView(app, table, logView, statusBar)
		updateStatusBar(statusBar, "Мониторинг запущен")
	})

	viewBtn := tview.NewButton("Просмотр")
	viewBtn.SetSelectedFunc(func() {
		viewmonitoringflag = true
		client.Flag = false
		client.SubToInfoTopic()
		StartViewMonitoring(app, table, logView, statusBar)
		updateStatusBar(statusBar, "Режим просмотра")
	})

	stopBtn := tview.NewButton("Стоп")
	stopBtn.SetSelectedFunc(func() {
		client.Flag = false
		viewmonitoringflag = false
		updateStatusBar(statusBar, "Остановлено")
	})

	saveBtn := tview.NewButton("Сохранить")
	saveBtn.SetSelectedFunc(func() {
		saveInfoHandler(getAllLogData(logView))
		updateStatusBar(statusBar, "Данные сохранены в log.txt")
	})

	clearBtn := tview.NewButton("Очистить")
	clearBtn.SetSelectedFunc(func() {
		logView.SetText("")
		updateStatusBar(statusBar, "Лог очищен")
	})

	ableBtn := tview.NewButton("Открыть Лог")
	ableBtn.SetSelectedFunc(func() {
		disablelogflag = false
		updateStatusBar(statusBar, "Лог открыт")
	})

	disableBtn := tview.NewButton("Закрыть Лог")
	disableBtn.SetSelectedFunc(func() {
		logView.Clear()
		disablelogflag = true
		updateStatusBar(statusBar, "Лог закрыт")
	})

	disconnectBtn := tview.NewButton("Отключиться")
	disconnectBtn.SetSelectedFunc(func() {
		disconnectHandler(pages)
	})

	// Добавляем кнопки
	buttons.AddItem(startBtn, 10, 1, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(viewBtn, 10, 1, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(stopBtn, 10, 1, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(saveBtn, 10, 1, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(clearBtn, 10, 1, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(disableBtn, 13, 3, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(ableBtn, 13, 3, true)
	buttons.AddItem(tview.NewBox(), 1, 0, false)
	buttons.AddItem(disconnectBtn, 12, 1, true)

	buttons.SetBorder(true).SetTitle(" Управление ")
	return buttons
}

func getAllLogData(logView *tview.TextView) string {
	return logView.GetText(false)
}

// Обновленные функции мониторинга
func StartInfoTextView(app *tview.Application, table *tview.Table, logView *tview.TextView, statusBar *tview.TextView) {
	go func() {
		for {
			if client.Flag {

				time.Sleep(refreshInterval)

				systemData := getSystemInfo()

				app.QueueUpdateDraw(func() {
					// Записываем в лог
					if !disablelogflag {

						fmt.Fprintf(logView, "[%s] %s\n",
							time.Now().Format("15:04:05"),
							systemData)
					}

					// Парсим данные и обновляем таблицу
					metrics := parseMetrics(systemData)
					clientID := client.GetCurrentUser()
					updateClientInTable(table, clientID, metrics)

					updateStatusBar(statusBar, "Отправка данных...")
				})

				client.Publish("topic/info", systemData)
			}
		}
	}()
}

func StartViewMonitoring(app *tview.Application, table *tview.Table, logView *tview.TextView, statusBar *tview.TextView) {
	go func() {
		for {
			if viewmonitoringflag {

				time.Sleep(refreshInterval)

				if client.DisplayInfo != "" {
					app.QueueUpdateDraw(func() {
						currentData := client.DisplayInfo
						client.DisplayInfo = ""

						// Записываем в лог
						if !disablelogflag {

							fmt.Fprintf(logView, "[%s] %s\n",
								time.Now().Format("15:04:05"),
								currentData)

						}

						// Парсим входящие данные и обновляем таблицу
						metrics := parseMetrics(currentData)
						clientID := extractClientID(currentData)
						if clientID != "" {
							updateClientInTable(table, clientID, metrics)
						}

						err := client.SaveMonitoringData(clientID, currentData)
						if err != nil {
							fmt.Fprintf(logView, "[red]Ошибка сохранения в БД: %s[white]\n", err.Error())
						} else {
							//fmt.Fprintf(logView, "[green]Данные сохранены в БД[white]\n")
						}

						updateStatusBar(statusBar, "Получены новые данные")
					})
				}
			}
		}
	}()
}

func parseMetrics(data string) map[string]string {
	metrics := make(map[string]string)
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Memory usage=") {
			metrics["memory"] = extractValueAfterEquals(line)
		} else if strings.Contains(line, "CPU usage (combined)=") {
			metrics["cpu"] = extractValueAfterEquals(line)
		} else if strings.Contains(line, "Disk") {
			metrics["disk"] = extractValueAfterEquals(line)
		} else if strings.Contains(line, "Network") {
			metrics["network"] = extractValueAfterEquals(line)
		}
	}

	// Если нет combined CPU, пытаемся найти другое CPU значение
	if metrics["cpu"] == "" {
		for _, line := range lines {
			if strings.Contains(line, "CPU") && strings.Contains(line, "=") && !strings.Contains(line, "Usage[") {
				metrics["cpu"] = extractValueAfterEquals(line)
				break
			}
		}
	}

	return metrics
}

func extractValueAfterEquals(line string) string {
	// Извлекаем значение после знака равно "Memory usage=77%"
	parts := strings.Split(line, "=")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func extractClientID(data string) string {
	// Парсим client ID из данных
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Client:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// Остальные функции остаются без изменений
func updateHeader(header *tview.TextView) {
	clientID := client.GetCurrentUser()
	if clientID == "" {
		fmt.Fprintf(header, "[red]Не подключено к MQTT[white]")
		return
	}
	fmt.Fprintf(header, "Клиент: %s | Брокер: %s", clientID, client.Broker)
}

func updateStatusBar(statusBar *tview.TextView, message string) {
	now := time.Now().Format("15:04:05")
	var status string
	if client.Flag {
		status = "[green]ВКЛ[white]"
	} else if viewmonitoringflag {
		status = "[blue]ПРОСМОТР[white]"
	} else {
		status = "[red]ВЫКЛ[white]"
	}
	statusBar.SetText(fmt.Sprintf("[%s] %s | %s", now, status, message))
}

func getSystemInfo() string {
	memUsage := systeminfo.GetMemUsage()
	percent := systeminfo.GetPercent()
	percentEvery := systeminfo.GetPercentEvery()

	return fmt.Sprintf("Client: %s\n%s%s%s",
		client.GetCurrentUser(),
		memUsage,
		percent,
		percentEvery)
}

func saveInfoHandler(data string) {
	fo, err := os.Create("log.txt")
	if err != nil {
		return
	}
	defer fo.Close()
	fo.Write([]byte(data))
}

//end

func flaghandler(app *tview.Application) {
	// for {
	// 	if client.Flag == true {
	// 		StartInfoTextView(app)
	// 	}
	// }

}

// func saveInfoHandler(str string) {

// 	fo, err := os.Create("log.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if _, err := fo.Write([]byte(str)); err != nil {

// 	}

// 	// close fo on exit and check for its returned error
// 	defer func() {
// 		if err := fo.Close(); err != nil {
// 			panic(err)
// 		}
// 	}()
// }

func disconnectHandler(pages *tview.Pages) {
	client.Flag = false
	disconnectflag = true
	client.DisconnectMQTT()
	pages.SwitchToPage("MqttClientsPage")
}
