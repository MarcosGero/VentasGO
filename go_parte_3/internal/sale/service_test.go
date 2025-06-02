package sale

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestUpdateStatus_InvalidTransition(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	// Usar almacenamiento local real en memoria
	storage := NewLocalStorage()
	// Crear ventas directamente con estado aprobado y rechazado
	approvedSale := &Sale{
		ID:        "sale-approved",
		UserID:    "user-123",
		Amount:    100.0,
		Status:    StatusApproved,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
	rejectedSale := &Sale{
		ID:        "sale-rejected",
		UserID:    "user-456",
		Amount:    200.0,
		Status:    StatusRejected,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
	//Cargo las ventas en el storage
	_ = storage.Set(approvedSale)
	_ = storage.Set(rejectedSale)
	service := NewService(storage, logger, "http://mock-api.local")
	// Intentar rechazar una venta ya aprobada
	_, err := service.UpdateStatus("sale-approved", StatusRejected)
	if err != ErrBadTrans {
		t.Errorf("expected ErrBadTrans when rejecting approved sale, got %v", err)
	}
	// Intentar aprobar una venta ya rechazada
	_, err = service.UpdateStatus("sale-rejected", StatusApproved)
	if err != ErrBadTrans {
		t.Errorf("expected ErrBadTrans when approving rejected sale, got %v", err)
	}
	//Intentar cambiar una venta aprobada a pendiente
	_, err = service.UpdateStatus("sale-approved", StatusPending)
	if err != ErrBadStatus {
		t.Errorf("expected ErrBadStatus when setting approved sale to pending, got %v", err)
	}
	//Intentar cambiar una venta rechazada a pendiente
	_, err = service.UpdateStatus("sale-rejected", StatusPending)
	if err != ErrBadStatus {
		t.Errorf("expected ErrBadStatus when setting rejected sale to pending, got %v", err)
	}
}
