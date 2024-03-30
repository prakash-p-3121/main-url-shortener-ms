package impl

import (
	"database/sql"
	"errors"
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/mysqllib"
	"sync"
)

type UrlRepositoryImpl struct {
	ShardConnectionsMap *sync.Map
}

func (repository *UrlRepositoryImpl) CreateShortUrl(shardID *int64, req *url_model.ShortUrl) errorlib.AppError {
	db, err := mysqllib.RetrieveShardConnectionByShardID(repository.ShardConnectionsMap, *shardID)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	qry := `INSERT INTO short_urls (id, id_bit_count, long_url, long_url_hash, short_url) VALUES (?, ?, ?, ?, ?);`
	_, err = db.Exec(qry, req.ID, req.IDBitCount, req.LongUrl, req.LongUrlHash, req.ShortUrl)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	return nil
}

func (repository *UrlRepositoryImpl) FindShortUrlIDByLongUrl(shardID *int64, longUrl *string) (string, errorlib.AppError) {
	db, err := mysqllib.RetrieveShardConnectionByShardID(repository.ShardConnectionsMap, *shardID)
	if err != nil {
		return "", errorlib.NewInternalServerError(err.Error())
	}

	qry := `SELECT short_url_id FROM long_to_short_url_mappings WHERE long_url=?;`
	row := db.QueryRow(qry, *longUrl)
	var shortUrlID string
	err = row.Scan(&shortUrlID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errorlib.NewNotFoundError("short-url-id-not-found-for-long-url=" + *longUrl)
	}
	if err != nil {
		return "", errorlib.NewInternalServerError(err.Error())
	}
	return shortUrlID, nil
}

func (repository *UrlRepositoryImpl) FindShortUrlByID(shardID *int64, shortUrlID string) (*url_model.ShortUrl, errorlib.AppError) {
	db, err := mysqllib.RetrieveShardConnectionByShardID(repository.ShardConnectionsMap, *shardID)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	qry := `SELECT id, id_bit_count, long_url, long_url_hash, short_url, created_at FROM short_urls WHERE id=?;`
	row := db.QueryRow(qry, shortUrlID)
	var shortUrl url_model.ShortUrl
	err = row.Scan(
		&shortUrl.ID,
		&shortUrl.IDBitCount,
		&shortUrl.LongUrl,
		&shortUrl.LongUrlHash,
		&shortUrl.ShortUrl,
		&shortUrl.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errorlib.NewNotFoundError("short-url-not-found-for-id=" + shortUrlID)
	}
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	return &shortUrl, nil
}
