package url_service

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
)

type UrlService interface {
	ShortenUrl(req *url_model.ShortenUrlReq) (*url_model.ShortenUrlResp, errorlib.AppError)
	FindLongUrl(urlHash *string) (*url_model.FindLongUrlResp, errorlib.AppError)
	FindTopDomains(count uint64) ([]*url_model.DomainCount, errorlib.AppError)
}
