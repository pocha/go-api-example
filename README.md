go-api-example
---

This repository provides an example http api written in Golang that I used to try and apply testing practices described in [this article](http://nathanleclaire.com/blog/2015/10/10/interfaces-and-composition-for-effective-unit-testing-in-golang/). I plan to do a future blog post on the work I've done here.

The API provides a single endpoint `/repos?user=<github_username>` that wraps up a call to GitHub's API to return a JSON array of repos for a given user.

### Getting Started

- `go get github.com/kpurdon/go-api-example`
- run `go-api-example`
- `curl localhost:8080/repos?user=<github_username>`

### Project Structure

```
├── README.md
├── internal
│   └── repos
│       ├── repos.go         - provides internal functions for calling the GitHub API
│       └── repos_test.go    - tests the internal function ("mocking" out the GitHub API calls)
├── main.go                  - provides the server and handler initalization
└── main_test.go             - tests the handler ("mocking" out the internal function calls)
```
