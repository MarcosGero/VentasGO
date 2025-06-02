package tests

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"parte3/api"
	"parte3/internal/sale"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestCreateSale_UserNotFound
// Test unitario: intenta crear una venta con un usuario inexistente y espera ErrUserNotFound.
func TestCreateSale_UserNotFound(t *testing.T) {
	// Mock user API que siempre devuelve 404 Not Found
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mockServer := httptest.NewServer(mockMux)
	defer mockServer.Close()

	st := sale.NewLocalStorage()
	svc := sale.NewService(st, zap.NewNop(), mockServer.URL)
	_, err := svc.Create("usuario-no-existe", 100)
	require.ErrorIs(t, err, sale.ErrUserNotFound)
}

// TestSalesFlow_HappyPath
// Test de integración: levanta Gin y prueba un flujo completo POST → PATCH → GET en el happy path.
func TestSalesFlow_HappyPath(t *testing.T) {
	// 1) Mock del servicio de usuarios en localhost:8080
	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"` + r.URL.Path[len("/users/"):] + `"}`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// Forzar el mock a escuchar en el puerto 8080
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	require.NoError(t, err)
	mockServer := &httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: mockMux},
	}
	mockServer.Start()
	defer mockServer.Close()

	// 2) Crear el router Gin con rutas de usuarios y ventas
	gin.SetMode(gin.TestMode)
	app := gin.New()
	api.InitRoutes(app)

	// 3) POST /sales
	recorder := httptest.NewRecorder()
	saleBody := `{"user_id":"u123","amount":150}`
	req, _ := http.NewRequest(http.MethodPost, "/sales", bytes.NewBufferString(saleBody))
	req.Header.Set("Content-Type", "application/json")
	app.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)

	// 4) Leer la venta creada
	var s sale.Sale
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &s))
	require.Contains(t,
		[]sale.Status{sale.StatusPending, sale.StatusApproved, sale.StatusRejected},
		s.Status,
	)

	// 5) Si quedó pending, actualizar a approved
	if s.Status == sale.StatusPending {
		patchRec := httptest.NewRecorder()
		patchBody := `{"status":"approved"}`
		patchReq, _ := http.NewRequest(http.MethodPatch, "/sales/"+s.ID, bytes.NewBufferString(patchBody))
		patchReq.Header.Set("Content-Type", "application/json")
		app.ServeHTTP(patchRec, patchReq)
		require.Equal(t, http.StatusOK, patchRec.Code)
		require.NoError(t, json.Unmarshal(patchRec.Body.Bytes(), &s))
		require.Equal(t, sale.StatusApproved, s.Status)
	}

	// 6) GET /sales?user_id=u123
	getRec := httptest.NewRecorder()
	getReq, _ := http.NewRequest(http.MethodGet, "/sales?user_id=u123", nil)
	app.ServeHTTP(getRec, getReq)
	require.Equal(t, http.StatusOK, getRec.Code)

	var out sale.SearchResult
	require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &out))
	require.Equal(t, 1, out.Metadata.Quantity)
}
