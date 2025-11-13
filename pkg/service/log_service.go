package service

import (
	"encoding/json"
	"go/hioto-logger/config"
	"go/hioto-logger/pkg/dto"
	"go/hioto-logger/pkg/model"
	"go/hioto-logger/pkg/utils"
	"log"

	"gorm.io/gorm"
)

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{
		db: db,
	}
}

func (s *LogService) GetAllLogs() {
	log.Print("Getting all logs...")

	var logs []model.Log

	if err := s.db.Raw(`
        SELECT
            id, 
            input_guid, 
            input_name, 
            input_value, 
            output_guid, 
            output_value, 
            time 
        FROM logs
	`).Scan(&logs).Error; err != nil {
		log.Printf("Failed to get logs: %v", err)
		return
	}

	if len(logs) == 0 {
		log.Print("No logs is found")
		return
	}

	var payload []dto.LogsReponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogsReponseDto{
			ID:          log.ID,
			InputName:   log.InputName,
			InputGuid:   log.InputGuid,
			InputValue:  log.InputValue,
			OutputGuid:  log.OutputGuid,
			OutputValue: log.OutputValue,
			Time:        log.Time,
			MacServer:   config.MAC_ADDRESS.GetValue(),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal logs: %v", err)
		return
	}

	utils.PublishToRmq(config.RMQ_CLOUD_INSTANCE.GetValue(), body, config.RMQ_LOGS_AKTUATOR_QUEUE.GetValue(), "amq.direct")

	if err := s.db.Exec("DELETE FROM logs").Error; err != nil {
		log.Printf("Failed to delete logs: %v", err)
	}
}

func (s *LogService) GetAllLogAktuators() {
	log.Print("Getting all log aktuators...")

	var logs []model.LogAktuator

	if err := s.db.Raw(`SELECT * FROM log_aktuators`).Scan(&logs).Error; err != nil {
		log.Printf("Failed to get logs: %v", err)
		return
	}

	if len(logs) == 0 {
		log.Print("No logs found")
		return
	}

	var payload []dto.LogAktuatorReponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogAktuatorReponseDto{
			ID:        log.ID,
			Guid:      log.InputGuid,
			Name:      log.Name,
			Value:     log.Value,
			Time:      log.Time,
			MacServer: config.MAC_ADDRESS.GetValue(),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to json marshal logs: %v", err)
		return
	}

	utils.PublishToRmq(
		config.RMQ_CLOUD_INSTANCE.GetValue(),
		body,
		config.RMQ_LOGS_AKTUATOR_QUEUE.GetValue(),
		"amq.direct",
	)

	if err := s.db.Exec("DELETE FROM log_aktuators").Error; err != nil {
		log.Printf("Failed to delete logs: %v", err)
		return
	}
}

func (s *LogService) GetAllMonitoringHistory() {
	log.Print("Getting all monitoring history...")

	var logs []model.MonitoringHistory

	if err := s.db.Raw(`SELECT * FROM monitoring_histories`).Scan(&logs).Error; err != nil {
		log.Printf("Failed to get monitoring history: %v", err)
		return
	}

	if len(logs) == 0 {
		log.Print("No monitoring history found")
		return
	}

	var payload []dto.LogMonitoringDeviceHistoryResponseDto

	for _, log := range logs {
		payload = append(payload, dto.LogMonitoringDeviceHistoryResponseDto{
			ID:         log.ID,
			DeviceGuid: log.DeviceGuid,
			DeviceName: log.DeviceName,
			DeviceType: log.DeviceType,
			Value:      log.Value,
			Time:       log.Time,
			MacServer:  config.MAC_ADDRESS.GetValue(),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to json marshal monitoring history: %v", err)
		return
	}

	utils.PublishToRmq(config.RMQ_CLOUD_INSTANCE.GetValue(), body, config.RMQ_QUEUE_MONITORING_RESPONSE.GetValue(), "amq.direct")

	if err := s.db.Exec("DELETE FROM monitoring_histories").Error; err != nil {
		log.Printf("Failed to delete monitoring history: %v", err)
		return
	}
}
