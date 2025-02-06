-- +goose Up
-- +goose StatementBegin
CREATE TABLE lead_whs(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	phone VARCHAR(15) UNIQUE NOT NULL,
	jid VARCHAR(100),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lead_whs;
-- +goose StatementEnd
