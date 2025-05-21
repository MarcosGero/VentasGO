package sale

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(userID string, amount float64) (*Sale, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) Create(userID string, amount float64) (*Sale, error) {
	if !userExists(userID) {
		return nil, ErrUserNotFound
	}

	if amount <= 0 {
		return nil, ErrBadAmount
	}

	now := time.Now()
	sale := &Sale{
		ID:        uuid.New().String(),
		UserID:    userID,
		Amount:    amount,
		Status:    randomStatus(),
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}

	if err := s.storage.Set(sale); err != nil {
		return nil, err
	}

	return sale, nil
}

func userExists(id string) bool {
	url := fmt.Sprintf("http://localhost:8080/users/%s", id)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false
	}

	return true
}

func randomStatus() Status {
	statuses := []Status{StatusPending, StatusApproved, StatusRejected}
	rand.Seed(time.Now().UnixNano())
	return statuses[rand.Intn(len(statuses))]
}
