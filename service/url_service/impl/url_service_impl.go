package impl

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/main-url-shortener-ms/repository/url_repository"
)

type UrlServiceImpl struct {
	UrlRepository url_repository.UrlRepository
}

func (service *UrlServiceImpl) ShortenUrl(req *url_model.ShortenUrlReq) errorlib.AppError {
	appErr := req.Validate()
	if appErr != nil {
		return appErr
	}

	return nil
}
