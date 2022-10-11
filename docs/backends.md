### S3 State

Currently supported all providers supported by minio-go client.

##### Auth

```
ATP_BACKEND: s3
ATP_S3_ENDPOINT: s3.endpoint.com
ATP_S3_BUCKET: tf-states
ATP_S3_ACCESS_KEY: your-access-key
ATP_S3_SECRET_KEY: your-secret-key
ATP_USE_SSL: true
```

##### Examples

###### Path Annotation

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: tf-example
  annotations:
    atp.kubernetes.io/path: "path/to/state/example.tfstate"
type: Opaque
data:
  username: <terraform:username>
  password: <terraform:password>
```

###### Inline Path

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: tf-example
type: Opaque
data:
  username: <terraform:path/to/state/example.tfstate#username>
  password: <terraform:path/to/state/example.tfstate#password>
```
###### Sub key

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: test-secret
  annotations:
    atp.kubernetes.io/path: "path/to/state/example.tfstate"
type: Opaque
stringData:
  password: <terraform:parent | jsonPath {.child}>
