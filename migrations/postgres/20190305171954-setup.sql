
-- +migrate Up
create table IF NOT EXISTS protest (
   oid  SERIAL PRIMARY KEY,
   title TEXT,
   protest TEXT,
   description TEXT,
   gather_at timestamp,
   created_at timestamp,
   updated_at timestamp NULL,
   deleted_at timestamp NULL
);
create table IF NOT EXISTS step (
   oid  SERIAL PRIMARY KEY,
   protest_id INTEGER REFERENCES protest(oid),
   place TEXT,
   details TEXT,
   gather_at timestamp,
   lat DECIMAL(10,8),
   lng DECIMAL(11,8)
);

-- +migrate Down
DROP TABLE IF EXISTS step;
DROP TABLE IF EXISTS protest;
