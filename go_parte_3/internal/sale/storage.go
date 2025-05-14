package sale

import "errors"

var (
	ErrNotFound     = errors.New("sale not found")
	ErrEmptyID      = errors.New("empty sale ID")
	ErrBadStatus    = errors.New("invalid status")
	ErrBadAmount    = errors.New("amount must be greater than 0")
	ErrBadTrans     = errors.New("invalid status transition")
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
	Set(*Sale) error
	Read(id string) (*Sale, error)
	Delete(id string) error
	All() []*Sale                 // para búsqueda
	ByUser(userID string) []*Sale // atajo útil
}

type LocalStorage struct{ m map[string]*Sale }

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{m: map[string]*Sale{}}
}

func (l *LocalStorage) Set(s *Sale) error {
	if s.ID == "" {
		return ErrEmptyID
	}
	l.m[s.ID] = s
	return nil
}

func (l *LocalStorage) Read(id string) (*Sale, error) {
	s, ok := l.m[id]
	if !ok {
		return nil, ErrNotFound
	}
	return s, nil
}

func (l *LocalStorage) Delete(id string) error {
	if _, ok := l.m[id]; !ok {
		return ErrNotFound
	}
	delete(l.m, id)
	return nil
}

func (l *LocalStorage) All() []*Sale {
	out := make([]*Sale, 0, len(l.m))
	for _, s := range l.m {
		out = append(out, s)
	}
	return out
}

func (l *LocalStorage) ByUser(uid string) []*Sale {
	out := []*Sale{}
	for _, s := range l.m {
		if s.UserID == uid {
			out = append(out, s)
		}
	}
	return out
}
