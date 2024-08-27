package replacer

import (
	"testing"

	bitwarden "github.com/bitwarden/sdk-go"
	"github.com/stretchr/testify/assert"
)

type TestLookup struct{}

func (t TestLookup) Get(secretID string) (*bitwarden.SecretResponse, error) {
	switch secretID {
	case "ce398fa2-5665-11ef-8916-97605d6da25b":
		projectID := "ddb13dae-5665-11ef-8583-f73233caa8df"
		return &bitwarden.SecretResponse{
			CreationDate:   "2022-11-17T15:55:18.005669100Z",
			ID:             secretID,
			Key:            "Human Readable Key",
			Note:           "",
			OrganizationID: "d4105690-5665-11ef-a058-c713a9374bb0",
			ProjectID:      &projectID,
			RevisionDate:   "2022-11-17T15:55:18.005669100Z",
			Value:          "my_secret_password",
		}, nil
	case "98b6c8ee-5666-11ef-ac37-8742ac5fc78f":
		projectID := "ddb13dae-5665-11ef-8583-f73233caa8df"
		return &bitwarden.SecretResponse{
			CreationDate:   "2019-05-11T15:55:18.005669100Z",
			ID:             secretID,
			Key:            "Other Key",
			Note:           "",
			OrganizationID: "d4105690-5665-11ef-a058-c713a9374bb0",
			ProjectID:      &projectID,
			RevisionDate:   "2019-05-11T15:55:18.005669100Z",
			Value:          "my_other_secret",
		}, nil
	default:
		return nil, nil
	}
}

var (
	template = `apiVersion: v1
kind: Pod
metadata:
  annotations:
    k8s.v1.cni.cncf.io/networks: kube-system/iot
  name: ubuntu
  labels:
    app: ubuntu
spec:
  containers:
  - image: ubuntu
    command:
      - "sleep"
      - "604800"
    env:
      - name: PASSWORD
        value: <bw:ce398fa2-5665-11ef-8916-97605d6da25b>
	  - name: SECRET_ENV1
	    value: <bw:98b6c8ee-5666-11ef-ac37-8742ac5fc78f>
    imagePullPolicy: IfNotPresent
    name: ubuntu
  restartPolicy: Always
`

	expected = `apiVersion: v1
kind: Pod
metadata:
  annotations:
    k8s.v1.cni.cncf.io/networks: kube-system/iot
  name: ubuntu
  labels:
    app: ubuntu
spec:
  containers:
  - image: ubuntu
    command:
      - "sleep"
      - "604800"
    env:
      - name: PASSWORD
        value: my_secret_password
	  - name: SECRET_ENV1
	    value: my_other_secret
    imagePullPolicy: IfNotPresent
    name: ubuntu
  restartPolicy: Always
`
)

func TestReplacer(t *testing.T) {
	r := Replacer{
		client: TestLookup{},
	}

	result, err := r.Replace(template)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
