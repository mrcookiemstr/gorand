package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewApiCallResult(t *testing.T) {
	newApiCallResult := NewApiCallResult([]int64{1, 2, 3}, nil)

	assert := assert.New(t)
	assert.NotNil(newApiCallResult)
	assert.True(len(newApiCallResult.IntArr) == 3)
	assert.Nil(newApiCallResult.Err)
}

func Test_FetchIntsResultArr(t *testing.T) {
	apiResultA := NewApiCallResult([]int64{1, 2, 3}, nil)
	apiResultB := NewApiCallResult([]int64{4, 5, 6}, nil)
	apiResultC := NewApiCallResult([]int64{7, 8, 9}, nil)
	apiResultArr := []ApiCallResult{apiResultA, apiResultB, apiResultC}

	fetchIntsResult := FetchIntsResultArr(apiResultArr)

	assert := assert.New(t)
	assert.NotNil(fetchIntsResult)
	assert.Equal(4, len(fetchIntsResult))

	//items from api
	const itemsCount = int(3)
	const itemsStdDev = float64(1)

	assert.Equal(itemsCount, len(fetchIntsResult[0].Data))
	assert.Equal(itemsStdDev, fetchIntsResult[0].StdDev)
	assert.Equal(itemsCount, len(fetchIntsResult[1].Data))
	assert.Equal(itemsStdDev, fetchIntsResult[1].StdDev)
	assert.Equal(itemsCount, len(fetchIntsResult[2].Data))
	assert.Equal(itemsStdDev, fetchIntsResult[2].StdDev)

	const blockCount = int(9)
	const blockStdDev = float64(2.7386127875258306)

	//block summary
	assert.Equal(blockCount, len(fetchIntsResult[3].Data))
	assert.Equal(blockStdDev, fetchIntsResult[3].StdDev)
}

func Test_IntToFloatArr(t *testing.T) {
	intArr := []int64{1, 2, 3}
	result := IntToFloatArr(intArr)

	floatArr := []float64{float64(1), float64(2), float64(3)}

	assert := assert.New(t)
	assert.NotNil(result)
	assert.Equal(result, floatArr)
}
