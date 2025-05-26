package sale

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

func TestCreateSaleWithNonExistentUser(t *testing.T) {
	mockHandler := http.NewServeMux()

	mockHandler.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mockServer := httptest.NewServer(mockHandler)
	defer mockServer.Close()

	logger, _ := zap.NewDevelopment()
	storage := NewLocalStorage()
	service := NewService(storage, logger, mockServer.URL)

	sale, err := service.Create("non-existent-user", 150.0)

	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
	if sale != nil {
		t.Errorf("expected no sale to be created, got %+v", sale)
	}
}
