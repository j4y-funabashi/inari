package acceptancetests_test

import (
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests"
	"github.com/j4y_funabashi/inari/apps/api/pkg/acceptancetests/specs"
	appconfig "github.com/j4y_funabashi/inari/apps/api/pkg/app_config"
)

func TestMediaCollections(t *testing.T) {
	testCases := []struct {
		desc        string
		setup       func(t *testing.T) specs.CollectionsDriver
		skipIfShort bool
	}{
		{
			desc:        "http",
			setup:       setupHttpDriver,
			skipIfShort: true,
		},
		{
			desc:        "app with fakes",
			setup:       setupAppWithFakesDriver,
			skipIfShort: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if testing.Short() && tC.skipIfShort {
				t.Skip("acceptance test")
			}

			t.Run("create custom collection and it appears in list", func(t *testing.T) {
				collectionsDriver := tC.setup(t)
				specs.CollectionsSpec(t, collectionsDriver)
			})
		})
	}
}

func setupAppWithFakesDriver(t *testing.T) specs.CollectionsDriver {
	inariApp := appconfig.New()
	return acceptancetests.AppWithFakesDriver{App: inariApp}
}

func setupHttpDriver(t *testing.T) specs.CollectionsDriver {
	containerEndpoint := acceptancetests.StartInariWebContainer(t)
	collectionsDriver := acceptancetests.CollectionsHttpDriver{BaseURL: containerEndpoint}
	return collectionsDriver
}
