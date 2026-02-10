package acceptancetests_test

import (
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests"
	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests/specs"
)

func TestMediaCollections(t *testing.T) {
	if testing.Short() {
		t.Skip("acceptance test")
	}

	t.Run("create custom collection and it appears in list", func(t *testing.T) {
		containerEndpoint := acceptancetests.StartInariWebContainer(t)

		collectionsDriver := acceptancetests.CollectionsActorDriver{BaseURL: containerEndpoint}

		specs.CollectionsSpec(t, collectionsDriver)
	})
}
