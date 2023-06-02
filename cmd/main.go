// main.go
package main

import (
	clients "bridge/pkg"
	"fmt"
	"os"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Lendo variáveis de ambiente
	broker := os.Getenv("MQTT_BROKER")
	topic := os.Getenv("MQTT_TOPIC")
	cloudID := os.Getenv("ELASTIC_CLOUD_ID")
	apiKey := os.Getenv("ELASTIC_API_KEY")

	// Verificando se as variáveis de ambiente estão definidas
	if broker == "" || topic == "" || cloudID == "" || apiKey == "" {
		fmt.Println("Error: environment variables MQTT_BROKER, MQTT_TOPIC, ELASTIC_CLOUD_ID and ELASTIC_API_KEY must be set")
		os.Exit(1)
	}

	// Parte 1: Ouvindo um tópico MQTT em Go
	c := clients.NewMQTTClient(broker)

	// Parte 2: Conectando com Elastic Cloud
	es := clients.NewElasticClient(cloudID, apiKey)

	clients.Subscribe(c, topic, func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("MSG: %s\n", msg.Payload())

		// Enviando dados para Elastic Cloud
		result, err := es.Index("search-data",
			strings.NewReader(`{"message": "`+string(msg.Payload())+`"}`))

		fmt.Printf("Pub to elastic status: %s\n", result.Status())

		if err != nil {
			// Handle error
			panic(err)
		}
	})

	for {
	}
}
