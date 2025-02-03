-- +goose Up
-- +goose StatementBegin
CREATE TABLE reports(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	path VARCHAR(100) NOT NULL,
	start_at TIMESTAMPTZ NOT NULL,
	end_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reports;
-- +goose StatementEnd
