-- +goose Up
-- +goose StatementBegin
ALTER TABLE leads
ADD COLUMN first_photo VARCHAR(255) NOT NULL DEFAULT 'assets/photo_not_found.html',
ADD COLUMN second_photo VARCHAR(255) NOT NULL DEFAULT 'assets/photo_not_found.html';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE leads
DROP COLUMN first_photo,
DROP COLUMN second_photo;
-- +goose StatementEnd
