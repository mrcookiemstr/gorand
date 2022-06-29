package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandleHttpRequest_GivenRootPath(t *testing.T) {
	url := "/some/endpoint"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()

	HandleHttpRequest(recorder, request)

	result := recorder.Result()
	defer result.Body.Close()

	status := result.StatusCode
	body, _ := ioutil.ReadAll(result.Body)

	assert := assert.New(t)
	assert.Equal(http.StatusNotFound, status)
	assert.NotNil(body)
}

func Test_HandleHttpRequest_GivenEndpoint_WhenParamsMissing(t *testing.T) {
	url := RandomEndpoint
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()

	HandleHttpRequest(recorder, request)

	result := recorder.Result()
	defer result.Body.Close()

	status := result.StatusCode
	body, _ := ioutil.ReadAll(result.Body)

	assert := assert.New(t)
	assert.Equal(http.StatusBadRequest, status)
	assert.NotNil(body)
}

func Test_HandleHttpRequest_GivenEndpoint_WhenEmptyParams(t *testing.T) {
	url := RandomEndpoint + "?requests=&length="
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()

	HandleHttpRequest(recorder, request)

	result := recorder.Result()
	defer result.Body.Close()

	status := result.StatusCode
	body, _ := ioutil.ReadAll(result.Body)

	assert := assert.New(t)
	assert.Equal(http.StatusBadRequest, status)
	assert.NotNil(body)
}

func Test_HandleHttpRequest_GivenEndpoint_WhenIntParams(t *testing.T) {
	url := RandomEndpoint + "?requests=2&length=3"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()

	HandleHttpRequest(recorder, request)

	result := recorder.Result()
	defer result.Body.Close()

	status := result.StatusCode
	body, _ := ioutil.ReadAll(result.Body)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, status)
	assert.NotNil(body)
}

func Test_GetParamAsInt_WhenCorrectParam(t *testing.T) {
	url := "/?abc=123"
	request := httptest.NewRequest(http.MethodGet, url, nil)

	param, err := GetParamAsInt(request, "abc")

	assert := assert.New(t)
	assert.NotNil(param)
	assert.Equal(123, param)
	assert.Nil(err)
}

func Test_GetParamAsInt_WhenEmptyParam(t *testing.T) {
	url := "/?abc="
	request := httptest.NewRequest(http.MethodGet, url, nil)

	param, err := GetParamAsInt(request, "abc")

	assert := assert.New(t)
	assert.NotNil(param)
	assert.Equal(-1, param)
	assert.NotNil(err)
	assert.Equal(ErrParamEmptyString, err)
}

func Test_GetParamAsInt_WhenParamNegativeInteger(t *testing.T) {
	url := "/?abc=-1"
	request := httptest.NewRequest(http.MethodGet, url, nil)

	param, err := GetParamAsInt(request, "abc")

	assert := assert.New(t)
	assert.NotNil(param)
	assert.Equal(-1, param)
	assert.NotNil(err)
	assert.Equal(ErrParamNegativeInt, err)
}

func Test_GetParamAsInt_WhenParamIsString(t *testing.T) {
	url := "/?abc=GhAc"
	request := httptest.NewRequest(http.MethodGet, url, nil)

	param, err := GetParamAsInt(request, "abc")

	assert := assert.New(t)
	assert.NotNil(param)
	assert.Equal(-1, param)
	assert.NotNil(err)
	assert.Equal(ErrParamParseError, err)
}
