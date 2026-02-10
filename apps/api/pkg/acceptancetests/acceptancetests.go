package acceptancetests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests/specs"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

const listCollectionsPath = "/api/timeline/months"

type CollectionsActorDriver struct {
	BaseURL string
}

func (d CollectionsActorDriver) ListCollections() ([]specs.MediaCollection, error) {
	collections := []specs.MediaCollection{}

	collectionsEndpoint, err := url.JoinPath(d.BaseURL, listCollectionsPath)
	if err != nil {
		return collections, fmt.Errorf("failed to create URL: %s", err.Error())
	}
	res, err := http.Get(collectionsEndpoint)
	if err != nil {
		return collections, fmt.Errorf("failed to GET url: %s :: %s", collectionsEndpoint, err.Error())
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return collections, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	err = json.Unmarshal(body, &collections)
	if err != nil {
		return collections, fmt.Errorf("failed to unmarshal json: %s :: %s", body, err.Error())
	}

	return collections, nil
}

type CreateCollectionRequest struct {
	Name string
}

func (d CollectionsActorDriver) CreateCustomCollection() (specs.MediaCollection, error) {
	collection := specs.MediaCollection{}

	collectionsEndpoint, err := url.JoinPath(d.BaseURL, listCollectionsPath)
	if err != nil {
		return collection, fmt.Errorf("failed to create URL: %s", err.Error())
	}

	request := CreateCollectionRequest{}
	requestJson, err := json.Marshal(&request)
	if err != nil {
		return collection, fmt.Errorf("failed to marshal json: %s :: %s", collectionsEndpoint, err.Error())
	}
	requestBody := bytes.NewBuffer(requestJson)

	res, err := http.Post(collectionsEndpoint, "application/json", requestBody)
	if err != nil {
		return collection, fmt.Errorf("failed to GET url: %s :: %s", collectionsEndpoint, err.Error())
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return collection, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	err = json.Unmarshal(body, &collection)
	if err != nil {
		return collection, fmt.Errorf("failed to unmarshal json: %s :: %s", body, err.Error())
	}

	return collection, nil
}

func StartInariWebContainer(t *testing.T) string {
	t.Helper()

	ctx := context.Background()

	logger := log.TestLogger(t)

	// start container
	ctr, err := testcontainers.Run(
		context.Background(),
		"",
		testcontainers.WithLogger(logger),
		testcontainers.WithWaitStrategy(
			wait.ForLog("inari web server running on port"),
		),
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{
			Context: "../../",
			Repo:    "inari-web-test",
			Tag:     "latest",
		}),
		testcontainers.WithExposedPorts("8080/tcp"),
	)

	t.Cleanup(func() {
		testcontainers.CleanupContainer(t, ctr)
	})

	if err != nil {
		t.Fatalf("failed to create container: %s", err.Error())
	}

	containerEndpoint, err := ctr.Endpoint(ctx, "http")
	if err != nil {
		t.Fatalf("failed to get container host: %s", err.Error())
	}

	return containerEndpoint
}
