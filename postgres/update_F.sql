DROP TABLE IF EXISTS news_articles;
CREATE TABLE IF NOT EXISTS news_articles(
    news_article_id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL,
    article_title TEXT,
    article_url TEXT,
    publication_time TIMESTAMPTZ,
    sentiment_score DECIMAL(3,2) NOT NULL DEFAULT 0,
    article_summary TEXT,
    finnhub_news_id TEXT
);

ALTER TABLE stocks ADD COLUMN overall_sentiment_score DECIMAL(3,2) NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_news_articles_finnhub_news_id ON news_articles(finnhub_news_id); 