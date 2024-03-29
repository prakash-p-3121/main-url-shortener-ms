package url_repository

import (
	"github.com/prakash-p-3121/main-url-shortener-ms/database"
	"github.com/prakash-p-3121/main-url-shortener-ms/repository/url_repository/impl"
)

func NewUrlRepository() UrlRepository {
	return impl.UrlRepositoryImpl{ShardConnectionsMap: database.GetShardConnectionsMap()}
}
