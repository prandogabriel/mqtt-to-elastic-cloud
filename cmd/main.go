// main.go
package main

import (
	clients "bridge/pkg"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
	// Lendo variáveis de ambiente
	broker := os.Getenv("MQTT_BROKER")
	topic := os.Getenv("MQTT_TOPIC")
	cloudID := os.Getenv("ELASTIC_CLOUD_ID")
	apiKey := os.Getenv("ELASTIC_API_KEY")
	indexName := os.Getenv("ELASTIC_INDEX_NAME")

	// Verificando se as variáveis de ambiente estão definidas
	if broker == "" || topic == "" || cloudID == "" || apiKey == "" || indexName == "" {
		fmt.Println("Error: environment variables MQTT_BROKER, MQTT_TOPIC, ELASTIC_CLOUD_ID, ELASTIC_INDEX_NAME and ELASTIC_API_KEY must be set")
		os.Exit(1)
	}

	// Parte 1: Ouvindo um tópico MQTT em Go
	c := clients.NewMQTTClient(broker)

	// Parte 2: Conectando com Elastic Cloud
	es := clients.NewElasticClient(cloudID, apiKey)

	clients.Subscribe(c, topic, func(client MQTT.Client, msg MQTT.Message) {
		uuid := uuid.New().String()             // generate a new UUID
		date := time.Now().Format(time.RFC3339) // get the current date in RFC3339 format

		// Parse the received message into a map structure
		var message map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			panic(err)
		}

		// Add the UUID and date to the message
		message["id"] = uuid
		message["date"] = date

		// Convert the message back to JSON
		jsonMessage, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(jsonMessage))

		// Index the message into Elasticsearch
		result, err := es.Index(indexName, strings.NewReader(string(jsonMessage)))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Pub to elastic status: %s %s\n", result.Status(), result.String())
	})

	for {
	}
}
