package url_service

import (
	"github.com/prakash-p-3121/main-url-shortener-ms/repository/url_repository"
	"github.com/prakash-p-3121/main-url-shortener-ms/service/url_service/impl"
)

func NewUrlService() UrlService {
	return impl.UrlServiceImpl{UrlRepository: url_repository.NewUrlRepository()}
}
