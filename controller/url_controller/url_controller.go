package product_controller

import "github.com/prakash-p-3121/restlib"

type UrlController interface {
	ShortenUrl(ctx restlib.RestContext)
	FindLongUrl(ctx restlib.RestContext)
	FindTopDomains(ctx restlib.RestContext)
}
