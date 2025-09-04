-- Order 1
INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, shardkey, sm_id, date_created, oof_shard)
VALUES ('11111111-1111-1111-1111-111111111111', 'TRACK111', 'WEB', 'ru', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'shard1', 1, '2025-09-03T10:00:00Z', '1');

INSERT INTO deliveries (order_uid, name, phone, city, address, email)
VALUES ('11111111-1111-1111-1111-111111111111', 'Ivan Ivanov', '+79998887766', 'Moscow', 'Red Square 1', 'ivan@example.com');

INSERT INTO payments (order_uid, transaction, currency, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee)
VALUES ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'RUB', 'card', 1200, 1690000001, 200, 1000, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('11111111-1111-1111-1111-111111111111', 1001, 'TRACK111', 1000, 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'T-shirt', 0, 'M', 1000, 5001, 'Nike', 1);

-- Order 2
INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, shardkey, sm_id, date_created, oof_shard)
VALUES ('22222222-2222-2222-2222-222222222222', 'TRACK222', 'WBIL', 'en', 'dddddddd-dddd-dddd-dddd-dddddddddddd', 'shard2', 2, '2025-09-03T11:00:00Z', '2');

INSERT INTO deliveries (order_uid, name, phone, city, address, email)
VALUES ('22222222-2222-2222-2222-222222222222', 'Petr Petrov', '+78120000000', 'Saint Petersburg', 'Nevsky Prospect 10', 'petr@example.com');

INSERT INTO payments (order_uid, transaction, currency, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee)
VALUES ('22222222-2222-2222-2222-222222222222', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'USD', 'sbp', 2000, 1690000002, 300, 1700, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('22222222-2222-2222-2222-222222222222', 2002, 'TRACK222', 1700, 'ffffffff-ffff-ffff-ffff-ffffffffffff', 'Shoes', 0, '42', 1700, 6002, 'Adidas', 1);

-- Order 3
INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, shardkey, sm_id, date_created, oof_shard)
VALUES ('33333333-3333-3333-3333-333333333333', 'TRACK333', 'WEB', 'ru', '11111111-2222-3333-4444-555555555555', 'shard3', 3, '2025-09-03T12:00:00Z', '3');

INSERT INTO deliveries (order_uid, name, phone, city, address, email)
VALUES ('33333333-3333-3333-3333-333333333333', 'Anna Smirnova', '+78451234567', 'Kazan', 'Bauman St 5', 'anna@example.com');

INSERT INTO payments (order_uid, transaction, currency, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee)
VALUES ('33333333-3333-3333-3333-333333333333', '66666666-7777-8888-9999-aaaaaaaaaaaa', 'EUR', 'paypal', 1500, 1690000003, 250, 1250, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('33333333-3333-3333-3333-333333333333', 3003, 'TRACK333', 1250, 'bbbbbbbb-cccc-dddd-eeee-ffffffffffff', 'Bag', 0, 'L', 1250, 7003, 'Gucci', 1);

-- Order 4
INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, shardkey, sm_id, date_created, oof_shard)
VALUES ('44444444-4444-4444-4444-444444444444', 'TRACK444', 'WBIL', 'ru', '66666666-6666-6666-6666-666666666666', 'shard4', 4, '2025-09-03T13:00:00Z', '4');

INSERT INTO deliveries (order_uid, name, phone, city, address, email)
VALUES ('44444444-4444-4444-4444-444444444444', 'Sergey Sidorov', '+73837776655', 'Novosibirsk', 'Lenina 50', 'sergey@example.com');

INSERT INTO payments (order_uid, transaction, currency, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee)
VALUES ('44444444-4444-4444-4444-444444444444', '77777777-7777-7777-7777-777777777777', 'RUB', 'yandex', 500, 1690000004, 100, 400, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('44444444-4444-4444-4444-444444444444', 4004, 'TRACK444', 400, '99999999-aaaa-bbbb-cccc-dddddddddddd', 'Hat', 0, 'L', 400, 8004, 'Puma', 1);

-- Order 5
INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, shardkey, sm_id, date_created, oof_shard)
VALUES ('55555555-5555-5555-5555-555555555555', 'TRACK555', 'WEB', 'en', '77777777-8888-9999-aaaa-bbbbbbbbbbbb', 'shard5', 5, '2025-09-03T14:00:00Z', '5');

INSERT INTO deliveries (order_uid, name, phone, city, address, email)
VALUES ('55555555-5555-5555-5555-555555555555', 'Elena Volkova', '+73432223344', 'Yekaterinburg', 'Malysheva 20', 'elena@example.com');

INSERT INTO payments (order_uid, transaction, currency, provider, amount, payment_dt, delivery_cost, goods_total, custom_fee)
VALUES ('55555555-5555-5555-5555-555555555555', '88888888-8888-8888-8888-888888888888', 'USD', 'qiwi', 3000, 1690000005, 500, 2500, 0);

INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ('55555555-5555-5555-5555-555555555555', 5005, 'TRACK555', 2500, 'aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee', 'Watch', 0, 'M', 2500, 9005, 'Casio', 1);
