-- +goose Up
-- +goose StatementBegin

-- Up migration
CREATE SCHEMA IF NOT EXISTS chat;

CREATE TABLE chat.chats(
                           id BIGSERIAL PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           type VARCHAR(50) NOT NULL,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE chat.chat_members(
                          id BIGSERIAL PRIMARY KEY,
                          chat_id BIGINT NOT NULL REFERENCES chat.chats(id) ON DELETE CASCADE,
                          user_id BIGINT NOT NULL,
                          joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          is_admin BOOLEAN DEFAULT FALSE,
                          UNIQUE (chat_id, user_id)
);

CREATE TABLE chat.messages(
                          id BIGSERIAL PRIMARY KEY,
                          chat_id BIGINT NOT NULL REFERENCES chat.chats(id) ON DELETE CASCADE,
                          from_user_id BIGINT NOT NULL,
                          text TEXT NOT NULL,
                          timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
-- +goose Down


-- +goose StatementBegin

-- Down migration
DROP TABLE IF EXISTS chat.messages;
DROP TABLE IF EXISTS chat.chat_members;
DROP TABLE IF EXISTS chat.chats;

-- +goose StatementEnd