DB_CONNECTION_STRING ?= postgres://bootandgo:bootandgo@localhost:5432/fleet?sslmode=disable

migrate-up:
	goose -dir migrations postgres "${DB_CONNECTION_STRING}" up

migrate-down:
	goose -dir migrations postgres "${DB_CONNECTION_STRING}" down

migrate-status:
	goose -dir migrations postgres "${DB_CONNECTION_STRING}" status
