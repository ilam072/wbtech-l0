CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    track_number VARCHAR(64) NOT NULL UNIQUE,
    entry VARCHAR(10) NOT NULL,
    locale VARCHAR(10) NOT NULL,
    internal_signature TEXT,
    customer_id VARCHAR(128) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey VARCHAR(2) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(2) NOT NULL
);

CREATE TABLE IF NOT EXISTS delivery (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(64) NOT NULL,
    address VARCHAR(64) NOT NULL,
    region VARCHAR(64) NOT NULL,
    email VARCHAR(128)
);

CREATE TABLE IF NOT EXISTS payment (
    transaction UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    request_id VARCHAR(64),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(32) NOT NULL,
    amount INT NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(20) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id BIGINT PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    track_number VARCHAR(64) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(64) NOT NULL,
    name VARCHAR(128) NOT NULL,
    sale INT NOT NULL CHECK(sale BETWEEN 0 AND 100),
    size VARCHAR(16) NOT NULL,
    total_price INT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(64) NOT NULL,
    status INT NOT NULL
);
