
-- +migrate Up
create table IF NOT EXISTS protest (
   `oid` INT AUTO_INCREMENT PRIMARY KEY,
   title TEXT,
   protest TEXT,
   description TEXT,
   gather_at DATETIME,
   created_at DATETIME,
   updated_at DATETIME NULL,
   deleted_at DATETIME NULL
);
create table IF NOT EXISTS step (
   `oid` INT AUTO_INCREMENT PRIMARY KEY,
   protest_id INTEGER,
   place TEXT,
   details TEXT,
   gather_at DATETIME,
   lat DECIMAL(10,8),
   lng DECIMAL(11,8)
);

-- +migrate Down
DROP TABLE IF EXISTS step;
DROP TABLE IF EXISTS protest;
