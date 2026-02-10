package specs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CollectionsActor interface {
	ListCollections() ([]MediaCollection, error)
	CreateCustomCollection() (MediaCollection, error)
}

type MediaCollection struct {
	Name string `json:"name,omitempty"`
}

func CollectionsSpec(t *testing.T, actor CollectionsActor) {
	newCollection, err := actor.CreateCustomCollection()
	assert.NoError(t, err)
	assert.Equal(t, newCollection.Name, "")

	collections, err := actor.ListCollections()
	assert.NoError(t, err)
	assert.NotEmpty(t, collections)
	assert.Contains(t, collections, newCollection)
}
