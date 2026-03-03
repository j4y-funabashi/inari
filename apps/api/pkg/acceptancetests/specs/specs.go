package specs

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type CollectionsDriver interface {
	ListCollections() ([]MediaCollection, error)
	CreateCustomCollection(collectionTitle, collectionType string) (MediaCollection, error)
}

type MediaCollection struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
}

func CollectionsSpec(t *testing.T, driver CollectionsDriver) {
	collectionTitle := uuid.New().String()
	collectionType := "hashtag"

	newCollection, err := driver.CreateCustomCollection(collectionTitle, collectionType)
	assert.NoError(t, err)
	assert.Equal(t, collectionTitle, newCollection.Title)
	assert.Equal(t, collectionType, newCollection.Type)

	collections, err := driver.ListCollections()
	assert.NoError(t, err)
	assert.NotEmpty(t, collections)
	t.Log(collections)
	assert.Contains(t, collections, newCollection)
}
