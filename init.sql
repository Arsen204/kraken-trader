CREATE TABLE orders
(
    order_id uuid NOT NULL,
    order_type text NOT NULL,
    symbol text NOT NULL,
    side text NOT NULL,
    quantity numeric NOT NULL,
    filled numeric NOT NULL,
    order_time timestamp NOT NULL,
    reduced_only boolean NOT NULL,
    PRIMARY KEY (order_id)
);
