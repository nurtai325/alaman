-- +goose Up
-- +goose StatementBegin
CREATE TABLE sales(
	id SERIAL PRIMARY KEY,
	type VARCHAR(100) NOT NULL,
	delivery_type VARCHAR(100),
	payment_at TIMESTAMPTZ,
	full_sum REAL NOT NULL,
	delivery_cost REAL NOT NULL,
	loan_cost REAL NOT NULL,
	items_sum REAL NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE sale_items(
	id SERIAL PRIMARY KEY,
	price REAL NOT NULL,
	product_name VARCHAR(100) NOT NULL,
	sale_count INT NOT NULL,
	quantity INT NOT NULL,
	sale_id INT NOT NULL,
	product_id INT NOT NULL,
	FOREIGN KEY(sale_id)
	REFERENCES sales(id),
	FOREIGN KEY(product_id)
	REFERENCES products(id),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE leads(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100),
	address TEXT,
	phone VARCHAR(15) NOT NULL,
	completed BOOLEAN NOT NULL DEFAULT false,
	user_id INT,
	sale_id INT,
	FOREIGN KEY(user_id)
	REFERENCES users(id),
	FOREIGN KEY(sale_id)
	REFERENCES sales(id),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	sold_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sale_items;
DROP TABLE IF EXISTS leads;
DROP TABLE IF EXISTS sales;
-- +goose StatementEnd
