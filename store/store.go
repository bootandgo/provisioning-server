package store

import (
	"errors"
	"provisioning-server/models"
)

var (
	ErrDuplicateSerialNumber = errors.New("serial number already exists")
	ErrServerNotFound        = errors.New("server not found")
)

type Store interface {
	CreateAdmin(admin *models.Admin) error
	FindAdminByUsername(username string) (*models.Admin, error)
	CreateServer(server *models.Server) error
	ListServers() ([]*models.Server, error)
	ApproveServer(serverID string) error
	FindServerBySerialNumber(serialNumber string) (*models.Server, error)
}
