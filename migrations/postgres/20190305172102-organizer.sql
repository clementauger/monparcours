
-- +migrate Up
ALTER TABLE protest ADD organizer TEXT;

-- +migrate Down
ALTER TABLE protest DROP COLUMN organizer;
