package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/sgade/randomorg"
	"gonum.org/v1/gonum/stat"
)

var (
	ErrRandApiKey    = errors.New("wrong api key")
	ErrRandApiClient = errors.New("error while api clinet")
	ErrApiIntRange   = errors.New("int range out of bound")
	ErrApiNetwork    = errors.New("network error")

	ErrFetchIntsTimeout = errors.New("timeout")
	ErrFetchIntsApi     = errors.New("api error")
)

type (
	ApiCallResult struct {
		IntArr []int64
		Err    error
	}

	FetchIntsResult struct {
		StdDev float64 `json:"stddev"`
		Data   []int64 `json:"data"`
	}
)

func NewApiCallResult(intArr []int64, err error) ApiCallResult {
	return ApiCallResult{IntArr: intArr, Err: err}
}

func FetchIntsResultArr(apiResultArr []ApiCallResult) []FetchIntsResult {
	fetchIntsResultArr := []FetchIntsResult{}

	allIntArr := []int64{}
	fetchFloatArr := []float64{}

	for _, item := range apiResultArr {
		floatArr := IntToFloatArr(item.IntArr)
		stdDev := stat.StdDev(floatArr, nil)

		fetchIntsResult := FetchIntsResult{
			Data:   item.IntArr,
			StdDev: stdDev,
		}

		allIntArr = append(allIntArr, item.IntArr...)
		fetchFloatArr = append(fetchFloatArr, floatArr...)
		fetchIntsResultArr = append(fetchIntsResultArr, fetchIntsResult)
	}

	if len(allIntArr) > 0 {
		stdDev := stat.StdDev(fetchFloatArr, nil)

		fetchIntsResult := FetchIntsResult{
			StdDev: stdDev,
			Data:   allIntArr,
		}

		fetchIntsResultArr = append(fetchIntsResultArr, fetchIntsResult)
	}

	return fetchIntsResultArr
}

func IntToFloatArr(intArr []int64) []float64 {
	floatArr := []float64{}

	for _, item := range intArr {
		floatArr = append(floatArr, float64(item))
	}

	return floatArr
}

func IntsFromRandomOrg(apiKey string, noOfInts int) ([]int64, error) {
	if apiKey == "" {
		return nil, ErrRandApiKey
	}

	randApi := randomorg.NewRandom(apiKey)
	resultValues, err := randApi.GenerateIntegers(noOfInts, 0, 999999)

	if err != nil && errors.Is(err, randomorg.ErrParamRange) {
		return []int64{}, ErrApiIntRange
	}

	if err != nil && errors.Is(err, randomorg.ErrAPIKey) {
		return []int64{}, ErrRandApiKey
	}

	if err != nil {
		return []int64{}, ErrApiNetwork
	}

	return resultValues, nil
}

func ApiCallRoutine(ctx context.Context, waitGrpoup *sync.WaitGroup, noOfInts int, resultsChan chan<- ApiCallResult) {
	defer func() {
		waitGrpoup.Done()
	}()

	result, err := IntsFromRandomOrg("665740fa-1cad-4fa5-8bdf-9bbfb2a25b7e", noOfInts)

	select {
	case <-ctx.Done():
		log.Println("ApiCallRoutine timeout", ctx.Err())
		return
	case resultsChan <- NewApiCallResult(result, err):
	}
}

func FetchInts(ctx context.Context, concurentRequestNo int, noOfInts int) ([]FetchIntsResult, error) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(concurentRequestNo)

	resultChan := make(chan ApiCallResult)
	defer close(resultChan)

	for n := 0; n < concurentRequestNo; n++ {
		go ApiCallRoutine(ctx, &waitGroup, noOfInts, resultChan)
	}

	apiResultArr := []ApiCallResult{}
	apiTimeoutError := false
	apiCallError := false

	for n := 0; n < concurentRequestNo; n++ {
		select {
		case <-ctx.Done():
			apiTimeoutError = true

		case apiResult := <-resultChan:
			if apiResult.Err != nil {
				apiCallError = true
			}

			apiResultArr = append(apiResultArr, apiResult)
		}
	}

	waitGroup.Wait()

	if apiTimeoutError {
		return nil, ErrFetchIntsTimeout
	}

	if apiCallError {
		return nil, ErrFetchIntsApi
	}

	return FetchIntsResultArr(apiResultArr), nil
}
