package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	logtest "github.com/Sirupsen/logrus/hooks/test"

	"github.com/kpurdon/go-api-example/internal/repos"
	"github.com/stretchr/testify/assert"
)

var l *log.Logger
var tlog *logtest.Hook

func init() {
	l, tlog = logtest.NewNullLogger()
}

type ReposTestClient struct {
	Repos []repos.Repo // TODO: should I not use type from internal package?
	Err   error
}

func (c ReposTestClient) Get(string) ([]repos.Repo, error) {
	return c.Repos, c.Err
}

func TestGetReposHandler(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description        string
		app                *App
		url                string
		expectedStatusCode int
		expectedBody       string
		expectedLog        string
	}{
		{
			description: "missing argument user",
			app: &App{
				Log:   l,
				repos: &ReposTestClient{},
			},
			url:                "/repos",
			expectedStatusCode: 400,
			expectedBody:       "MISSING_ARG_USER\n",
			expectedLog:        "missing argument user",
		}, {
			description: "error getting repos",
			app: &App{
				Log: l,
				repos: &ReposTestClient{
					Repos: []repos.Repo{},
					Err:   errors.New("fake test error"),
				},
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 500,
			expectedBody:       "INTERNAL_ERROR\n",
			expectedLog:        "fake test error",
		}, {
			description: "no repos found",
			app: &App{
				Log: l,
				repos: &ReposTestClient{
					Repos: []repos.Repo{},
					Err:   nil,
				},
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 200,
			expectedBody:       `[]`,
			expectedLog:        "",
		}, {
			description: "succesfull query",
			app: &App{
				Log: l,
				repos: &ReposTestClient{
					Repos: []repos.Repo{
						repos.Repo{
							Name:        "test",
							Description: "a test",
						},
					},
					Err: nil,
				},
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 200,
			expectedBody:       `[{"name":"test","description":"a test"}]`,
			expectedLog:        "",
		},
		// TODO: not all cases are tested, but this is enough of a sample
	}

	for _, tc := range tests {

		req, err := http.NewRequest("GET", tc.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		tc.app.GetReposHandler(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
		assert.Equal(tc.expectedBody, w.Body.String(), tc.description)

		if tc.expectedLog != "" {
			assert.Equal(tc.expectedLog, tlog.LastEntry().Message)
		}
	}
}
