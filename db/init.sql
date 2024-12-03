CREATE TABLE public.orders (
  uuid uuid NOT NULL,
  details jsonb NOT NULL,
  created_at timestamp NOT NULL
);
CREATE USER user_table_orders WITH PASSWORD 'myPassword';
GRANT ALL PRIVILEGES ON DATABASE docker to user_table_orders;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO user_table_orders;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO user_table_orders;
GRANT CONNECT ON DATABASE docker TO user_table_orders;
