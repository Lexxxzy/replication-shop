CREATE KEYSPACE IF NOT EXISTS cassandrakeyspace
            WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

USE cassandrakeyspace;

CREATE TABLE IF NOT EXISTS product_types
(
    id   int PRIMARY KEY,
    name text,
);

CREATE INDEX IF NOT EXISTS product_types ON product_types (name);

CREATE TABLE IF NOT EXISTS products
(
    id              int PRIMARY KEY,
    name            text,
    price           float,
    manufacturer    text,
    product_type_id int
) WITH compaction = {'class': 'SizeTieredCompactionStrategy', 'bucket_high': 1.5 };

CREATE INDEX IF NOT EXISTS product_name_index ON products (name);
CREATE INDEX IF NOT EXISTS product_type_id_index ON products (product_type_id);