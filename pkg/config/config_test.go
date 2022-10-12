package config_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/KazanExpress/argocd-terraform-plugin/pkg/config"
	"github.com/spf13/viper"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
	}{
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ENDPOINT":   "endpoint.com",
				"ATP_S3_BUCKET":     "bucket",
				"ATP_S3_ACCESS_KEY": "key",
				"ATP_S3_SECRET_KEY": "key",
			},
			"*backends.S3Backend",
		},
	}
	for _, tc := range testCases {
		for k, v := range tc.environment {
			os.Setenv(k, v.(string))
		}
		viper := viper.New()
		config, err := config.New(viper, &config.Options{})
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		xType := fmt.Sprintf("%T", config.Backend)
		if xType != tc.expectedType {
			t.Errorf("expected: %s, got: %s.", tc.expectedType, xType)
		}
		for k := range tc.environment {
			os.Unsetenv(k)
		}
	}
}

// Helper function that captures log output from a function call into a string
// Adapted from https://stackoverflow.com/a/26806093/170154
func captureOutput(f func()) string {
	var buf bytes.Buffer
	flags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0) // don't include any date or time in the logging messages
	f()
	log.SetOutput(os.Stderr)
	log.SetFlags(flags)
	return buf.String()
}

func TestNewConfigMissingParameter(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
	}{
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ACCESS_KEY": "github",
				"ATP_S3_BUCKET":     "token",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":     "s3",
				"ATP_S3_ENDPOINT": "token.com",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ACCESS_KEY": "userpass",
				"ATP_S3_SECRET_KEY": "username",
				"ATP_S3_ENDPOINT":   "endpoint.com",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ACCESS_KEY": "userpass",
				"ATP_S3_SECRET_KEY": "username",
				"ATP_S3_BUCKET":     "token",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ACCESS_KEY": "userpass",
				"ATP_S3_SECRET_KEY": "username",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_ACCESS_KEY": "userpass",
				"ATP_S3_BUCKET":     "password",
			},
			"*backends.S3Backend",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":       "s3",
				"ATP_S3_SECRET_KEY": "userpass",
			},
			"*backends.S3Backend",
		},
	}
	for _, tc := range testCases {
		for k, v := range tc.environment {
			os.Setenv(k, v.(string))
		}
		viper := viper.New()
		_, err := config.New(viper, &config.Options{})
		if err == nil {
			t.Fatalf("%s should not instantiate", tc.expectedType)
		}
		for k := range tc.environment {
			os.Unsetenv(k)
		}
	}
}
