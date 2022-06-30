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
