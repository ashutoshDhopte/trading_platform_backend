
DROP TABLE IF EXISTS stock_watchlist;
CREATE TABLE IF NOT EXISTS stock_watchlist(
    stock_watchlist_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(stock_id) ON DELETE RESTRICT,
    target_price_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);