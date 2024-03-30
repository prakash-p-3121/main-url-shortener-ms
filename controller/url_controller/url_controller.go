package url_controller

import "github.com/prakash-p-3121/restlib"

type UrlController interface {
	ShortenUrl(ctx restlib.RestContext)
	RedirectToLongUrl(ctx restlib.RestContext)
	FindTopDomains(ctx restlib.RestContext)
}
