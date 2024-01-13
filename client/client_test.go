package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestFetchUsers(t *testing.T) {

	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			t.Errorf("Expected GET request to /admin/v1/users, got %s request to %s", req.Method, req.URL.Path)
			http.Error(rw, "Invalid request", http.StatusBadRequest)
			return
		}
		data, err := os.ReadFile("./test_data/users_response.json")
		assert.NoError(t, err)
		rw.Write(data)
	}))

	defer server.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	spec := Spec{
		ClientId:     "abc",
		ClientSecret: "alksjdnfkljsadfkjnasd",
		ApiHost:      server.URL,
	}

	ct, err := New(context.TODO(), logger, &spec)
	assert.NoError(t, err)
	usersFromDUO, err := ct.GetUsers(context.TODO(), 1, 10)
	assert.NoError(t, err)
	assert.Greater(t, len(usersFromDUO), 0)
}
