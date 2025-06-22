DROP TABLE IF EXISTS news_articles;
CREATE TABLE IF NOT EXISTS news_articles(
    news_article_id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL,
    article_title TEXT,
    article_url TEXT,
    publication_time TIMESTAMPTZ,
    sentiment_score DECIMAL(3,2) NOT NULL DEFAULT 0
);

ALTER TABLE stocks ADD COLUMN overall_sentiment_score DECIMAL(3,2) NOT NULL DEFAULT 0;