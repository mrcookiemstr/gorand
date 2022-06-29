package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	RandomEndpoint = "/random/mean"
	RequestParam   = "requests"
	LengthParam    = "length"
)

var (
	ErrParamEmptyString = errors.New("empty string")
	ErrParamNegativeInt = errors.New("negative integer")
	ErrParamParseError  = errors.New("parsing error")
)

func HandleHttpRequest(responseWriter http.ResponseWriter, request *http.Request) {
	requstPath := request.URL.Path
	log.Println("url: ", requstPath)

	switch requstPath {
	case RandomEndpoint:
		HandleRandomEndpoint(responseWriter, request)
	default:
		http.NotFoundHandler().ServeHTTP(responseWriter, request)
	}
}

func HandleRandomEndpoint(responseWriter http.ResponseWriter, request *http.Request) {
	requests, requestsErr := GetParamAsInt(request, RequestParam)
	length, lengthErr := GetParamAsInt(request, LengthParam)

	if requestsErr != nil || lengthErr != nil {
		ResponseParamsError(responseWriter, requestsErr, lengthErr)
	} else {
		fetchIntResult, fetchIntErr := FetchInts(request.Context(), requests, length)

		if fetchIntErr != nil {
			ResponseFetchIntError(responseWriter, fetchIntErr)
		} else {
			ResponseFetchInt(responseWriter, fetchIntResult)
		}
	}
}

func ResponseParamsError(responseWriter http.ResponseWriter, requestsErr error, lengthErr error) {
	response := "Issue with params! \r\n"

	if requestsErr != nil {
		response = response + RequestParam + " - " + requestsErr.Error() + ".\r\n"
	}

	if lengthErr != nil {
		response = response + LengthParam + " - " + lengthErr.Error() + ".\r\n"
	}

	responseWriter.WriteHeader(400)
	responseWriter.Header().Set("Content-Type", "text/plain")
	io.WriteString(responseWriter, response)
}

func ResponseFetchIntError(responseWriter http.ResponseWriter, apiError error) {
	if apiError != nil && errors.Is(apiError, ErrFetchIntsTimeout) {
		responseWriter.WriteHeader(http.StatusGatewayTimeout)
		responseWriter.Header().Set("Content-Type", "text/plain")
		io.WriteString(responseWriter, "timeout while fetching data from random.org")
	} else if errors.Is(apiError, ErrFetchIntsApi) {
		responseWriter.WriteHeader(http.StatusBadGateway)
		responseWriter.Header().Set("Content-Type", "text/plain")
		io.WriteString(responseWriter, "random.org api issue")
	} else {
		responseWriter.WriteHeader(http.StatusServiceUnavailable)
		responseWriter.Header().Set("Content-Type", "text/plain")
		io.WriteString(responseWriter, "unknown server error")
	}
}

func ResponseFetchInt(responseWriter http.ResponseWriter, fetchIntResults []FetchIntsResult) {
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(fetchIntResults)
}

func GetParamAsInt(request *http.Request, paramName string) (int, error) {
	requestParam := request.URL.Query().Get(paramName)

	if requestParam == "" {
		return -1, ErrParamEmptyString
	}

	requestCount, atoiError := strconv.Atoi(requestParam)

	if atoiError != nil {
		return -1, ErrParamParseError
	}

	if requestCount <= 0 {
		return -1, ErrParamNegativeInt
	}

	return requestCount, nil
}
