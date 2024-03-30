CREATE TABLE domain_shortening_counts (
    long_domain_url TEXT NOT NULL,
    shortening_count BIGINT UNSIGNED NOT NULL,
    FULLTEXT  INDEX (long_domain_url)
);