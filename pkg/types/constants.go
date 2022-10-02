package types

const (
	// Environment Variable Prefix
	EnvArgoCDPrefix = "ARGOCD_ENV"

	// Environment Variable Constants
	EnvAvpType             = "ATP_TYPE"
	EnvAvpRoleID           = "ATP_ROLE_ID"
	EnvAvpSecretID         = "ATP_SECRET_ID"
	EnvAvpAuthType         = "ATP_AUTH_TYPE"
	EnvAvpGithubToken      = "ATP_GITHUB_TOKEN"
	EnvAvpK8sRole          = "ATP_K8S_ROLE"
	EnvAvpK8sMountPath     = "ATP_K8S_MOUNT_PATH"
	EnvAvpMountPath        = "ATP_MOUNT_PATH"
	EnvAvpK8sTokenPath     = "ATP_K8S_TOKEN_PATH"
	EnvAvpIBMAPIKey        = "ATP_IBM_API_KEY"
	EnvAvpIBMInstanceURL   = "ATP_IBM_INSTANCE_URL"
	EnvAvpKvVersion        = "ATP_KV_VERSION"
	EnvAvpPathPrefix       = "ATP_PATH_PREFIX"
	EnvAWSRegion           = "AWS_REGION"
	EnvVaultAddress        = "VAULT_ADDR"
	EnvYCLKeyID            = "ATP_YCL_KEY_ID"
	EnvYCLServiceAccountID = "ATP_YCL_SERVICE_ACCOUNT_ID"
	EnvYCLPrivateKey       = "ATP_YCL_PRIVATE_KEY"
	EnvAvpUsername         = "ATP_USERNAME"
	EnvAvpPassword         = "ATP_PASSWORD"
	EnvAvpTFS3Bucket       = "ATP_TF_S3_BUCKET"
	EnvAvpTFS3Endpoint     = "ATP_TF_S3_ENDPOINT"
	EnvAvpTFS3AccessKey    = "ATP_TF_S3_ACCESS_KEY"
	EnvAvpTFS3SecretKey    = "ATP_TF_S3_SECRET_KEY"
	EnvAvpTFS3UseSSL       = "ATP_TF_S3_USE_SSL"

	// Backend and Auth Constants
	VaultBackend              = "vault"
	IBMSecretsManagerbackend  = "ibmsecretsmanager"
	AWSSecretsManagerbackend  = "awssecretsmanager"
	GCPSecretManagerbackend   = "gcpsecretmanager"
	AzureKeyVaultbackend      = "azurekeyvault"
	Sopsbackend               = "sops"
	YandexCloudLockboxbackend = "yandexcloudlockbox"
	OnePasswordConnect        = "1passwordconnect"
	K8sAuth                   = "k8s"
	ApproleAuth               = "approle"
	GithubAuth                = "github"
	TokenAuth                 = "token"
	UserPass                  = "userpass"
	IAMAuth                   = "iam"
	AwsDefaultRegion          = "us-east-2"
	GCPCurrentSecretVersion   = "latest"
	IBMMaxRetries             = 3
	IBMRetryIntervalSeconds   = 20
	IBMMaxPerPage             = 200
	IBMIAMCredentialsType     = "iam_credentials"
	IBMImportedCertType       = "imported_cert"
	IBMPublicCertType         = "public_cert"
	TerraformStateBackend     = "terraform"

	// Supported annotations
	ATPPathAnnotation          = "avp.kubernetes.io/path"
	ATPIgnoreAnnotation        = "avp.kubernetes.io/ignore"
	ATPRemoveMissingAnnotation = "avp.kubernetes.io/remove-missing"
	ATPSecretVersionAnnotation = "avp.kubernetes.io/secret-version"
	VaultKVVersionAnnotation   = "avp.kubernetes.io/kv-version"

	// Kube Constants
	ArgoCDNamespace = "argocd"
)
