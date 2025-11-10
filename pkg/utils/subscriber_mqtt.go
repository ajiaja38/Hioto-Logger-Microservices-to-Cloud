package utils

import (
	"context"
	"go/hioto-logger/config"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageHandler func([]byte)

func ConsumeMQTTTopic(ctx context.Context, instanceName, topic string, handlerFunc MessageHandler) {
	client, err := config.GetMqttInstance(instanceName)

	if err != nil {
		log.Println(err)
		return
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		go handlerFunc(msg.Payload())
	}

	if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe: %v", token.Error())
		client.Disconnect(250)
		return
	}

	log.Printf("Subscribed to topic: %s", topic)

	<-ctx.Done()

	log.Printf("MQTT context done, cleaning up...")
	client.Unsubscribe(topic)
	client.Disconnect(250)
}
