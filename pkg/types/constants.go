package types

const (
	// Environment Variable Prefix
	EnvArgoCDPrefix = "ARGOCD_ENV"

	// Environment Variable Constants
	EnvAtpBackend     = "ATP_BACKEND"
	EnvAtpS3Bucket    = "ATP_S3_BUCKET"
	EnvAtpS3Endpoint  = "ATP_S3_ENDPOINT"
	EnvAtpS3AccessKey = "ATP_S3_ACCESS_KEY"
	EnvAtpS3SecretKey = "ATP_S3_SECRET_KEY"
	EnvAtpS3UseSSL    = "ATP_S3_USE_SSL"

	// Backend and Auth Constants
	S3Backend = "s3"

	// Supported annotations
	ATPPathAnnotation          = "atp.kubernetes.io/path"
	ATPIgnoreAnnotation        = "atp.kubernetes.io/ignore"
	ATPRemoveMissingAnnotation = "atp.kubernetes.io/remove-missing"

	// Kube Constants
	ArgoCDNamespace = "argocd"
)
