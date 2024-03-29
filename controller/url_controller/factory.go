package product_controller

import (
	"github.com/prakash-p-3121/main-url-shortener-ms/controller/url_controller/impl"
	"github.com/prakash-p-3121/main-url-shortener-ms/service/url_service"
)

func NewUrlController() UrlController {
	return &impl.UrlControllerImpl{UrlService: url_service.NewUrlService()}
}
