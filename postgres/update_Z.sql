BEGIN;

CREATE TABLE IF NOT EXISTS stocks_staging (LIKE stocks INCLUDING ALL);
INSERT INTO stocks_staging SELECT * FROM stocks WHERE NOT EXISTS (SELECT 1 FROM stocks_staging LIMIT 1);

CREATE TABLE IF NOT EXISTS news_articles_staging (LIKE news_articles INCLUDING ALL);
INSERT INTO news_articles_staging SELECT * FROM news_articles WHERE NOT EXISTS (SELECT 1 FROM news_articles_staging LIMIT 1);

CREATE TABLE IF NOT EXISTS stock_watchlist_staging (LIKE stock_watchlist INCLUDING ALL);
INSERT INTO stock_watchlist_staging SELECT * FROM stock_watchlist WHERE NOT EXISTS (SELECT 1 FROM stock_watchlist_staging LIMIT 1);

CREATE TABLE IF NOT EXISTS holdings_staging (LIKE holdings INCLUDING ALL);
INSERT INTO holdings_staging SELECT * FROM holdings WHERE NOT EXISTS (SELECT 1 FROM holdings_staging LIMIT 1);

CREATE TABLE IF NOT EXISTS orders_staging (LIKE orders INCLUDING ALL);
INSERT INTO orders_staging SELECT * FROM orders WHERE NOT EXISTS (SELECT 1 FROM orders_staging LIMIT 1);

-- Only alter if column isn't already the target type
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'stocks_staging'
          AND column_name = 'overall_sentiment_score'
          AND data_type <> 'numeric'
    ) THEN
        ALTER TABLE stocks_staging
            ALTER COLUMN overall_sentiment_score TYPE DECIMAL(3,2)
            USING overall_sentiment_score::DECIMAL(3,2);
    END IF;
END $$;

ALTER TABLE news_articles_staging
    ADD COLUMN IF NOT EXISTS finnhub_news_id TEXT;

CREATE INDEX IF NOT EXISTS idx_news_articles_staging_finnhub_news_id
    ON news_articles_staging(finnhub_news_id);

CREATE INDEX IF NOT EXISTS idx_news_articles_staging_ticker
    ON news_articles_staging(ticker);

CREATE INDEX IF NOT EXISTS idx_news_articles_staging_publication_time
    ON news_articles_staging(publication_time DESC);

CREATE INDEX IF NOT EXISTS idx_news_articles_staging_ticker_pub_time
    ON news_articles_staging(ticker, publication_time DESC);

CREATE INDEX IF NOT EXISTS idx_holdings_staging_user_id
    ON holdings_staging(user_id);

CREATE INDEX IF NOT EXISTS idx_holdings_staging_stock_id
    ON holdings_staging(stock_id);

CREATE INDEX IF NOT EXISTS idx_orders_staging_user_id
    ON orders_staging(user_id);

CREATE INDEX IF NOT EXISTS idx_orders_staging_stock_id
    ON orders_staging(stock_id);

CREATE INDEX IF NOT EXISTS idx_orders_staging_created_at
    ON orders_staging(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_stock_watchlist_staging_user_id
    ON stock_watchlist_staging(user_id);

CREATE INDEX IF NOT EXISTS idx_stock_watchlist_staging_stock_id
    ON stock_watchlist_staging(stock_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_active_watchlist_staging_entry
    ON stock_watchlist_staging(user_id, stock_id)
    WHERE is_active = TRUE;

COMMIT;