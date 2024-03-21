COPY public.product_types (id, name) FROM '/statements/dml/dummy_data/product_types.csv' WITH (FORMAT csv, HEADER true);

COPY public.products (id, name, price, manufacturer, product_type_id) FROM '/statements/dml/dummy_data/products.csv' WITH (FORMAT csv, HEADER true);

COPY public.users (id, name, email, password, created_at) FROM '/statements/dml/dummy_data/users.csv' WITH (FORMAT csv, HEADER true);