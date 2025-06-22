
-- Table for Users
DROP TABLE IF EXISTS users;
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL DEFAULT 'default_user', -- For V1, we can have a default user
    email TEXT UNIQUE,                                  -- Optional for V1, can be NULL
    hashed_password TEXT,                               -- Not used in V1
    cash_balance_cents BIGINT NOT NULL DEFAULT 0, -- e.g., $10,000.00 stored as 1,000,000 cents
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    notifications_on BOOLEAN DEFAULT FALSE
);

-- Table for Mock Stocks
DROP TABLE IF EXISTS stocks;
CREATE TABLE IF NOT EXISTS stocks (
    stock_id SERIAL PRIMARY KEY,
    ticker TEXT UNIQUE NOT NULL,                        -- e.g., "FAKE_AAPL"
    name TEXT NOT NULL,                                 -- e.g., "Fake Apple Inc."
    opening_price_cents BIGINT,                         -- For V2: daily opening price
    current_price_cents BIGINT,
    min_price_generator_cents BIGINT,                   -- For V2: lower bound for dynamic price generator
    max_price_generator_cents BIGINT,                   -- For V2: upper bound for dynamic price generator
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    overall_sentiment_score INTEGER NOT NULL DEAFULT 0
);

-- Table for User's Portfolio Holdings (Current Stock Positions)
DROP TABLE IF EXISTS holdings;
CREATE TABLE IF NOT EXISTS holdings (
    holding_id SERIAL PRIMARY KEY,                      -- Surrogate key
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(stock_id) ON DELETE RESTRICT,
    quantity BIGINT NOT NULL DEFAULT 0,
    average_cost_per_share_cents BIGINT NOT NULL DEFAULT 0, -- Crucial for V2 P&L. For V1, can be set to buy price.
    created_at TIMESTAMPTZ DEFAULT NOW(),              -- When the holding was first initiated
    updated_at TIMESTAMPTZ DEFAULT NOW(),              -- When quantity or avg_cost was last changed
    CONSTRAINT unique_user_stock_holding UNIQUE (user_id, stock_id) -- Ensures one holding record per user per stock
);

-- Table for Transaction History (Log of all executed buy/sell orders)
DROP TABLE IF EXISTS orders;
CREATE TABLE IF NOT EXISTS orders (
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(stock_id) ON DELETE RESTRICT,
    trade_type TEXT NOT NULL,
    order_status TEXT NOT NULL,
    quantity BIGINT NOT NULL CHECK (quantity > 0),
    price_per_share_cents BIGINT NOT NULL,
    total_order_value_cents BIGINT NOT NULL,      -- Calculated: quantity * price_per_share_cents_at_execution
    created_at TIMESTAMPTZ DEFAULT NOW(),
    notes TEXT                                          -- Optional, for any specific details
);

-- Optional: Indexes for frequently queried columns (PostgreSQL automatically creates indexes for PRIMARY KEY and UNIQUE constraints)
-- Consider adding indexes on foreign keys and columns used in WHERE clauses or ORDER BY for performance as your data grows.
-- Example:
-- CREATE INDEX IF NOT EXISTS idx_portfolio_holdings_user_id ON portfolio_holdings(user_id);
-- CREATE INDEX IF NOT EXISTS idx_portfolio_holdings_stock_id ON portfolio_holdings(stock_id);
-- CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
-- CREATE INDEX IF NOT EXISTS idx_transactions_stock_id ON transactions(stock_id);
-- CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp DESC);


-- Initial Data for V1 MVP

-- Insert a default user for the MVP
INSERT INTO users (username, email, cash_balance_cents)
VALUES ('default_user', 'user@example.com', 10000000)
    ON CONFLICT (username) DO NOTHING;

-- Insert some mock stocks for V1
INSERT INTO stocks (ticker, name, opening_price_cents, current_price_cents, min_price_generator_cents, max_price_generator_cents)
VALUES
    ('AAPL', 'Apple Inc.', 15000, 15000, 14000, 16000),    -- $175.50
    ('GOOGL', 'Alphabet Inc.', 25000, 25000, 24000, 26000), -- $2800.25
    ('MSFT', 'Microsoft Corp.', 30000, 30000, 29000, 31000),   -- $330.75
    ('AMZN', 'Amazon.com Inc.', 35000, 35000, 34000, 36000),   -- $3300.45
    ('TSLA', 'Tesla Inc.', 20000, 20000, 19000, 21000)      -- $250.00
    ON CONFLICT (ticker) DO NOTHING;

-- Function to automatically update 'updated_at' columns
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for 'users' table
CREATE OR REPLACE TRIGGER set_timestamp_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Triggers for 'stocks' table
CREATE OR REPLACE TRIGGER set_timestamp_stocks
    BEFORE UPDATE ON stocks
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Triggers for 'portfolio_holdings' table
CREATE OR REPLACE TRIGGER set_timestamp_holdings
    BEFORE UPDATE ON holdings
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

-- Note: `transactions` table typically doesn't have an `updated_at` as transactions are immutable once created.
-- Its `timestamp` field records the creation time.

DROP TABLE IF EXISTS stock_watchlist;
CREATE TABLE IF NOT EXISTS stock_watchlist(
    stock_watchlist_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    stock_id INTEGER NOT NULL REFERENCES stocks(stock_id) ON DELETE RESTRICT,
    target_price_cents BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

DROP TABLE IF EXISTS news_articles;
CREATE TABLE IF NOT EXISTS news_articles(
    news_article_id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL,
    article_title TEXT,
    article_url TEXT,
    publication_time TIMESTAMPTZ,
    sentiment_score DECIMAL(3,2) NOT NULL DEFAULT 0
);
