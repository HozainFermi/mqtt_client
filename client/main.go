package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "1357902479"
// 	dbname   = "Devices"
// 	sslMode  = "disable"
// )

type Device struct {
	ClientId   string
	TempCPU    float32
	TempGPU    float32
	UpdateTime time.Time
}
type User struct {
	Username string
	Password string
	ClientId string
	Topic    string
	Action   string
	Access   string
}

var Flag bool = false

// "postgres://postgres:1357902479@localhost:5432/Devices?sslmode=disable"
func ConnectDB(user string, password string, dbname string) *pg.DB {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, "localhost", 5432, dbname, "disable")
	opt, errors := pg.ParseURL(url)
	if errors != nil {
		log.Fatal("Error connecting to database:", errors)
	}

	db := pg.Connect(opt)
	if db == nil {
		log.Fatal("Faild to connect to the database")
	} else {
		//log.Print("Succsesfuly connected to the database")
	}
	return db
}

var db *pg.DB
var client mqtt.Client

func CreareSchema() error {
	err := db.Model((*Device)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return err
	}
	erro := db.Model((*User)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if erro != nil {
		return erro
	}
	return err
}

func Init(user string, password string, dbname string) {
	db = ConnectDB(user, password, dbname)
	err := CreareSchema()
	if err != nil {
		log.Fatal("Faild to create table", err)
	}
}

func GetUsers() []User {
	var users []User
	err := db.Model(&users).Select()
	if err != nil {
		log.Fatal(err)
	}
	return users

}

func AddNewUsers(users []User) {

	_, err := db.Model(users).Insert()
	if err != nil {
		panic(err)
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if msg.Topic() == "topic/commands" && string(msg.Payload()) == "StartMonitoring" {
		Flag = true
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	} else if msg.Topic() == "topic/commands" && string(msg.Payload()) == "StopMonitoring" {
		Flag = false
		//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	go func() {
		time.Sleep(2000)
		fmt.Println("Connected")
	}()
	//fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func ConnectMqtt(clientID string, username string, password string) {
	var broker = "r44a800d.ala.eu-central-1.emqxsl.com"
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	tlsConfig := NewTlsConfig()
	opts.SetTLSConfig(tlsConfig)
	// other options
	opts.SetClientID(clientID) //"go_mqtt_client"
	opts.SetUsername(username) //"PC"
	opts.SetPassword(password) //"public"
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	Sub(client, "topic/commands")
	//Publish(client)

	//client.Disconnect(250)
}

func Publish(topic string, message string) {

	text := fmt.Sprintf("Message %s", message)
	token := client.Publish(topic, 0, false, text)
	token.Wait()
	//time.Sleep(time.Second)

}

func Sub(client mqtt.Client, topicname string) {
	topic := topicname
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	//fmt.Printf("Subscribed to topic: %s", topic)
}
func DisconnectMQTT() {
	client.Disconnect(250)
}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	path, erro := filepath.Abs("./emqxsl-ca.crt")
	if erro != nil {
		panic(erro)
	}
	ca, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	return &tls.Config{
		RootCAs: certpool,
	}
}
