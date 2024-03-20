CREATE KEYSPACE IF NOT EXISTS cassandrakeyspace
            WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE cassandrakeyspace;

CREATE TABLE IF NOT EXISTS product_types
(
    id   int PRIMARY KEY,
    name text,
);

CREATE TABLE IF NOT EXISTS products
(
    id              int PRIMARY KEY,
    name            text,
    price           decimal,
    manufacturer    text,
    product_type_id int
);
