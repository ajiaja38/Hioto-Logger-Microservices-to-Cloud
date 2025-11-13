package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

const (
	DB_PATH     EnvKey = "DB_PATH"
	MAC_ADDRESS EnvKey = "MAC_ADDRESS"

	// RMQ cloud
	RMQ_CLOUD_URI                 EnvKey = "RMQ_HIOTO_CLOUD_URL"
	RMQ_CLOUD_INSTANCE            EnvKey = "RMQ_HIOTO_CLOUD_INSTANCE"
	RMQ_QUEUE_MONITORING_RESPONSE EnvKey = "MONITORING_RESPONSE_QUEUE"
	RMQ_LOGS_QUEUE                EnvKey = "LOGS_QUEUE"
	RMQ_LOGS_AKTUATOR_QUEUE       EnvKey = "LOGS_AKTUATOR_QUEUE"
	RMQ_QUEUE_UPDATE_RESPONSE     EnvKey = "UPDATE_RES_CLOUD"

	// MQTT LOCAL
	MQTT_LOCAL_HOST          EnvKey = "MQTT_LOCAL_HOST"
	MQTT_LOCAL_USERNAME      EnvKey = "MQTT_LOCAL_USERNAME"
	MQTT_LOCAL_PASSWORD      EnvKey = "MQTT_LOCAL_PASSWORD"
	MQTT_LOCAL_INSTANCE_NAME EnvKey = "MQTT_LOCAL_INSTANCE_NAME"
	MQTT_TOPIC_STATUS_DEVICE EnvKey = "STATUS_DEVICE_TOPIC"
)

func LoadEnv() error {
	return godotenv.Load()
}

func (e EnvKey) GetValue() string {
	return os.Getenv(string(e))
}
