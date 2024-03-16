--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cart; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    added_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.cart OWNER TO postgres;

--
-- Name: cart_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cart_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.cart_id_seq OWNER TO postgres;

--
-- Name: cart_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cart_id_seq OWNED BY public.cart.id;


--
-- Name: cart_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart_items (
    cart_id integer NOT NULL,
    product_id integer NOT NULL,
    quantity integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.cart_items OWNER TO postgres;

--
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
    order_id integer NOT NULL,
    product_id integer NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    price_at_order numeric(10,2) NOT NULL
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    delivery_address character(256),
    order_date timestamp with time zone DEFAULT now()
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: product_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_types (
    id integer NOT NULL,
    name character(64) NOT NULL
);


ALTER TABLE public.product_types OWNER TO postgres;

--
-- Name: product_types_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.product_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.product_types_id_seq OWNER TO postgres;

--
-- Name: product_types_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.product_types_id_seq OWNED BY public.product_types.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer NOT NULL,
    name character(128) NOT NULL,
    price numeric(10,2) NOT NULL,
    manufacturer character(64),
    product_type_id integer
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character(64) NOT NULL,
    email character(64) NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: cart id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart ALTER COLUMN id SET DEFAULT nextval('public.cart_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: product_types id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_types ALTER COLUMN id SET DEFAULT nextval('public.product_types_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);

--
-- Data for Name: product_types; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.product_types (id, name) VALUES
(1, 'Наушники'),
(2, 'Колонки'),
(3, 'Микрофоны'),
(4, 'Кабели');

--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, price, manufacturer, product_type_id) FROM stdin;
1	Sony WH-1000XM4                                                                                                                 	350.00	Sony                                                            	1
2	Bose QC 35 II                                                                                                                   	300.00	Bose                                                            	1
3	Sennheiser HD 560S                                                                                                              	199.99	Sennheiser                                                      	1
4	Audio-Technica ATH-M50x                                                                                                         	149.00	Audio-Technica                                                  	1
5	Beyerdynamic DT 770 PRO                                                                                                         	159.00	Beyerdynamic                                                    	1
6	AKG K702                                                                                                                        	349.00	AKG                                                             	1
7	Philips SHP9500                                                                                                                 	75.00	Philips                                                         	1
8	Grado SR80e                                                                                                                     	125.00	Grado                                                           	1
9	Hifiman Sundara                                                                                                                 	349.00	Hifiman                                                         	1
10	Focal Listen Professional                                                                                                       	299.00	Focal                                                           	1
11	Sony MDR7506                                                                                                                    	100.00	Sony                                                            	1
12	Jabra Elite 85h                                                                                                                 	250.00	Jabra                                                           	1
13	JBL Flip 5                                                                                                                      	120.00	JBL                                                             	2
14	Marshall Emberton                                                                                                               	150.00	Marshall                                                        	2
15	Sony SRS-XB12                                                                                                                   	58.00	Sony                                                            	2
16	Ultimate Ears BOOM 3                                                                                                            	129.99	Ultimate Ears                                                   	2
17	Anker Soundcore 2                                                                                                               	39.99	Anker                                                           	2
18	Bose SoundLink Mini II                                                                                                          	199.00	Bose                                                            	2
19	Sonos Move                                                                                                                      	399.00	Sonos                                                           	2
20	Bang & Olufsen Beoplay A1                                                                                                       	250.00	Bang & Olufsen                                                  	2
21	Yamaha WX-010 MusicCast                                                                                                         	179.95	Yamaha                                                          	2
22	Harman Kardon Onyx Studio 6                                                                                                     	479.00	Harman Kardon                                                   	2
23	JBL Charge 4                                                                                                                    	140.00	JBL                                                             	2
24	Marshall Stockwell II                                                                                                           	219.99	Marshall                                                        	2
25	Shure SM7B                                                                                                                      	400.00	Shure                                                           	3
26	Rode NT-USB                                                                                                                     	170.00	Rode                                                            	3
27	Audio-Technica AT2020                                                                                                           	99.00	Audio-Technica                                                  	3
28	Blue Yeti USB                                                                                                                   	129.99	Blue                                                            	3
29	AKG Pro Audio C214                                                                                                              	399.00	AKG                                                             	3
30	Behringer C-1                                                                                                                   	49.00	Behringer                                                       	3
31	Rode NT1-A                                                                                                                      	229.00	Rode                                                            	3
32	Sennheiser MD 421 II                                                                                                            	379.95	Sennheiser                                                      	3
33	Neumann TLM 102                                                                                                                 	699.00	Neumann                                                         	3
34	Shure Beta 58A                                                                                                                  	159.00	Shure                                                           	3
35	Electro-Voice RE20                                                                                                              	449.00	Electro-Voice                                                   	3
36	Zoom H1n Handy Recorder                                                                                                         	119.99	Zoom                                                            	3
37	Аудиокабель Jack 3.5mm                                                                                                          	10.00	Generic                                                         	4
38	USB Type-C                                                                                                                      	8.00	Generic                                                         	4
39	HDMI 2.0                                                                                                                        	15.00	Generic                                                         	4
40	Lightning to USB                                                                                                                	25.00	Apple                                                           	4
41	DisplayPort to DisplayPort                                                                                                      	20.00	Generic                                                         	4
42	XLR 3-Pin                                                                                                                       	10.00	Generic                                                         	4
43	AUX 3.5mm (10m)                                                                                                                 	12.00	Generic                                                         	4
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password, created_at) FROM stdin;
2e982c82-4169-404d-b059-9f269bb882e0	test1                                                           	test@gmail.com                                                  	$2a$14$LmwDpwvlVj71I/99UA2.v.axGUWkmMh6lIHXnkwEivWW/flbkcXP2	2024-03-09 13:48:45.052859+03
\.


--
-- Name: cart_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cart_id_seq', 1, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 8, true);


--
-- Name: product_types_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.product_types_id_seq', 4, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 43, true);


--
-- Name: cart_items cart_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_pkey PRIMARY KEY (cart_id, product_id);


--
-- Name: cart cart_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (order_id, product_id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: product_types product_types_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_types
    ADD CONSTRAINT product_types_name_key UNIQUE (name);


--
-- Name: product_types product_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_types
    ADD CONSTRAINT product_types_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_cart_items_cart_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_cart_items_cart_id ON public.cart_items USING btree (cart_id);


--
-- Name: idx_cart_items_product_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_cart_items_product_id ON public.cart_items USING btree (product_id);


--
-- Name: idx_products_product_type_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_product_type_id ON public.products USING btree (product_type_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: cart_items cart_items_cart_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.cart(id);


--
-- Name: cart_items cart_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: cart cart_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: products products_product_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES public.product_types(id);


--
-- PostgreSQL database dump complete
--

