package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kpurdon/go-api-example/internal/repos"
	"github.com/stretchr/testify/assert"
)

type ReposTestClient struct {
	Repos []repos.Repo
	Err   error
}

func (c ReposTestClient) Get(string) ([]repos.Repo, error) {
	return c.Repos, c.Err
}

func TestGetReposHandler(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description        string
		reposClient        *ReposTestClient
		url                string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description:        "missing argument user",
			reposClient:        &ReposTestClient{},
			url:                "/repos",
			expectedStatusCode: 400,
			expectedBody:       "MISSING_ARG_USER\n",
		}, {
			description: "error getting repos",
			reposClient: &ReposTestClient{
				Repos: []repos.Repo{},
				Err:   errors.New("fake test error"),
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 500,
			expectedBody:       "INTERNAL_ERROR\n",
		}, {
			description: "no repos found",
			reposClient: &ReposTestClient{
				Repos: []repos.Repo{},
				Err:   nil,
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 200,
			expectedBody:       `[]`,
		}, {
			description: "succesfull query",
			reposClient: &ReposTestClient{
				Repos: []repos.Repo{
					repos.Repo{Name: "test", Description: "a test"},
				},
				Err: nil,
			},
			url:                "/user?user=fakeuser",
			expectedStatusCode: 200,
			expectedBody:       `[{"name":"test","description":"a test"}]`,
		},
		// TODO not all cases are covered
	}

	for _, tc := range tests {
		app := &App{repos: tc.reposClient}

		req, err := http.NewRequest("GET", tc.url, nil)
		assert.NoError(err)

		w := httptest.NewRecorder()
		app.GetReposHandler(w, req)

		assert.Equal(tc.expectedStatusCode, w.Code, tc.description)
		assert.Equal(tc.expectedBody, w.Body.String(), tc.description)
	}
}
