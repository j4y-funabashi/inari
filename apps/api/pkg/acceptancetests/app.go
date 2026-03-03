package acceptancetests

import (
	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests/specs"
	"github.com/j4y_funabashi/inari/apps/api/pkg/app"
)

type AppWithFakesDriver struct {
	App app.App
}

func (d AppWithFakesDriver) ListCollections() ([]specs.MediaCollection, error) {
	collections := []specs.MediaCollection{}

	collectionType := app.CollectionTypeHashTag
	c, err := d.App.ListCollections(collectionType)
	if err != nil {
		return collections, err
	}
	for _, col := range c {
		collections = append(collections, appCollectionToMediaCollection(col))
	}

	return collections, nil
}

func (d AppWithFakesDriver) CreateCustomCollection(collectionTitle, collectionType string) (specs.MediaCollection, error) {
	c, err := d.App.CreateCollection(collectionTitle, app.CollectionType(collectionType))
	return appCollectionToMediaCollection(c), err
}

func appCollectionToMediaCollection(c app.Collection) specs.MediaCollection {
	return specs.MediaCollection{
		ID:    c.ID,
		Title: c.Title,
		Type:  string(c.Type),
	}
}
