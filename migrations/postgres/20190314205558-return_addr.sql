
-- +migrate Up
ALTER TABLE contact_message RENAME COLUMN returnaddr TO return_addr;

-- +migrate Down
ALTER TABLE contact_message RENAME COLUMN return_addr TO returnaddr;
