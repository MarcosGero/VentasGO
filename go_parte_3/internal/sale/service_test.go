package sale_test

import (
	//"errors"
	"net/http"
	"parte3/internal/sale"
	"testing"

	//"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

type fakeHTTP struct{ code int }

func (f fakeHTTP) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: f.code, Body: http.NoBody}, nil
}

func TestCreate_UserNotFound(t *testing.T) {
	st := sale.NewLocalStorage()
	svc := sale.NewService(st, nil, "http://fake")
	// stub resty transport
	//svc.Client().GetClient().Transport = fakeHTTP{code: http.StatusNotFound}

	_, err := svc.Create("abc", 10)
	require.ErrorIs(t, err, sale.ErrUserNotFound)
}
