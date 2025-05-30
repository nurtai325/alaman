-- +goose Up
-- +goose StatementBegin
CREATE TABLE products(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) UNIQUE NOT NULL,
	in_stock INT NOT NULL,
	price INT NOT NULL DEFAULT 0,
	sale_count INT NOT NULL DEFAULT 1,
	stock_price INT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

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
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
