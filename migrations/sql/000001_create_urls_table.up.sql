CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    short_url VARCHAR(50) NOT NULL,
    full_url VARCHAR(1000) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_urls_short_url ON urls(short_url);
SELECT setval('urls_id_seq', 100000, true);
