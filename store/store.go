package store

import (
	"errors"
	"provisioning-server/models"
)

var (
	ErrDuplicateSerialNumber = errors.New("serial number already exists")
	ErrServerNotFound        = errors.New("server not found")
	ErrInvalidInvitation     = errors.New("invalid or expired invitation token")
	ErrInvitationExists      = errors.New("invitation already exists")
)

type Store interface {
	CreateAdmin(admin *models.Admin) error
	FindAdminByUsername(username string) (*models.Admin, error)
	CreateServer(server *models.Server) error
	ListServers() ([]*models.Server, error)
	ApproveServer(serverID string) error
	FindServerBySerialNumber(serialNumber string) (*models.Server, error)
	GetServerByID(serverID string) (*models.Server, error)
	CreateInvitation(invite *models.Invitation) error
	GetInvitation(token string) (*models.Invitation, error)
	MarkInvitationUsed(token string) error
	GetRootAdmin() (*models.Admin, error)
}
