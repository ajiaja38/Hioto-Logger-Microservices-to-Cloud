package service

import (
	"encoding/json"
	"go/hioto-logger/pkg/dto"
	"go/hioto-logger/pkg/model"
	"go/hioto-logger/pkg/utils"
	"log"
	"os"

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
			MacServer:   os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal logs: %v", err)
		return
	}

	utils.PublishToRmq(os.Getenv("RMQ_HIOTO_CLOUD_INSTANCE"), body, os.Getenv("LOGS_QUEUE"), "amq.direct")

	queryDeleteAll := s.db.Exec("DELETE FROM logs")

	if queryDeleteAll.Error != nil {
		log.Printf("Failed to delete logs: %v", queryDeleteAll.Error)
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
			MacServer: os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to json marshal logs: %v", err)
		return
	}

	utils.PublishToRmq(os.Getenv("RMQ_HIOTO_CLOUD_INSTANCE"), body, os.Getenv("LOGS_AKTUATOR_QUEUE"), "amq.direct")

	queryDeleteAll := s.db.Exec("DELETE FROM log_aktuators")

	if queryDeleteAll.Error != nil {
		log.Printf("Failed to delete logs: %v", queryDeleteAll.Error)
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
			MacServer:  os.Getenv("MAC_ADDRESS"),
		})
	}

	body, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to json marshal monitoring history: %v", err)
		return
	}

	utils.PublishToRmq(os.Getenv("RMQ_HIOTO_CLOUD_INSTANCE"), body, os.Getenv("MONITORING_RESPONSE_QUEUE"), "amq.direct")

	if err := s.db.Exec("DELETE FROM monitoring_histories").Error; err != nil {
		log.Printf("Failed to delete monitoring history: %v", err)
		return
	}
}
