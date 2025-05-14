package sale

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	storage Storage
	client  *resty.Client
	logger  *zap.Logger
	userURL string
	randSrc *rand.Rand
}

func NewService(st Storage, logger *zap.Logger, userAPIBaseURL string) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync()
	}
	return &Service{
		storage: st,
		client:  resty.New(),
		logger:  logger,
		userURL: userAPIBaseURL,
		randSrc: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Service) Create(userID string, amount float64) (*Sale, error) {
	if amount <= 0 {
		return nil, ErrBadAmount
	}
	if ok, err := s.userExists(userID); err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrUserNotFound
	}

	status := s.randomStatus()

	sale := &Sale{
		ID:        uuid.NewString(),
		UserID:    userID,
		Amount:    amount,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}
	if err := s.storage.Set(sale); err != nil {
		return nil, err
	}
	s.logger.Info("sale created", zap.Any("sale", sale))
	return sale, nil
}

func (s *Service) UpdateStatus(id string, newStatus Status) (*Sale, error) {
	if !IsValidStatus(newStatus) || newStatus == StatusPending {
		return nil, ErrBadStatus
	}
	sale, err := s.storage.Read(id)
	if err != nil {
		return nil, err
	}
	if sale.Status != StatusPending {
		return nil, ErrBadTrans
	}
	sale.Status = newStatus
	sale.UpdatedAt = time.Now()
	sale.Version++
	s.logger.Info("sale updated", zap.Any("sale", sale))
	return sale, nil
}

type SearchResult struct {
	Metadata struct {
		Quantity    int     `json:"quantity"`
		Approved    int     `json:"approved"`
		Rejected    int     `json:"rejected"`
		Pending     int     `json:"pending"`
		TotalAmount float64 `json:"total_amount"`
	} `json:"metadata"`
	Results []*Sale `json:"results"`
}

func (s *Service) Search(userID string, statusFilter *Status) (*SearchResult, error) {
	if statusFilter != nil && !IsValidStatus(*statusFilter) {
		return nil, ErrBadStatus
	}
	sales := s.storage.ByUser(userID)

	out := &SearchResult{}
	for _, sale := range sales {
		if statusFilter != nil && sale.Status != *statusFilter {
			continue
		}
		out.Results = append(out.Results, sale)
		out.Metadata.Quantity++
		out.Metadata.TotalAmount += sale.Amount
		switch sale.Status {
		case StatusApproved:
			out.Metadata.Approved++
		case StatusRejected:
			out.Metadata.Rejected++
		case StatusPending:
			out.Metadata.Pending++
		}
	}
	return out, nil
}

/*** helpers ***/

func (s *Service) randomStatus() Status {
	switch s.randSrc.Intn(3) { // 0,1,2
	case 0:
		return StatusPending
	case 1:
		return StatusApproved
	default:
		return StatusRejected
	}
}

func (s *Service) userExists(id string) (bool, error) {
	resp, err := s.client.R().SetDoNotParseResponse(true).Get(s.userURL + "/users/" + id)
	if err != nil {
		return false, err
	}
	defer resp.RawBody().Close()

	return resp.StatusCode() == http.StatusOK, nil
}
