package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttInstance struct {
	client mqtt.Client
}

type MqttConfig struct {
	InstanceName string
	Host         string
	Username     string
	Password     string
	ClientId     string
}

var mqttInstance = make(map[string]*MqttInstance)
var mqttMu sync.Mutex

func initializeMqtt(mqttConfig *MqttConfig) error {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	opts := mqtt.NewClientOptions().
		AddBroker(mqttConfig.Host).
		SetUsername(mqttConfig.Username).
		SetPassword(mqttConfig.Password).
		SetClientID(mqttConfig.ClientId).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(1 * time.Second)

	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("🔓 MQTT %s connected successfully", mqttConfig.InstanceName)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("⚠️ MQTT %s connection lost: %v", mqttConfig.InstanceName, err)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] Failed to connect to MQTT broker: %v", mqttConfig.InstanceName, token.Error())
	}

	mqttInstance[mqttConfig.InstanceName] = &MqttInstance{
		client: client,
	}

	return nil
}

func GetMqttInstance(instanceName string) (mqtt.Client, error) {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	instance, ok := mqttInstance[instanceName]

	if !ok {
		return nil, fmt.Errorf("MQTT instance %s not found", instanceName)
	}

	return instance.client, nil
}

func CloseAllMqttInstances() {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	for name, instance := range mqttInstance {
		if instance.client.IsConnectionOpen() {
			instance.client.Disconnect(250)
			log.Printf("🔒 MQTT %s connection closed", name)
		}
	}
}

func CreateMqttInstance() {
	if err := initializeMqtt(&MqttConfig{
		InstanceName: MQTT_LOCAL_INSTANCE_NAME.GetValue(),
		Host:         MQTT_LOCAL_HOST.GetValue(),
		Username:     MQTT_LOCAL_USERNAME.GetValue(),
		Password:     MQTT_LOCAL_PASSWORD.GetValue(),
		ClientId:     "listener-mqtt-local",
	}); err != nil {
		log.Println(err)
	}
}
