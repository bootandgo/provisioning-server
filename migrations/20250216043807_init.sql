-- +goose Up
CREATE TABLE IF NOT EXISTS admins (
    username VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS servers (
    id UUID PRIMARY KEY,
    serial_number VARCHAR(255) NOT NULL UNIQUE,
    ip_address VARCHAR(15) NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'approved')),
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_servers_status ON servers(status);

-- +goose Down
DROP TABLE IF EXISTS servers;
DROP TABLE IF EXISTS admins;
