package url_model

import "time"

type ShortUrl struct {
	ID          string    `json:"id"`
	IDBitCount  uint64    `json:"id-bit-count"`
	LongUrl     string    `json:"long-url"`
	LongUrlHash string    `json:"long-url-hash"`
	ShortUrl    string    `json:"short-url"`
	CreatedAt   time.Time `json:"created-at"`
}

type LongUrlMapping struct {
	LongUrl    string `json:"long-url"`
	ShortUrlID string `json:"short-url-id"`
}

type DomainCount struct {
	DomainUrl      string `json:"domain-url"`
	ShortenedCount uint64 `json:"shortened-count"`
}
