-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_changes(
	id SERIAL PRIMARY KEY,
	quantity INT NOT NULL,
	is_income BOOLEAN NOT NULL,
	product_id INT NOT NULL,
	FOREIGN KEY(product_id)
	REFERENCES products(id),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_changes;
-- +goose StatementEnd
