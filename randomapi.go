package main

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/sgade/randomorg"
	"gonum.org/v1/gonum/stat"
)

type (
	RandomOrgClinetFunc func(int) ([]int64, error)

	ApiCallResult struct {
		IntArr []int64
		Err    error
	}

	FetchIntsResult struct {
		StdDev float64 `json:"stddev"`
		Data   []int64 `json:"data"`
	}
)

const (
	RandomOrgKey = "665740fa-1cad-4fa5-8bdf-9bbfb2a25b7e"
)

var (
	RandomOrgClient = func(noOfInts int) ([]int64, error) {
		api := randomorg.NewRandom(RandomOrgKey)
		result, err := api.GenerateIntegers(noOfInts, 0, 999999)

		if err != nil && (errors.Is(err, randomorg.ErrParamRange) || errors.Is(err, randomorg.ErrAPIKey)) {
			panic(err)
		} else if err != nil {
			return []int64{}, ErrApiNetwork
		}

		return result, nil
	}

	ErrApiNetwork = errors.New("network error")

	ErrFetchIntsTimeout = errors.New("timeout")
	ErrFetchIntsApi     = errors.New("api error")
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

func ApiCallRoutine(ctx context.Context, waitGrpoup *sync.WaitGroup, client RandomOrgClinetFunc, noOfInts int, resultsChan chan<- ApiCallResult) {
	defer waitGrpoup.Done()

	result, err := client(noOfInts)

	select {
	case <-ctx.Done():
		log.Println("ApiCallRoutine timeout", ctx.Err())
		return
	case resultsChan <- NewApiCallResult(result, err):
	}
}

func FetchInts(ctx context.Context, client RandomOrgClinetFunc, concurentRequestNo int, noOfInts int) ([]FetchIntsResult, error) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(concurentRequestNo)

	resultChan := make(chan ApiCallResult)
	defer close(resultChan)

	for n := 0; n < concurentRequestNo; n++ {
		go ApiCallRoutine(ctx, &waitGroup, client, noOfInts, resultChan)
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
