-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
