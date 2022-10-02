package config_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func TestNewConfigNoType(t *testing.T) {
	viper := viper.New()
	_, err := config.New(viper, &config.Options{})
	expectedError := "Must provide a supported Vault Type, received "

	if err.Error() != expectedError {
		t.Errorf("expected error %s to be thrown, got %s", expectedError, err)
	}
}

func TestNewConfigNoAuthType(t *testing.T) {
	os.Setenv("ATP_BACKEND", "vault")
	viper := viper.New()
	_, err := config.New(viper, &config.Options{})
	expectedError := "Must provide a supported Authentication Type, received "

	if err.Error() != expectedError {
		t.Errorf("expected error %s to be thrown, got %s", expectedError, err)
	}
	os.Unsetenv("ATP_BACKEND")
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

func TestNewConfigAwsRegionWarning(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
		expectedLog  string
	}{
		{ // this test issues a warning for missing AWS_REGION env var
			map[string]interface{}{
				"ATP_BACKEND":           "awssecretsmanager",
				"AWS_ACCESS_KEY_ID":     "id",
				"AWS_SECRET_ACCESS_KEY": "key",
			},
			"*backends.AWSSecretsManager",
			"warning: AWS_REGION env var not set, using AWS region us-east-2\n",
		},
		{ // no warning is issued
			map[string]interface{}{
				"ATP_BACKEND":           "awssecretsmanager",
				"AWS_REGION":            "us-west-1",
				"AWS_ACCESS_KEY_ID":     "id",
				"AWS_SECRET_ACCESS_KEY": "key",
			},
			"*backends.AWSSecretsManager",
			"",
		},
	}

	for _, tc := range testCases {
		for k, v := range tc.environment {
			os.Setenv(k, v.(string))
		}
		viper.Set("verboseOutput", true)

		v := viper.New()
		output := captureOutput(func() {
			config, err := config.New(v, &config.Options{})
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			xType := fmt.Sprintf("%T", config.Backend)
			if xType != tc.expectedType {
				t.Errorf("expected: %s, got: %s.", tc.expectedType, xType)
			}
		})

		if !strings.Contains(output, tc.expectedLog) {
			t.Errorf("Unexpected warning issued. Expected: %s, actual: %s", tc.expectedLog, output)
		}

		for k := range tc.environment {
			os.Unsetenv(k)
		}
	}
}

