
-- +migrate Up
ALTER TABLE protest ADD author_id TEXT;

-- +migrate Down
ALTER TABLE protest DROP COLUMN author_id;
