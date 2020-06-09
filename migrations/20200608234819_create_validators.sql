-- +goose Up
CREATE TABLE validators (
  id         BIGSERIAL NOT NULL PRIMARY KEY,
  address    TEXT NOT NULL,
  name       TEXT,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE validators;
