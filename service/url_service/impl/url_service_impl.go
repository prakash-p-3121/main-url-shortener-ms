package impl

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	database_clustermgt_client "github.com/prakash-p-3121/database-clustermgt-client"
	"github.com/prakash-p-3121/errorlib"
	"github.com/prakash-p-3121/idgenclient"
	"github.com/prakash-p-3121/main-url-shortener-ms/cfg"
	"github.com/prakash-p-3121/main-url-shortener-ms/database"
	"github.com/prakash-p-3121/main-url-shortener-ms/model/url_model"
	"github.com/prakash-p-3121/main-url-shortener-ms/repository/url_repository"
	"log"
	"net/url"
	"strings"
)

type UrlServiceImpl struct {
	UrlRepository url_repository.UrlRepository
}

const (
	shortUrlDomain string = "http://localhost:3000/"
)

func (service *UrlServiceImpl) ShortenUrl(req *url_model.ShortenUrlReq) (*url_model.ShortenUrlResp, errorlib.AppError) {
	appErr := req.Validate()
	if appErr != nil {
		return nil, appErr
	}

	databaseClstrMgtMsCfg, err := cfg.GetMsConnectionCfg("database-clustermgt-ms")
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	client := database_clustermgt_client.NewDatabaseClusterMgtClient(databaseClstrMgtMsCfg.Host,
		uint(databaseClstrMgtMsCfg.Port))

	shortUrlFound := true
	shortUrl, appErr := service.findShortUrlByLongUrl(client, req.LongUrl)
	if appErr != nil {
		var notFoundErrorImpl *errorlib.NotFoundErrorImpl
		ok := errors.As(appErr, &notFoundErrorImpl)
		if ok {
			shortUrlFound = false
		} else {
			return nil, appErr
		}
	}
	if shortUrlFound {
		log.Println("shortUrlFound=true")
		return &url_model.ShortenUrlResp{ShortUrl: shortUrl.ShortUrl}, nil
	}
	log.Println("ShortUrl NotFound")

	idGenMSCfg, err := cfg.GetMsConnectionCfg("idgenms")
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	idGenClient := idgenclient.NewIDGenClient(idGenMSCfg.Host, uint(idGenMSCfg.Port))
	resp, appErr := idGenClient.NextID(database.ShortUrlsTable)
	if appErr != nil {
		return nil, appErr
	}

	/*
		     	Short URL Generation Algorithm
				1. Compute Long Url Hash
			    2. Build (ShortUrlID, LongUrlHash) CSV
			    3. Base64 URL Encode JSON

	*/

	urlHash, err := service.computeHash(req.LongUrl)
	if err != nil {
		return nil, errorlib.NewInternalServerError("compute-short-url-err=" + err.Error())
	}
	log.Println("url_hash=" + urlHash)

	shortUrl = &url_model.ShortUrl{
		ID:          resp.ID,
		IDBitCount:  uint64(resp.BitCount),
		LongUrl:     *req.LongUrl,
		LongUrlHash: urlHash,
	}
	shortUrlResp, appErr := service.buildShortUrl(shortUrl)
	if appErr != nil {
		return nil, appErr
	}

	shortUrl.ShortUrl = shortUrlResp.ShortUrl
	log.Println("shortUrl=", shortUrl.ShortUrl)

	shardPtr, appErr := client.FindShard(database.ShortUrlsTable, resp.ID)
	if appErr != nil {
		return nil, appErr
	}

	appErr = service.UrlRepository.CreateShortUrl(shardPtr.ID, shortUrl)
	if appErr != nil {
		return nil, appErr
	}

	trimmedLongUrl := service.findLongUrlShardKey(req.LongUrl)
	shardPtr, appErr = client.FindShard(database.LongToShortUrlsMappingsTable, trimmedLongUrl)
	if appErr != nil {
		return nil, appErr
	}
	appErr = service.UrlRepository.CreateLongUrlToShortUrlIDMapping(shardPtr.ID,
		&shortUrl.LongUrl,
		&shortUrl.ID)
	if appErr != nil {
		return nil, appErr
	}

	shortenedUrlDomain, err := service.findDomain(req.LongUrl)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	log.Println("shortened url domain= " + shortenedUrlDomain)
	appErr = service.UrlRepository.IncrShortenedDomainCount(&shortenedUrlDomain)
	if appErr != nil {
		return nil, appErr
	}

	return shortUrlResp, nil
}

func (service *UrlServiceImpl) findDomain(urlStr *string) (string, error) {
	parsedUrl, err := url.Parse(*urlStr)
	if err != nil {
		return "", err
	}
	return parsedUrl.Hostname(), nil
}

