
-- +migrate Up
ALTER TABLE protest ADD organizer TEXT;

-- +migrate Down
