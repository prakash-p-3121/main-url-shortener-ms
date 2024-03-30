CREATE TABLE domain_shortening_counts (
    long_url_domain TEXT NOT NULL PRIMARY KEY,
    shortening_count BIGINT UNSIGNED NOT NULL
);