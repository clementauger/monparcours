
-- +migrate Up
ALTER TABLE protest ADD public INT;
ALTER TABLE protest ADD password TEXT;

-- +migrate Down
ALTER TABLE protest DROP COLUMN public;
ALTER TABLE protest DROP COLUMN password;