func TestNewConfigMissingParameter(t *testing.T) {
	testCases := []struct {
		environment  map[string]interface{}
		expectedType string
	}{
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "github",
				"ATP_GH_TOKEN":  "token",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "token",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "approle",
				"ATP_ROLEID":    "role_id",
				"ATP_SECRET_ID": "secret_id",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "k8s",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "userpass",
				"ATP_USERNAME":  "username",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "userpass",
				"ATP_PASSWORD":  "password",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":   "vault",
				"ATP_AUTH_TYPE": "userpass",
			},
			"*backends.Vault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":     "ibmsecretsmanager",
				"ATP_IAM_API_KEY": "token",
			},
			"*backends.IBMSecretsManager",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND": "ibmsecretsmanager",
				"VAULT_ADDR":  "http://vault",
			},
			"*backends.IBMSecretsManager",
		},
		{ //  WebIdentityEmptyRoleARNErr will occur if 'AWS_WEB_IDENTITY_TOKEN_FILE' was set but 'AWS_ROLE_ARN' was not set.
			map[string]interface{}{
				"ATP_BACKEND":                 "awssecretsmanager",
				"AWS_REGION":                  "us-west-1",
				"AWS_WEB_IDENTITY_TOKEN_FILE": "/var/run/secrets/eks.amazonaws.com/serviceaccount/token",
			},
			"*backends.AWSSecretsManager",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":     "azurekeyvault",
				"AZURE_TENANT_ID": "test",
				"AZURE_CLIENT_ID": "test",
			},
			"*backends.AzureKeyVault",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":                "yandexcloudlockbox",
				"ATP_YCL_KEY_ID":             "test",
				"ATP_YCL_SERVICE_ACCOUNT_ID": "test",
			},
			"*backends.YandexCloudLockbox",
		},
		{
			map[string]interface{}{
				"ATP_BACKEND":      "1passwordconnect",
				"OP_CONNECT_TOKEN": "token",
			},
			"*backends.OnePasswordConnect",
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

func TestExternalConfig(t *testing.T) {
	os.Setenv("ATP_BACKEND", "vault")
	viper := viper.New()
	viper.SetDefault("VAULT_ADDR", "http://my-vault:8200/")
	config.New(viper, &config.Options{})
	if os.Getenv("VAULT_ADDR") != "http://my-vault:8200/" {
		t.Errorf("expected VAULT_ADDR env to be set from external config, was instead: %s", os.Getenv("VAULT_ADDR"))
	}
	os.Unsetenv("ATP_BACKEND")
	os.Unsetenv("VAULT_ADDR")
}

const avpConfig = `ATP_BACKEND: awssecretsmanager
AWS_ACCESS_KEY_ID: AKIAIOSFODNN7EXAMPLE
AWS_SECRET_ACCESS_KEY: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
AWS_REGION: us-west-2`

var expectedEnvVars = map[string]string{
	"ATP_BACKEND":           "", // shouldn't be an env var
	"AWS_ACCESS_KEY_ID":     "AKIAIOSFODNN7EXAMPLE",
	"AWS_SECRET_ACCESS_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	"AWS_REGION":            "us-west-2",
}

func TestExternalConfigAWS(t *testing.T) {
	// Test setting AWS_* env variables from external ATP config, note setting
	// env vars is necessary to pass ATP config entries to the AWS golang SDK
	tmpFile, err := ioutil.TempFile("", "avpConfig.*.yaml")
	if err != nil {
		t.Errorf("Cannot create temporary file %s", err)
	}

	defer os.Remove(tmpFile.Name()) // clean up the file afterwards

	if _, err = tmpFile.WriteString(avpConfig); err != nil {
		t.Errorf("Failed to write to temporary file %s", err)
	}

	viper := viper.New()
	if _, err = config.New(viper, &config.Options{ConfigPath: tmpFile.Name()}); err != nil {
		t.Errorf("config.New returned error: %s", err)
	}

	if viper.GetString("ATP_BACKEND") != "awssecretsmanager" {
		t.Errorf("expected ATP_BACKEND to be set from external config, was instead: %s", viper.GetString("ATP_BACKEND"))
	}

	for envVar, expected := range expectedEnvVars {
		if actual := os.Getenv(envVar); actual != expected {
			t.Errorf("expected %s env to be %s, was instead: %s", envVar, expected, actual)
		}
	}

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_REGION")
}

func TestExternalConfigSOPS(t *testing.T) {
	const avpSOPSConfig = `ATP_BACKEND: sops
SOPS_AGE_KEY_FILE: age`

	expectedSOPSEnvVars := map[string]string{
		"ATP_BACKEND":       "", // shouldn't be an env var
		"SOPS_AGE_KEY_FILE": "age",
	}

	// Test setting SOPS_* env variables from external ATP config, note setting
	// env vars is necessary to pass ATP config entries to SOPS
	tmpFile, err := ioutil.TempFile("", "avpSOPSConfig.*.yaml")
	if err != nil {
		t.Errorf("Cannot create temporary file %s", err)
	}

	defer os.Remove(tmpFile.Name()) // clean up the file afterwards

	if _, err = tmpFile.WriteString(avpSOPSConfig); err != nil {
		t.Errorf("Failed to write to temporary file %s", err)
	}

	viper := viper.New()
	if _, err = config.New(viper, &config.Options{ConfigPath: tmpFile.Name()}); err != nil {
		t.Errorf("config.New returned error: %s", err)
	}

	if viper.GetString("ATP_BACKEND") != "sops" {
		t.Errorf("expected ATP_BACKEND to be set from external config, was instead: %s", viper.GetString("ATP_BACKEND"))
	}

	for envVar, expected := range expectedSOPSEnvVars {
		if actual := os.Getenv(envVar); actual != expected {
			t.Errorf("expected %s env to be %s, was instead: %s", envVar, expected, actual)
		}
	}

	os.Unsetenv("SOPS_AGE_KEY_FILE")
}
