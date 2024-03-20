USE cassandrakeyspace;

COPY product_types (id, name)
FROM '/statements/dml/dummy_data/product_types.csv' WITH HEADER=TRUE;

COPY products (id, name, price, manufacturer, product_type_id)
FROM '/statements/dml/dummy_data/products.csv' WITH HEADER=TRUE;