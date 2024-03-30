package impl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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
	"strings"
)

type UrlServiceImpl struct {
	UrlRepository url_repository.UrlRepository
}

const (
	shortUrlDomain    string = "https://infra.cloud/"
	encryptDecryptKey string = "HelloWorld@123"
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
		return &url_model.ShortenUrlResp{ShortUrl: shortUrl.ShortUrl}, nil
	}

	idGenMSCfg, err := cfg.GetMsConnectionCfg("idgenms")
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	idGenClient := idgenclient.NewIDGenClient(idGenMSCfg.Host, uint(idGenMSCfg.Port))
	resp, appErr := idGenClient.NextID(database.ShortUrlsTable)
	if appErr != nil {
		return nil, appErr
	}

	urlHash, err := service.computeHash(req.LongUrl)
	if err != nil {
		return nil, errorlib.NewInternalServerError("compute-short-url-err=" + err.Error())
	}

	shortUrlResp, appErr := service.buildShortUrl(shortUrl)
	if appErr != nil {
		return nil, appErr
	}

	shortUrl = &url_model.ShortUrl{
		ID:          resp.ID,
		IDBitCount:  uint64(resp.BitCount),
		LongUrl:     *req.LongUrl,
		LongUrlHash: urlHash,
		ShortUrl:    shortUrlResp.ShortUrl,
	}

	shardPtr, appErr := client.FindShard(database.ShortUrlsTable, resp.ID)
	if appErr != nil {
		return nil, appErr
	}

	appErr = service.UrlRepository.CreateShortUrl(shardPtr.ID, shortUrl)
	if appErr != nil {
		return nil, appErr
	}

	return shortUrlResp, nil
}

func (service *UrlServiceImpl) findShortUrlByLongUrl(client database_clustermgt_client.DatabaseClusterMgtClient,
	longUrl *string) (*url_model.ShortUrl, errorlib.AppError) {
	trimmedLongUrl, exists := strings.CutPrefix(*longUrl, "https://www.")
	if !exists {
		trimmedLongUrl, _ = strings.CutPrefix(*longUrl, "http://www.")
	}
	shardPtr, appErr := client.FindShard(database.LongToShortUrlsMappingsTable, trimmedLongUrl)
	if appErr != nil {
		return nil, appErr
	}
	shortUrlID, appErr := service.UrlRepository.FindShortUrlIDByLongUrl(shardPtr.ID, longUrl)
	if appErr != nil {
		return nil, appErr
	}

	shardPtr, appErr = client.FindShard(database.ShortUrlsTable, shortUrlID)
	if appErr != nil {
		return nil, appErr
	}
	shortUrl, appErr := service.UrlRepository.FindShortUrlByID(shardPtr.ID, shortUrlID)
	if appErr != nil {
		return nil, appErr
	}
	return shortUrl, nil
}

func (service *UrlServiceImpl) computeHash(longUrl *string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(*longUrl))
	if err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	return string(hash), nil
}

func (service *UrlServiceImpl) buildShortUrl(shortUrl *url_model.ShortUrl) (*url_model.ShortenUrlResp, errorlib.AppError) {
	components := url_model.ShortUrlComponents{
		ID:      shortUrl.ID,
		UrlHash: shortUrl.LongUrlHash,
	}
	jsonData, err := json.Marshal(components)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	cipherText, appErr := service.encrypt(jsonData)
	if appErr != nil {
		return nil, appErr
	}
	urlEncodedStr := base64.URLEncoding.EncodeToString(cipherText)
	log.Println("url_encoded_str=", urlEncodedStr)
	shortUrlStr := fmt.Sprintf(shortUrlDomain+"%s", urlEncodedStr)
	return &url_model.ShortenUrlResp{ShortUrl: shortUrlStr}, nil
}

func (service *UrlServiceImpl) encrypt(jsonData []byte) ([]byte, errorlib.AppError) {
	iv := make([]byte, aes.BlockSize) // Random initialization vector
	_, err := rand.Read(iv)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	block, err := aes.NewCipher([]byte(encryptDecryptKey))
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	plainText := jsonData
	cipherText := gcm.Seal(nil, iv, plainText, nil)
	log.Println("cipherText=", cipherText)
	return cipherText, nil
}

func (service *UrlServiceImpl) FindLongUrl(encodedShortUrl *string) (*url_model.FindLongUrlResp, errorlib.AppError) {
	cipherText, err := base64.URLEncoding.DecodeString(*encodedShortUrl)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	cipherStr := string(cipherText)
	plainText, appErr := service.decrypt(&cipherStr)
	if appErr != nil {
		return nil, appErr
	}
	var components url_model.ShortUrlComponents
	err = json.Unmarshal(plainText, &components)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

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

func (service *UrlServiceImpl) decrypt(cipherStr *string) ([]byte, errorlib.AppError) {
	cipherText := *cipherStr
	iv := cipherText[:aes.BlockSize] // Extract initialization vector (IV)
	cipherText = cipherText[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(encryptDecryptKey))
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}

	plainText, err := gcm.Open(nil, []byte(iv), []byte(cipherText), nil)
	if err != nil {
		return nil, errorlib.NewInternalServerError(err.Error())
	}
	log.Println("decrypted plain text=", plainText)
	return plainText, nil
}
