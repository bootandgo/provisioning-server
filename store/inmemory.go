package store

import (
	"fmt"
	"sync"

	"provisioning-server/models"
)

type InMemoryStore struct {
	admins  map[string]*models.Admin
	servers map[string]*models.Server
	mu      sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		admins:  make(map[string]*models.Admin),
		servers: make(map[string]*models.Server),
	}
}

func (s *InMemoryStore) CreateAdmin(admin *models.Admin) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.admins[admin.Username]; exists {
		return fmt.Errorf("username already exists")
	}
	s.admins[admin.Username] = admin
	return nil
}

func (s *InMemoryStore) FindAdminByUsername(username string) (*models.Admin, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	admin, exists := s.admins[username]
	if !exists {
		return nil, fmt.Errorf("admin not found")
	}
	return admin, nil
}

func (s *InMemoryStore) CreateServer(server *models.Server) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[server.ID] = server
	return nil
}

func (s *InMemoryStore) ListServers() ([]*models.Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	servers := make([]*models.Server, 0, len(s.servers))
	for _, server := range s.servers {
		servers = append(servers, server)
	}
	return servers, nil
}

func (s *InMemoryStore) ApproveServer(serverID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	server, exists := s.servers[serverID]
	if !exists {
		return fmt.Errorf("server not found")
	}
	if server.Status != "pending" {
		return fmt.Errorf("server is not pending approval")
	}
	server.Status = "approved"
	return nil
}

func (s *InMemoryStore) FindServerBySerialNumber(serialNumber string) (*models.Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, server := range s.servers {
		if server.SerialNumber == serialNumber {
			return server, nil
		}
	}
	return nil, ErrServerNotFound
}
