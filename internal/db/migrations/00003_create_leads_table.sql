-- +goose Up
-- +goose StatementBegin
CREATE TABLE sales(
	id SERIAL PRIMARY KEY,
	full_price NUMERIC(15, 2),
	type VARCHAR(100),
	items INT[],
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
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS leads;
DROP TABLE IF EXISTS sales;
-- +goose StatementEnd
