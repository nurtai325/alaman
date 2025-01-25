-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	phone VARCHAR(15) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	role VARCHAR(50) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
