package services

import (
	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/models"
)

type ArchiveService struct {
}

func NewArchiveService() *ArchiveService {
	return &ArchiveService{}
}

func (s *ArchiveService) GetArchivePageList(where map[string]interface{}, page int, pageSize int) (list []models.Archive, count int) {
	err := models.PageList(&models.Archive{}, where, page, pageSize, &list, &count)
	if err != nil {
		logger.Error("GetArchivePageList Error:", err)
		return nil, 0
	}
	return
}

func (s *ArchiveService) CreateArchive(monitor *models.Archive) bool {
	err := models.Create(monitor)
	if err != nil {
		logger.Error("CreateArchive Error:", err)
		return false
	}
	return true
}
