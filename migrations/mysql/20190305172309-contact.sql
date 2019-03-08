
-- +migrate Up
create table IF NOT EXISTS contact_message (
   `oid` INT AUTO_INCREMENT PRIMARY KEY,
   returnaddr TEXT,
   subject TEXT,
   body TEXT,
   created_at DATETIME,
   updated_at DATETIME NULL,
   deleted_at DATETIME NULL
);

-- +migrate Down
DROP TABLE IF EXISTS contact_message;
