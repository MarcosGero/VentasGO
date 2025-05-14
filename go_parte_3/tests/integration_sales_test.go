package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"parte3/api"
	"parte3/internal/sale"
	"testing"
)

// Happy path: POST → PATCH → GET
func TestSalesFlow(t *testing.T) {
	app := gin.Default()
	api.InitRoutes(app)

	// 1) crear usuario
	resUser := post(app, "/users", `{"name":"Ana","address":"Rioja","nickname":"ani"}`, http.StatusCreated, t)

	var u struct{ ID string }
	require.NoError(t, json.Unmarshal(resUser.Body.Bytes(), &u))

	// 2) crear venta
	saleJSON := `{"user_id":"` + u.ID + `","amount":123}`
	resSale := post(app, "/sales", saleJSON, http.StatusCreated, t)

	var s sale.Sale
	require.NoError(t, json.Unmarshal(resSale.Body.Bytes(), &s))
	require.Equal(t, sale.StatusApproved, s.Status) // o pending / rejected → random

	// 3) si quedó pending lo pasamos a approved
	if s.Status == sale.StatusPending {
		body := `{"status":"approved"}`
		resUpdate := patch(app, "/sales/"+s.ID, body, http.StatusOK, t)
		require.NoError(t, json.Unmarshal(resUpdate.Body.Bytes(), &s))
		require.Equal(t, sale.StatusApproved, s.Status)
	}

	// 4) buscar
	searchRes := get(app, "/sales?user_id="+u.ID, http.StatusOK, t)

	var out sale.SearchResult
	require.NoError(t, json.Unmarshal(searchRes.Body.Bytes(), &out))
	require.Equal(t, 1, out.Metadata.Quantity)
}

func post(e *gin.Engine, url, body string, code int, t *testing.T) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	return do(e, r, code, t)
}
func patch(e *gin.Engine, url, body string, code int, t *testing.T) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(http.MethodPatch, url, bytes.NewBufferString(body))
	return do(e, r, code, t)
}
func get(e *gin.Engine, url string, code int, t *testing.T) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	return do(e, r, code, t)
}
func do(e *gin.Engine, r *http.Request, code int, t *testing.T) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)
	require.Equal(t, code, w.Code)
	return w
}
