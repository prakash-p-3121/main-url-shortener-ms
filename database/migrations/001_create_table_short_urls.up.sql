CREATE TABLE short_urls (
    id VARBINARY(2000) NOT NULL PRIMARY KEY,
    id_bit_count BIGINT UNSIGNED NOT NULL,
    long_url TEXT NOT NULL,
    long_url_hash TEXT NOT NULL,
    short_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);