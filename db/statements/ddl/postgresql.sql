-- Create tables
CREATE TABLE public.cart
(
    id       integer NOT NULL,
    user_id  uuid    NOT NULL,
    added_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE TABLE public.cart_items
(
    cart_id    integer           NOT NULL,
    product_id integer           NOT NULL,
    quantity   integer DEFAULT 1 NOT NULL,
    PRIMARY KEY (cart_id, product_id)
);

CREATE TABLE public.order_items
(
    order_id       integer           NOT NULL,
    product_id     integer           NOT NULL,
    quantity       integer DEFAULT 1 NOT NULL,
    price_at_order numeric(10, 2)    NOT NULL,
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE public.orders
(
    id               integer NOT NULL,
    user_id          uuid    NOT NULL,
    delivery_address character(256),
    order_date       timestamp with time zone DEFAULT now(),
    PRIMARY KEY (id)
);

CREATE TABLE public.product_types
(
    id   integer       NOT NULL,
    name character(64) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name)
);

CREATE TABLE public.products
(
    id              integer        NOT NULL,
    name            character(128) NOT NULL,
    price           numeric(10, 2) NOT NULL,
    manufacturer    character(64),
    product_type_id integer,
    PRIMARY KEY (id),
    FOREIGN KEY (product_type_id) REFERENCES public.product_types (id)
);

CREATE TABLE public.users
(
    id         uuid                     DEFAULT gen_random_uuid() NOT NULL,
    name       character(64)                                      NOT NULL,
    email      character(64)                                      NOT NULL,
    password   text                                               NOT NULL,
    created_at timestamp with time zone DEFAULT now()             NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (email)
);

-- Create sequences
CREATE SEQUENCE public.cart_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE public.product_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Set default values for sequences
SELECT pg_catalog.setval('public.cart_id_seq', 1, true);
SELECT pg_catalog.setval('public.orders_id_seq', 8, true);
SELECT pg_catalog.setval('public.product_types_id_seq', 4, true);
SELECT pg_catalog.setval('public.products_id_seq', 43, true);
