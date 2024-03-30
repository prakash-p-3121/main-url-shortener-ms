package url_repository

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
)

type UrlRepository interface {
	CreateShortUrl(shardID *int64, req *url_model.ShortUrl) errorlib.AppError
	CreateLongUrlToShortUrlIDMapping(shardID *int64, longUrl, shortUrlID *string) errorlib.AppError
	FindShortUrlIDByLongUrl(shardID *int64, longUrl *string) (string, errorlib.AppError)
	FindShortUrlByID(shardID *int64, shortUrlID string) (*url_model.ShortUrl, errorlib.AppError)
}
