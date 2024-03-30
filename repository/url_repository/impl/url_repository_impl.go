package impl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/mysqllib"
	"log"
	"sync"
)

type UrlRepositoryImpl struct {
	ShardConnectionsMap   *sync.Map
	SingleStoreConnection *sql.DB
}

func (repository *UrlRepositoryImpl) CreateShortUrl(shardID *int64, req *url_model.ShortUrl) errorlib.AppError {
	db, err := mysqllib.RetrieveShardConnectionByShardID(repository.ShardConnectionsMap, *shardID)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	log.Println("SHORT_URL_RECORD=", *req)
	qry := `INSERT INTO short_urls (id, id_bit_count, long_url, long_url_hash, short_url) VALUES (?, ?, ?, ?, ?);`
	_, err = db.Exec(qry, req.ID, req.IDBitCount, req.LongUrl, req.LongUrlHash, req.ShortUrl)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	return nil
}

func (repository *UrlRepositoryImpl) CreateLongUrlToShortUrlIDMapping(shardID *int64, longUrl, shortUrlID *string) errorlib.AppError {
	db, err := mysqllib.RetrieveShardConnectionByShardID(repository.ShardConnectionsMap, *shardID)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	qry := `INSERT INTO long_to_short_url_mappings (long_url, short_url_id) SELECT ?, ? WHERE NOT 
    		EXISTS ( SELECT 1 FROM long_to_short_url_mappings WHERE long_url = ? ); `
	_, err = db.Exec(qry, *longUrl, *shortUrlID, *longUrl)
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

func (repository *UrlRepositoryImpl) IncrShortenedDomainCount(domain *string) errorlib.AppError {
	db := repository.SingleStoreConnection
	qry := `INSERT INTO domain_shortening_counts (long_url_domain, shortening_count) VALUES (?, ?) 
            ON DUPLICATE KEY UPDATE shortening_count=shortening_count+1 ;`
	_, err := db.Exec(qry, domain, 1)
	if err != nil {
		return errorlib.NewInternalServerError(err.Error())
	}
	return nil
}

func (repository *UrlRepositoryImpl) FindTopDomains(count uint64) ([]*url_model.DomainCount, errorlib.AppError) {
	db := repository.SingleStoreConnection
	qry := "SELECT long_url_domain, shortening_count FROM domain_shortening_counts ORDER BY shortening_count DESC LIMIT %d ;"
	qry = fmt.Sprintf(qry, count)
	rows, err := db.Query(qry)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	resultList := make([]*url_model.DomainCount, 0)
	for rows.Next() {
		var result url_model.DomainCount
		err := rows.Scan(&result.DomainUrl, &result.ShortenedCount)
		if err != nil {
			return nil, errorlib.NewInternalServerError(err.Error())
		}
		resultList = append(resultList, &result)
	}
	return resultList, nil
}
