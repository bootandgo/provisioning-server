-- +goose Up
ALTER TABLE admins ADD COLUMN is_root BOOLEAN NOT NULL DEFAULT FALSE;

-- Create invitations table
CREATE TABLE invitations (
    token TEXT PRIMARY KEY,
    created_by TEXT NOT NULL REFERENCES admins(username),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create index for faster lookups
CREATE INDEX idx_invitations_expires ON invitations(expires_at);
CREATE INDEX idx_invitations_used ON invitations(used);

-- +goose Down
DROP TABLE invitations;
ALTER TABLE admins DROP COLUMN is_root;
