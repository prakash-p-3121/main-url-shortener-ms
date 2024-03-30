CREATE TABLE long_to_short_url_mappings (
   long_url TEXT NOT NULL,
   short_url_id VARBINARY(2000) NOT NULL,
   FULLTEXT INDEX(long_url)
);