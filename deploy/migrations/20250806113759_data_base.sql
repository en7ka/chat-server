-- +goose Up
create table data (
      id SERIAL PRIMARY KEY,
      username VARCHAR(100) NOT NULL UNIQUE,
      email VARCHAR(100) NOT NULL UNIQUE,
      created_at timestamp not null default now()
);

-- +goose Down
drop table data;