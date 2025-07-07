DROP TABLE IS EXISTS stock_ohlcv;
CREATE TABLE IF NOT EXISTS stock_ohlcv(
    stock_ohlcv_id SERIAL PRIMARY KEY,
    stock_name TEXT NOT NULL,
    timestamp TIMESTAMPTZ,
    open DECIMAL(10,4) NOT NULL DEFAULT 0,
    high DECIMAL(10,4) NOT NULL DEFAULT 0,
    low DECIMAL(10,4) NOT NULL DEFAULT 0,
    close DECIMAL(10,4) NOT NULL DEFAULT 0,
    volume INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT unique_stock_ohlcv_timestamp UNIQUE (stock_name, timestamp)
);