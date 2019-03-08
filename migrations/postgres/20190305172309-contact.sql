
-- +migrate Up
create table IF NOT EXISTS contact_message (
   oid  SERIAL PRIMARY KEY,
   returnaddr TEXT,
   subject TEXT,
   body TEXT,
   created_at timestamp,
   updated_at timestamp NULL,
   deleted_at timestamp NULL
);

-- +migrate Down
DROP TABLE IF EXISTS contact_message;
