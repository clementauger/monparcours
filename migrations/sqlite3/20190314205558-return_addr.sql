
-- +migrate Up
PRAGMA foreign_keys=off;

DROP TABLE IF EXISTS _contact_message_old;

ALTER TABLE contact_message RENAME TO _contact_message_old;

CREATE TABLE contact_message (
  return_addr TEXT,
  subject TEXT,
  body TEXT,
  created_at DATETIME,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL
);

INSERT INTO contact_message (return_addr, subject, body, created_at, updated_at, deleted_at)
  SELECT returnaddr, subject, body, created_at, updated_at, deleted_at
  FROM _contact_message_old;

DROP TABLE IF EXISTS _contact_message_old;

PRAGMA foreign_keys=on;

-- +migrate Down

PRAGMA foreign_keys=off;

DROP TABLE IF EXISTS _contact_message_old;

ALTER TABLE contact_message RENAME TO _contact_message_old;

CREATE TABLE contact_message (
  returnaddr TEXT,
  subject TEXT,
  body TEXT,
  created_at DATETIME,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL
);

INSERT INTO contact_message (returnaddr, subject, body, created_at, updated_at, deleted_at)
  SELECT return_addr, subject, body, created_at, updated_at, deleted_at
  FROM _contact_message_old;

DROP TABLE IF EXISTS _contact_message_old;

PRAGMA foreign_keys=on;
