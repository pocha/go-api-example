package repos

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Client provides an interface to the repos package
type Client interface {
	Get(string) ([]Repo, error)
}

// ReposClient provides an implmentation of the Client interface
type ReposClient struct{}

// Repo representes GitHub repository
type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Get calls the GitHub API to and returns a Repo object for a given user
func (c ReposClient) Get(user string) ([]Repo, error) {
	var r []Repo

	reposURL := fmt.Sprintf("https://api.github.com/users/%s/repos", user)

	res, err := http.Get(reposURL)
	if err != nil {
		return nil, fmt.Errorf("github api: unknown error, %s", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("github api: error decoding response %s", err)
		}
	case 404:
		return nil, fmt.Errorf("github api: no results found")
	default:
		return nil, fmt.Errorf("github api: unknown error")
	}

	return r, nil
}
