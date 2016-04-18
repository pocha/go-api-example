package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/kpurdon/go-api-example/internal/repos"
)

// App defines the application container
type App struct {
	Log   *log.Logger
	repos repos.Client
}

// GetReposHandler returns a list of (public) repositories for a given GitHub user
func (a *App) GetReposHandler(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if user == "" {
		a.Log.Error("missing argument user")
		http.Error(w, "MISSING_ARG_USER", 400)
		return
	}

	repos, err := a.repos.Get(user)
	if err != nil {
		a.Log.Error(err)
		http.Error(w, "INTERNAL_ERROR", 500)
		return
	}

	b, err := json.Marshal(repos)
	if err != nil {
		a.Log.Errorf("error marshaling repos: %s", err)
		http.Error(w, "INTERNAL_ERROR", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	app := &App{
		Log:   log.New(),
		repos: repos.ReposClient{},
	}

	http.HandleFunc("/repos", app.GetReposHandler)

	app.Log.Println("listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
