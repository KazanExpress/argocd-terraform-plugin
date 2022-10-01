package backends

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/argoproj-labs/argocd-vault-plugin/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/minio/minio-go/v7"
)

type MinioClient interface {
	GetObject(ctx context.Context, bucket, path string, opt minio.GetObjectOptions) (io.Reader, error)
}

// TerraformState is a struct for working with a Terraform State backend
type TerraformState struct {
	client MinioClient
	bucket string
}

type minioClientWrapper struct {
	mcl *minio.Client
}

func (mcw *minioClientWrapper) GetObject(ctx context.Context, bucket, path string, opt minio.GetObjectOptions) (io.Reader, error) {
	return mcw.mcl.GetObject(ctx, bucket, path, opt)
}

func WrapMinioClient(c *minio.Client) MinioClient {
	return &minioClientWrapper{mcl: c}
}

// NewTerraformStateBackend initializes a new Terraform S3 State backend
func NewTerraformStateBackend(client MinioClient, bucket string) *TerraformState {
	return &TerraformState{
		client: client,
		bucket: bucket,
	}
}

// Login does nothing as a "login" is handled on the instantiation of the minio client
func (ycl *TerraformState) Login() error {
	return nil
}

// GetSecrets gets secrets from lockbox and returns the formatted data
func (ycl *TerraformState) GetSecrets(path string, version string, _ map[string]string) (map[string]interface{}, error) {

	var options = minio.GetObjectOptions{}

	if version != "" {
		options.VersionID = version
	}

	utils.VerboseToStdErr("Terraform S3 State getting object %s at version %s", path, version)
	obj, err := ycl.client.GetObject(context.Background(), ycl.bucket, path, options)
	if err != nil {
		return nil, fmt.Errorf("mc get object: %w", err)
	}

	utils.VerboseToStdErr("Terraform S3 State got object %v", obj)

	var state terraform.State
	err = json.NewDecoder(obj).Decode(&state)
	if err != nil {
		return nil, err
	}

	results := make(map[string]interface{})
	for key, output := range state.RootModule().Outputs {
		results[key] = output.Value
	}

	return results, nil
}

// GetIndividualSecret will get the specific secret (placeholder) from the lockbox backend
func (ycl *TerraformState) GetIndividualSecret(path, key, version string, _ map[string]string) (interface{}, error) {
	secrets, err := ycl.GetSecrets(path, version, nil)
	if err != nil {
		return nil, err
	}

	secret, found := secrets[key]
	if !found {
		return nil, fmt.Errorf("secretID: %s, key: %s, version: %s not found", path, key, version)
	}

	return secret, nil
}
