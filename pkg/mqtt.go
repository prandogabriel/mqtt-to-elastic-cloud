package pkg

import (
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTTClient(broker string) MQTT.Client {
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID("mqtt_client")

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

func Subscribe(c MQTT.Client, topic string, callback MQTT.MessageHandler) {
	if token := c.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
