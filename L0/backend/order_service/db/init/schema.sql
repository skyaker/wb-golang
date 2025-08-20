CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE orders (
  order_uid        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  track_number     text NOT NULL,
  entry            text NOT NULL,
  locale           text NOT NULL,
  internal_signature text,
  customer_id      UUID NOT NULL,
  delivery_service text,
  shardkey         text NOT NULL,
  sm_id            integer NOT NULL CHECK (sm_id > 0),
  date_created     timestamptz NOT NULL,
  oof_shard        text NOT NULL
);

CREATE TABLE deliveries (
  id        UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_uid UUID NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
  name      text NOT NULL,
  phone     text,
  zip       text,
  city      text,
  address   text NOT NULL,
  region    text,
  email     text,
  CONSTRAINT contact_required CHECK (phone IS NOT NULL OR email IS NOT NULL)
);

CREATE TABLE payments (
  id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_uid     UUID NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
  transaction   UUID NOT NULL UNIQUE,
  request_id    UUID,
  currency      char(3) NOT NULL,
  provider      text NOT NULL,
  amount        integer NOT NULL CHECK (amount >= 0),
  payment_dt    BIGINT NOT NULL,
  bank          text,
  delivery_cost integer NOT NULL CHECK (delivery_cost >= 0),
  goods_total   integer NOT NULL CHECK (goods_total >= 0),
  custom_fee    integer CHECK (custom_fee >= 0) DEFAULT 0
);

CREATE TABLE items (
  id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_uid    UUID NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
  chrt_id      bigint NOT NULL,
  track_number text NOT NULL,
  price        integer NOT NULL CHECK (price >= 0),
  rid          UUID NOT NULL,
  name         text NOT NULL,
  sale         integer NOT NULL CHECK (sale >= 0),
  size         text,
  total_price  integer NOT NULL CHECK (total_price >= 0),
  nm_id        bigint NOT NULL,
  brand        text,
  status       integer CHECK (status >= 0),
  CONSTRAINT unique_items_in_order UNIQUE (order_uid, chrt_id)
);
