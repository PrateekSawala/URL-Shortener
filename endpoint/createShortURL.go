package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"url-shortener/domain"
	"url-shortener/domain/logging"
	"url-shortener/service"
)

func CreateShortURLHandler(writer http.ResponseWriter, request *http.Request) {
	log := logging.Log("CreateShortURLHandler")
	log.Tracef("Start")
	defer log.Tracef("End")

	// Get Params
	routeQueryParameterUrl := request.URL.Query().Get("url")
	if routeQueryParameterUrl == "" {
		log.Debugf("url query parameter is empty")
		http.Error(writer, domain.EmptyInput, http.StatusBadRequest)
		return
	}

	// Check URL
	if !strings.Contains(routeQueryParameterUrl, "http") {
		routeQueryParameterUrl = fmt.Sprintf("http://%s", routeQueryParameterUrl)
	}

	err := IsURLAleadyShortened(request.Host, routeQueryParameterUrl)
	if err != nil {
		log.Debugf("Url %s is already a shortUrl, error %s", routeQueryParameterUrl, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = IsURLValid(routeQueryParameterUrl)
	if err != nil {
		log.Debugf("Url %s is invalidUrl, error %s", routeQueryParameterUrl, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	generateShortURLResp, err := service.GenerateShortURL(request.Host, routeQueryParameterUrl)
	if err != nil {
		log.Debugf("service.GenerateShortURL error: %s", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert object to json
	jsonResp, err := json.Marshal(generateShortURLResp)
	if err != nil {
		log.Debugf("json.Marshal error : %s", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response header
	writer.Header().Set("Content-Type", "application/json")

	// Returning json response
	_, err = writer.Write(jsonResp)
	if err != nil {
		log.Debugf("writer.Write error : %s", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	return
}
