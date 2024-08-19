-- Buat tabel 'products'
CREATE TABLE products
(
    id          SERIAL PRIMARY KEY,
    merchant_id INTEGER NOT NULL,
    image       VARCHAR(255),
    name        VARCHAR(255),
    qty         INTEGER,
    price       DOUBLE PRECISION,
    CONSTRAINT fk_merchant
        FOREIGN KEY (merchant_id)
            REFERENCES merchants (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_products_merchant_id ON products (merchant_id);
