package config

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/argoproj-labs/argocd-vault-plugin/pkg/backends"
	"github.com/argoproj-labs/argocd-vault-plugin/pkg/kube"
	"github.com/argoproj-labs/argocd-vault-plugin/pkg/types"
	"github.com/argoproj-labs/argocd-vault-plugin/pkg/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

// Options options that can be passed to a Config struct
type Options struct {
	SecretName string
	ConfigPath string
}

// Config is used to decide the backend and auth type
type Config struct {
	Backend types.Backend
}

var backendPrefixes []string = []string{
	"vault",
	"aws",
	"azure",
	"google",
	"sops",
	"op_connect",
}

// New returns a new Config struct
func New(v *viper.Viper, co *Options) (*Config, error) {
	// Set Defaults
	v.SetDefault(types.EnvAvpKvVersion, "2")

	// Read in config file or kubernetes secret and set as env vars
	err := readConfigOrSecret(co.SecretName, co.ConfigPath, v)
	if err != nil {
		return nil, err
	}

	// Instantiate Env
	utils.VerboseToStdErr("reading configuration from environment, overriding any previous settings")
	v.AutomaticEnv()

	utils.VerboseToStdErr("AVP configured with the following settings:\n")
	for k, viperValue := range v.AllSettings() {
		utils.VerboseToStdErr("%s: %s\n", k, viperValue)
	}

	authType := v.GetString(types.EnvAvpAuthType)

	var auth types.AuthType
	var backend types.Backend

	switch v.GetString(types.EnvAvpType) {
	case types.TerraformStateBackend:
		{
			if !v.IsSet(types.EnvAvpTFS3AccessKey) ||
				!v.IsSet(types.EnvAvpTFS3Bucket) ||
				!v.IsSet(types.EnvAvpTFS3Endpoint) ||
				!v.IsSet(types.EnvAvpTFS3SecretKey) {
				return nil, fmt.Errorf(
					"%s, %s, %s and %s are required for terraform state backend",
					types.EnvAvpTFS3AccessKey,
					types.EnvAvpTFS3Bucket,
					types.EnvAvpTFS3Endpoint,
					types.EnvAvpTFS3SecretKey,
				)
			}

			client, err := minio.New(v.GetString(types.EnvAvpTFS3Endpoint), &minio.Options{
				Creds: credentials.NewStaticV4(
					v.GetString(types.EnvAvpTFS3AccessKey),
					v.GetString(types.EnvAvpTFS3SecretKey),
					""),
				Secure: v.GetBool(types.EnvAvpTFS3UseSSL),
			})

			if err != nil {
				return nil, fmt.Errorf("failed to create minio client: %w", err)
			}

			backend = backends.NewS3Backend(backends.WrapMinioClient(client), v.GetString(types.EnvAvpTFS3Bucket))
		}
	default:
		return nil, fmt.Errorf("Must provide a supported Vault Type, received %s", v.GetString(types.EnvAvpType))
	}

	return &Config{
		Backend: backend,
	}, nil
}

func readConfigOrSecret(secretName, configPath string, v *viper.Viper) error {
	// If a secret name is passed, pull config from Kubernetes
	if secretName != "" {
		utils.VerboseToStdErr("reading configuration from secret %s", secretName)

		localClient, err := kube.NewClient()
		if err != nil {
			return err
		}
		yaml, err := localClient.ReadSecret(secretName)
		if err != nil {
			return err
		}
		v.SetConfigType("yaml")
		v.ReadConfig(bytes.NewBuffer(yaml))
	}

	// If a config file path is passed, read in that file and overwrite all other
	if configPath != "" {
		utils.VerboseToStdErr("reading configuration from config file %s, overriding any previous settings", configPath)

		v.SetConfigFile(configPath)
		err := v.ReadInConfig()
		if err != nil {
			return err
		}
	}

	// Check for ArgoCD 2.4 prefixed environment variables
	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, types.EnvArgoCDPrefix) {
			envVarPair := strings.SplitN(envVar, "=", 2)
			key := strings.TrimPrefix(envVarPair[0], types.EnvArgoCDPrefix+"_")
			val := envVarPair[1]
			v.Set(key, val)
		}
	}

	for k, viperValue := range v.AllSettings() {
		for _, prefix := range backendPrefixes {
			if strings.HasPrefix(k, prefix) {
				var value string
				switch viperValue.(type) {
				case bool:
					value = strconv.FormatBool(viperValue.(bool))
				default:
					value = viperValue.(string)
				}
				os.Setenv(strings.ToUpper(k), value)
				utils.VerboseToStdErr("Setting %s to %s for backend SDK", strings.ToUpper(k), value)
			}
		}
	}

	return nil
}
