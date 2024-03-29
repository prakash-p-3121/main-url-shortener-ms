package url_model

import (
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/restlib"
)

type ShortenUrlReq struct {
	LongUrl *string `json:"long-url"`
}

func (req *ShortenUrlReq) Validate() errorlib.AppError {
	if req.LongUrl == nil {
		return errorlib.NewBadReqError("long-url-nil")
	}
	if restlib.TrimAndCheckForEmptyString(req.LongUrl) {
		return errorlib.NewBadReqError("long-url-empty")
	}
	return nil
}

type ShortenUrlResp struct {
	ShortUrl string `json:"short-url"`
}
