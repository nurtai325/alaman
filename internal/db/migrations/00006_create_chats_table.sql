-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats(
	id SERIAL PRIMARY KEY,
	lead_id INT NOT NULL,
	user_id INT NOT NULL,
	FOREIGN KEY(lead_id)
	REFERENCES leads(id),
	FOREIGN KEY(user_id)
	REFERENCES users(id),
	updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE messages(
	id SERIAL PRIMARY KEY,
	text TEXT,
	path VARCHAR(255),
	type VARCHAR(50) NOT NULL,
	is_sent BOOLEAN NOT NULL,
	audio_length INT NOT NULL,
	chat_id INT NOT NULL,
	FOREIGN KEY(chat_id)
	REFERENCES chats(id),
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
