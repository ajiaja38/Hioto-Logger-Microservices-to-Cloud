package service

import (
	"encoding/json"
	"go/hioto-logger/config"
	"go/hioto-logger/pkg/dto"
	"go/hioto-logger/pkg/enum"
	"go/hioto-logger/pkg/model"
	"go/hioto-logger/pkg/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

var location *time.Location

func init() {
	location = time.FixedZone("WIB", 7*60*60)
}

type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService(db *gorm.DB) *DeviceService {
	return &DeviceService{
		db: db,
	}
}

func (s *DeviceService) UpdateStatusDevice(guid string, status string) {
	var device model.Registration

	if err := s.db.First(&device, "guid = ?", guid).Error; err != nil {
		log.Printf("Device not found: %v 💥", err)
		return
	}

	device.StatusDevice = enum.EDeviceStatus(status)
	device.LastSeen = time.Now().In(location)

	if err := s.db.Save(&device).Error; err != nil {
		log.Printf("Error updating status device: %v 💥", err)
		return
	}

	bodyToCloud := dto.ResCloudDeviceDto{
		ResponseDeviceDto: dto.ResponseDeviceDto{
			ID:           device.ID,
			Guid:         device.Guid,
			Mac:          device.Mac,
			Type:         device.Type,
			Quantity:     device.Quantity,
			Name:         device.Name,
			Version:      device.Version,
			Minor:        device.Minor,
			Status:       device.Status,
			StatusDevice: string(device.StatusDevice),
			LastSeen:     device.LastSeen,
			CreatedAt:    device.CreatedAt,
			UpdatedAt:    device.UpdatedAt,
		},
		MacServer: config.MAC_ADDRESS.GetValue(),
	}

	jsonBody, err := json.Marshal(bodyToCloud)

	if err != nil {
		log.Printf("Error marshaling JSON: %v 💥", err)
		return
	}

	utils.PublishToRmq(config.RMQ_CLOUD_INSTANCE.GetValue(), jsonBody, config.RMQ_QUEUE_UPDATE_RESPONSE.GetValue(), "amq.direct")

	log.Printf("Status device successfully updated: %s ✅", status)
}

func (s *DeviceService) CheckInactiveDevice() {
	treshold := time.Now().Add(-10 * time.Second)

	if err := s.db.Model(&model.Registration{}).Where("last_seen < ?", treshold).Update("status_device", enum.OFF).Error; err != nil {
		log.Printf("Error checking for inactive device: %v 💥", err)
	}

	log.Printf("Inactive devices marked as offline 🔻")
}
