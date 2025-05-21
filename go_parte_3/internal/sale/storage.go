package sale

import (
	"errors"
	"sync"
)

var (
	ErrNotFound     = errors.New("sale not found")
	ErrBadStatus    = errors.New("invalid status")
	ErrBadAmount    = errors.New("amount must be greater than 0")
	ErrBadTrans     = errors.New("invalid status transition")
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
	Set(*Sale) error
	Read(id string) (*Sale, error)
	Delete(id string) error
	All() []*Sale
	ByUser(userID string) []*Sale
}

type LocalStorage struct {
	mtx sync.Mutex
	m   map[string]*Sale
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: make(map[string]*Sale),
	}
}

func (l *LocalStorage) Set(s *Sale) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if s.ID == "" {
		return errors.New("sale ID cannot be empty")
	}

	l.m[s.ID] = s
	return nil
}

func (l *LocalStorage) Read(id string) (*Sale, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	s, ok := l.m[id]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

func (l *LocalStorage) Delete(id string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.m[id]; !ok {
		return ErrNotFound
	}
	delete(l.m, id)
	return nil
}

func (l *LocalStorage) All() []*Sale {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	out := make([]*Sale, 0, len(l.m))
	for _, s := range l.m {
		out = append(out, s)
	}
	return out
}

func (l *LocalStorage) ByUser(uid string) []*Sale {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	out := []*Sale{}
	for _, s := range l.m {
		if s.UserID == uid {
			out = append(out, s)
		}
	}
	return out
}
