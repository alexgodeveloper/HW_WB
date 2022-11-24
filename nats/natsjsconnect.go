package nats

import (
	usergrade "HW_WB"
	"HW_WB/storage"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

var Nc *nats.Conn
var Js nats.JetStreamContext

// Соединение с jetstream
func NatsJsConnect() {
	var err error
	Nc, err = nats.Connect(nats.DefaultURL)
	fmt.Println(err)
	var e error
	Js, e = Nc.JetStream(nats.PublishAsyncMaxPending(256))
	fmt.Println(e)
	log.Println(Js.AddStream(&nats.StreamConfig{
		Name:     "UserGrades",
		Subjects: []string{"UserGrades.*"},
		NoAck:    false,
	}))
}

type Msg struct {
	AppName string              `json:"appName"`
	Data    usergrade.UserGrade `json:"data"`
}

// Публикация  в stream
func Publish(topic string, appName string, ug usergrade.UserGrade) {
	var m Msg
	m.AppName = appName
	m.Data = ug
	ugJson, _ := json.Marshal(m)
	_, err := Js.Publish(topic, ugJson)
	if err != nil {
		fmt.Println(err)
	}
}

// Подписка на stream
func Subscribe(topic string, appName string) {
	_, err := Js.Subscribe(topic, func(m *nats.Msg) {
		m.Ack()
		switch topic {
		case "UserGrades.*":
			var msg Msg
			json.Unmarshal(m.Data, &msg)
			if msg.AppName == appName {
				fmt.Println("Чтение своих сообщений")
			} else {

				userGrade := msg.Data
				storage.SetStore(userGrade)
				fmt.Println("Записано в стор")
			}

		default:
			fmt.Println("Нет такого топика")
		}
		log.Printf("monitor service subscribes from subject:%s\n", m.Subject)
	}, nats.Durable("MONITOR"), nats.ManualAck())

	if err != nil {
	} else {
		fmt.Println("Received a JetStream message")
	}
}
