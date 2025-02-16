package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"provisioning-server/models"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateAdmin(admin *models.Admin) error {
	query := `
		INSERT INTO admins (username, password, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (username) DO NOTHING`

	_, err := s.db.Exec(query, admin.Username, admin.Password, admin.CreatedAt)
	if err != nil {
		return fmt.Errorf("create admin error: %w", err)
	}
	return nil
}

func (s *PostgresStore) FindAdminByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	query := `SELECT username, password, created_at FROM admins WHERE username = $1`
	err := s.db.Get(&admin, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("admin not found")
		}
		return nil, fmt.Errorf("find admin error: %w", err)
	}
	return &admin, nil
}

func (s *PostgresStore) CreateServer(server *models.Server) error {
	query := `
		INSERT INTO servers (id, serial_number, ip_address, status, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(query,
		server.ID,
		server.SerialNumber,
		server.IPAddress,
		server.Status,
		server.CreatedAt,
	)
	return err
}

func (s *PostgresStore) ListServers() ([]*models.Server, error) {
	var servers []*models.Server
	query := `SELECT id, serial_number, ip_address, status, created_at FROM servers`
	err := s.db.Select(&servers, query)
	return servers, err
}

func (s *PostgresStore) ApproveServer(serverID string) error {
	query := `
		UPDATE servers
		SET status = 'approved'
		WHERE id = $1 AND status = 'pending'`

	result, err := s.db.Exec(query, serverID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no server found with id %s in pending state", serverID)
	}

	return nil
}

func (s *PostgresStore) FindServerBySerialNumber(serialNumber string) (*models.Server, error) {
	server := &models.Server{}
	row := s.db.QueryRow(
		`SELECT id, serial_number, ip_address, status, created_at
         FROM servers WHERE serial_number = $1`,
		serialNumber,
	)

	err := row.Scan(&server.ID, &server.SerialNumber, &server.IPAddress,
		&server.Status, &server.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrServerNotFound
	}
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
