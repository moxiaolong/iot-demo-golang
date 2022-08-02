package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"time"
)

func Conn(clientId string) MQTT.Client {

	log.Println("MQTT ClientId", clientId)
	username := ""
	password := ""
	server := "tcp://" + "127.0.0.1" + ":" + strconv.Itoa(1883)
	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientId).SetCleanSession(true)
	connOpts.SetUsername(username)
	connOpts.SetPassword(password)
	connOpts.SetAutoReconnect(true)
	connOpts.SetMaxReconnectInterval(5)

	MqttClient := MQTT.NewClient(connOpts)

	for {
		if MqttClient.IsConnected() {
			time.Sleep(time.Second * 15)
			continue
		}
		log.Println("Connecting... to mqtt:", server)
		if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		} else {
			log.Println("Connected to mqtt:", server)
			return MqttClient
		}
	}

}