func (service *UrlServiceImpl) findLongUrlShardKey(longUrl *string) string {
	var cutPrefixes []string = []string{"https://www.", "http://www.", "https://", "http://"}
	var trimmedLongUrl string = *longUrl
	var exists bool
	for _, cutPrefix := range cutPrefixes {
		trimmedLongUrl, exists = strings.CutPrefix(*longUrl, cutPrefix)
		if exists {
			break
		}
	}
	return trimmedLongUrl
}

func (service *UrlServiceImpl) findShortUrlByLongUrl(client database_clustermgt_client.DatabaseClusterMgtClient,
	longUrl *string) (*url_model.ShortUrl, errorlib.AppError) {

	trimmedLongUrl := service.findLongUrlShardKey(longUrl)

	log.Println("findShortUrlByLongUrl:trimmedLongUrl=", trimmedLongUrl)
	shardPtr, appErr := client.FindShard(database.LongToShortUrlsMappingsTable, trimmedLongUrl)
	if appErr != nil {
		return nil, appErr
	}
	shortUrlID, appErr := service.UrlRepository.FindShortUrlIDByLongUrl(shardPtr.ID, longUrl)
	if appErr != nil {
		return nil, appErr
	}

	log.Println("findShortUrlByLongUrl:shortUrlID=", shortUrlID)

	shardPtr, appErr = client.FindShard(database.ShortUrlsTable, shortUrlID)
	if appErr != nil {
		return nil, appErr
	}
	shortUrl, appErr := service.UrlRepository.FindShortUrlByID(shardPtr.ID, shortUrlID)
	if appErr != nil {
		return nil, appErr
	}
	log.Println("findShortUrlByLongUrl:shortUrl=", shortUrl.ID)

	return shortUrl, nil
}

func (service *UrlServiceImpl) computeHash(longUrl *string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(*longUrl))
	if err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func (service *UrlServiceImpl) buildShortUrl(shortUrl *url_model.ShortUrl) (*url_model.ShortenUrlResp, errorlib.AppError) {
	shortUrlCombinedCmpnts := fmt.Sprintf("%s,%s", shortUrl.ID, shortUrl.LongUrlHash)
	urlEncStr := base64.URLEncoding.EncodeToString([]byte(shortUrlCombinedCmpnts))
	log.Println("UrlEncoded Str=", urlEncStr)

	shortUrlStr := fmt.Sprintf(shortUrlDomain+"%s", urlEncStr)
	return &url_model.ShortenUrlResp{ShortUrl: shortUrlStr}, nil
}

func (service *UrlServiceImpl) FindLongUrl(encodedShortUrl *string) (*url_model.FindLongUrlResp, errorlib.AppError) {

	/*
		Long URL Finding Algorithm

		1. Base64 URL Decode
		2. CSV Components  to Struct
		3. Access database to find the Long URL

	*/

	log.Println("EncodedShortUrl=", encodedShortUrl)

	csvUrl, err := base64.URLEncoding.DecodeString(*encodedShortUrl)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	var components url_model.ShortUrlComponents
	componentList := strings.SplitN(string(csvUrl), ",", 2)
	components.ID = componentList[0]
	components.UrlHash = componentList[1]

	log.Println("Components###")
	log.Println(components)

	databaseClstrMgtMsCfg, err := cfg.GetMsConnectionCfg("database-clustermgt-ms")
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	client := database_clustermgt_client.NewDatabaseClusterMgtClient(databaseClstrMgtMsCfg.Host,
		uint(databaseClstrMgtMsCfg.Port))

	shardPtr, appErr := client.FindShard(database.ShortUrlsTable, components.ID)
	if appErr != nil {
		return nil, appErr
	}

	shortUrl, appErr := service.UrlRepository.FindShortUrlByID(shardPtr.ID, components.ID)
	if appErr != nil {
		return nil, appErr
	}

	if shortUrl.LongUrlHash != components.UrlHash {
		return nil, errorlib.NewNotFoundError("long-url")
	}
	resp := url_model.FindLongUrlResp{LongUrl: shortUrl.LongUrl}

	return &resp, nil
}

func (service *UrlServiceImpl) FindTopDomains(count uint64) ([]*url_model.DomainCount, errorlib.AppError) {
	if count == 0 {
		badReqErr := errorlib.NewBadReqError("count==0")
		return nil, badReqErr
	}
	return service.UrlRepository.FindTopDomains(count)
}
