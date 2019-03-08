
-- +migrate Up
ALTER TABLE protest ADD interests INT;
ALTER TABLE protest ADD views INT;

-- +migrate Down
