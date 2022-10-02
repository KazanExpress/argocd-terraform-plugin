package backends_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/argoproj-labs/argocd-vault-plugin/pkg/backends"
	"github.com/minio/minio-go/v7"
)

type mockMinioClient struct {
	objects map[string]map[string][]byte
}

func newMockMinioClient() *mockMinioClient {
	return &mockMinioClient{
		objects: map[string]map[string][]byte{},
	}
}

func (m *mockMinioClient) GetObject(_ context.Context, bucketName, path string, opt minio.GetObjectOptions) (io.Reader, error) {
	if bucket, ok := m.objects[bucketName]; ok {
		if obj, ok := bucket[path]; ok {
			return bytes.NewReader(obj), nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockMinioClient) setObject(bucket, path string, data []byte) {
	if _, ok := m.objects[bucket]; !ok {
		m.objects[bucket] = make(map[string][]byte)
	}
	m.objects[bucket][path] = data
}

func TestTerraformState(t *testing.T) {

	const (
		strVal = "str"
	)
	objVal := map[string]interface{}{
		"key":      "value",
		"otherKey": []string{"item", "item2"},
	}

	bucketName := "argocd-test"
	path := "/test/case/1/terraform.state"
	st := backends.TFState{
		Outputs: map[string]*backends.TFOutput{
			"test_string": {
				Value: strVal,
			},
			"test_obj": {
				Value: objVal,
			},
		}}
	stateJson, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	mock := newMockMinioClient()
	mock.setObject(bucketName, path, stateJson)

	backend := backends.NewS3Backend(mock, bucketName)

	t.Run("Terraform State GetSecrets()", func(t *testing.T) {

		secrets, err := backend.GetSecrets(path, "", nil)
		if err != nil {
			t.Fatal(err)
		}

		if realStrVal, ok := secrets["test_string"]; ok {
			if strVal != realStrVal {
				t.Fatalf("test_string secret expected to be %v but received %v", strVal, realStrVal)
			}
		} else {
			t.Fatal("failed to get test_string secret")
		}

		if realObjVal, ok := secrets["test_obj"]; ok {
			if reflect.DeepEqual(realObjVal, objVal) {
				t.Fatalf("test_obj secret expected to be %v but received %v", objVal, realObjVal)
			}
		} else {
			t.Fatal("failed to get test_obj secret")
		}
	})

	t.Run("Terraform State GetIndividualSecret()", func(t *testing.T) {
		realStrVal, err := backend.GetIndividualSecret(path, "test_string", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		if strVal != realStrVal {
			t.Fatalf("test_string secret expected to be %v but received %v", strVal, realStrVal)
		}

		realObjVal, err := backend.GetIndividualSecret(path, "test_obj", "", nil)
		if reflect.DeepEqual(realObjVal, objVal) {
			t.Fatalf("test_obj secret expected to be %v but received %v", objVal, realObjVal)
		}
	})
}
