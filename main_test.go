package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ProvideHostName(t *testing.T) {
	result := ProvideHostName()
	assert := assert.New(t)

	assert.NotNil(result)
	assert.Equal(result, "localhost")
}

func Test_ProvideHttpRequestPort_WhenSetEnv(t *testing.T) {
	t.Setenv(envHttpPortName, "1234")

	httpPort := ProvideHttpRequestPort()
	assert := assert.New(t)

	assert.NotNil(httpPort)
	assert.Equal(httpPort, "1234")
}

func Test_ProvideHttpRequestPort_WhenSetEnvEmpty(t *testing.T) {
	httpPort := ProvideHttpRequestPort()
	assert := assert.New(t)

	assert.NotNil(httpPort)
	assert.Equal(httpPort, "8080")
}

func Test_ProvideRandomOrgApiKey_WhenSetEnv(t *testing.T) {
	t.Setenv(defaultRandomOrgKey, "abcd-1234")

	apiKey := ProvideRandomOrgApiKey()
	assert := assert.New(t)
	assert.NotNil(apiKey)
	assert.Equal(apiKey, "abcd-1234")
}

func Test_ProvideRandomOrgApiKey_WhenSetEnvEmpty(t *testing.T) {
	panic := false

	defer func() {
		if r := recover(); r != nil {
			panic = true
		}
	}()

	ProvideRandomOrgApiKey()

	assert := assert.New(t)
	assert.True(panic)
}
