package tests

import (
	"env/src/environment"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Env struct {
	PackageName string `env:"PACKAGE_NAME"`
	LogLevel    string `env:"LOG_LEVEL"`
	Iterations  int    `env:"ITERATIONS"`
	BaseURL     string `env:"BASE_URL"`
}

type EnvOptional struct {
	PackageName string `env:"PACKAGE_NAME"`
	BaseURL     string `env:"BASE_URL"`
}

type EnvInvalidType struct {
	InvalidType any `env:"INVALID_TYPE"`
}

func TestAddValue(t *testing.T) {
	var env Env
	err := environment.LoadFile(".env", &env, environment.Config{})
	assert.Nil(t, err)
	assert.Empty(t, env.BaseURL)
	assert.Equal(t, "env", env.PackageName, env)
	assert.Equal(t, "debug", env.LogLevel, env)
	assert.Equal(t, 10, env.Iterations)
}

func TestMissingEnvInFile(t *testing.T) {
	var env Env
	err := environment.LoadFile(".env", &env, environment.Config{
		Force: true,
	})
	assert.Error(t, err)
	assert.Equal(t, "missing value for BaseURL", err.Error())
}

func TestOptionalEnvStruct(t *testing.T) {
	var (
		env         Env
		envOptional EnvOptional
		err         error
	)
	err = environment.LoadFile(".env", &env, environment.Config{})
	assert.Nil(t, err)
	err = environment.LoadFile(".env", &envOptional, environment.Config{})
	assert.Nil(t, err)
}

func TestInvalidTypeEnvStruct(t *testing.T) {
	var (
		env EnvInvalidType
		err error
	)
	err = environment.LoadFile(".env", &env, environment.Config{})
	assert.Error(t, err)
	assert.Equal(t, "env: type \"interface\" not supported", err.Error())
}

func TestInvalidPath(t *testing.T) {
	var (
		env EnvInvalidType
		err error
	)
	err = environment.LoadFile(".env-invalid-path", &env, environment.Config{})
	assert.Error(t, err)
	assert.Equal(t, "open .env-invalid-path: no such file or directory", err.Error())
}
