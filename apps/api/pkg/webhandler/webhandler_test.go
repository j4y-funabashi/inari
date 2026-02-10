package webhandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j4y_funabashi/inari/apps/api/pkg/webhandler"
	"github.com/stretchr/testify/assert"
)

func TestCreateCollection(t *testing.T) {
	router := webhandler.NewWebHandler()

	req, err := http.NewRequest(http.MethodPost, "/api/timeline/months", nil)
	assert.NoError(t, err)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}
