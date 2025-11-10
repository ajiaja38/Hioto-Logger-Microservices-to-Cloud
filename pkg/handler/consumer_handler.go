package handler

import (
	"go/hioto-logger/pkg/service"
	"strings"
)

type ConsumerHandler struct {
	deviceService *service.DeviceService
}

func NewConsumerHandler(deviceService *service.DeviceService) *ConsumerHandler {
	return &ConsumerHandler{
		deviceService: deviceService,
	}
}

func (h *ConsumerHandler) ChangeStatusDevice(message []byte) {
	messageString := string(message)

	guid := strings.Split(messageString, "#")[0]
	status := strings.Split(messageString, "#")[1]

	h.deviceService.UpdateStatusDevice(guid, status)
}
