package dto

import (
	"go/hioto-logger/pkg/enum"
	"time"
)

type ResponseDeviceDto struct {
	ID           uint             `json:"id"`
	Guid         string           `json:"guid"`
	Mac          string           `json:"mac"`
	Type         enum.EDeviceType `json:"type"`
	Quantity     int              `json:"quantity"`
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	Minor        string           `json:"minor"`
	Status       string           `json:"status"`
	StatusDevice string           `json:"status_device"`
	LastSeen     time.Time        `json:"last_seen"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type ResCloudDeviceDto struct {
	ResponseDeviceDto
	MacServer string `json:"mac_server"`
}
