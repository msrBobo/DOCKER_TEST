CURRENT_DIR=$(shell pwd)
DB_URL := "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"

proto-gen:
	chmod +x ./scripts/genproto.sh
	./scripts/genproto.sh
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-force:
	migrate -path migrations -database "$(DB_URL)" -verbose force 1

migrate-file:
	migrate create -ext sql -dir migrations/ -seq create_comments_table


# -- CREATE TABLE IF NOT EXISTS users (
# --     id UUID NOT NULL PRIMARY KEY,
# --     first_name VARCHAR(64) NOT NULL,
# --     last_name VARCHAR(64) NOT NULL,
# --     email VARCHAR(64) NOT NULL  UNIQUE,
# --     password TEXT NOT NULL,
# --     refresh_token TEXT NOT NULL DEFAULT '',
# --     created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
# --     updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
# --     deleted_at TIMESTAMP WITHOUT TIME ZONE
# -- );