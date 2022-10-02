package backends

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/KazanExpress/argocd-terraform-plugin/pkg/utils"
	"github.com/minio/minio-go/v7"
)

// MinioClient is an interface to work with S3. It's needed to mock S3 communication during tests
type MinioClient interface {
	GetObject(ctx context.Context, bucket, path string, opt minio.GetObjectOptions) (io.Reader, error)
}

// TerraformState is a struct for working with a Terraform State backend
type TerraformState struct {
	client MinioClient
	bucket string
}

type TFOutput struct {
	Value interface{}
}

type TFState struct {
	Outputs map[string]*TFOutput `json:"outputs"`
}

type minioClientWrapper struct {
	mcl *minio.Client
}

func (mcw *minioClientWrapper) GetObject(ctx context.Context, bucket, path string, opt minio.GetObjectOptions) (io.Reader, error) {
	return mcw.mcl.GetObject(ctx, bucket, path, opt)
}

// WrapMinioClient wraps official minio.Client to match needed interface
func WrapMinioClient(c *minio.Client) MinioClient {
	return &minioClientWrapper{mcl: c}
}

// NewS3Backend initializes a new Terraform S3 State backend
func NewS3Backend(client MinioClient, bucket string) *TerraformState {
	return &TerraformState{
		client: client,
		bucket: bucket,
	}
}

// Login does nothing as a "login" is handled on the instantiation of the minio client
func (ycl *TerraformState) Login() error {
	return nil
}

// GetSecrets gets secrets from terraform state backend and returns the formatted data
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

	data, err := ioutil.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	var state TFState
	// state.Init()
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&state)
	if err != nil {
		utils.VerboseToStdErr("Terraform S3 State failed parsing json: %s: %w", string(data), err)

		return nil, fmt.Errorf("failed to decode state from json: %w", err)
	}

	results := make(map[string]interface{})
	for key, output := range state.Outputs {
		results[key] = output.Value
	}

	return results, nil
}

// GetIndividualSecret will get the specific secret (placeholder) from the terraform state backend
func (ycl *TerraformState) GetIndividualSecret(path, key, version string, _ map[string]string) (interface{}, error) {
	secrets, err := ycl.GetSecrets(path, version, nil)
	if err != nil {
		return nil, err
	}

	secret, found := secrets[key]
	if !found {
		utils.VerboseToStdErr("Terraform S3 State existing secrets: %v", secrets)

		return nil, fmt.Errorf("path: %s, key: %s, version: %s not found", path, key, version)
	}

	return secret, nil
}
