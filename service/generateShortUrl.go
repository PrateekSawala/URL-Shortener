package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"url-shortener/domain"
	"url-shortener/domain/logging"

	"github.com/google/uuid"
)

var (
	_, b, _, _    = runtime.Caller(0)
	URLRecordPath = filepath.Dir(path.Join(path.Dir(b))) + domain.URLRecordFolderPath
)

func GenerateShortURL(host string, longUrl string) (*domain.UrlInfo, error) {
	log := logging.Log("GenerateShortURL")
	log.Tracef("Start")
	defer log.Tracef("End")

	response := &domain.UrlInfo{}
	URLRecords, err := GetURLRecords()
	if err != nil {
		log.Debugf("GetURLRecords error: %s", err)
		return response, err
	}

	urlRecord := CheckUrl(URLRecords, longUrl)
	// Check urlRecord
	if urlRecord.ID == "" {
		urlRecord, err = SaveUrl(URLRecords, longUrl)
		if err != nil {
			log.Debugf("SaveUrl error: %s", err)
			return response, err
		}
	}
	response.Url = fmt.Sprintf("http://%s/%s/%s", host, domain.ShortUrl, urlRecord.ID)
	return response, nil
}

func SaveUrl(URLRecords []domain.UrlRecord, longUrl string) (domain.UrlRecord, error) {
	response := domain.UrlRecord{}

	// Generate UUid
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)

	// Add new record
	record := domain.UrlRecord{ID: uuid, LongUrl: longUrl}
	URLRecords = append(URLRecords, record)

	// Marshal the records
	result, err := json.MarshalIndent(URLRecords, "", "  ")
	if err != nil {
		return response, err
	}

	err = ioutil.WriteFile(URLRecordPath, result, 0644)
	if err != nil {
		return response, err
	}

	response = record
	return response, nil
}

func CheckUrl(URLRecords []domain.UrlRecord, url string) domain.UrlRecord {
	response := domain.UrlRecord{}
	// Loop over all URLRecords
	for _, urlRecord := range URLRecords {
		if urlRecord.LongUrl == url || urlRecord.ID == url {
			response = urlRecord
			break
		}
	}
	return response
}

func GetURLRecords() ([]domain.UrlRecord, error) {
	response := []domain.UrlRecord{}
	// Check if the URLRecord path exist or not
	_, err := os.Stat(URLRecordPath)
	if err == nil {
		// Reading file from the path
		record, err := ioutil.ReadFile(URLRecordPath)
		if err != nil {
			return response, err
		}
		// Find URL Records
		err = json.Unmarshal([]byte(record), &response)
		if err != nil {
			return response, err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// Create new record file
		_, err := os.Create(URLRecordPath)
		if err != nil {
			return response, err
		}
	} else {
		return response, err
	}
	return response, nil
}