package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gobuffalo/envy"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kkeuning/go-api-example/pkg/auth"
	"github.com/kkeuning/go-api-example/pkg/services"
	"github.com/kkeuning/go-api-example/pkg/services/users"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

//var db *sqlx.DB

func hello(e *services.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Hello, World!</h1>\n")
	})
}

func getUser(e *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))
		e.Log.Debug("*** Get User ***")
		id := mux.Vars(req)["id"]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an integer.")
			return
		}
		userID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Expected id as an integer.")
			return
		}
		payload := &users.ShowPayload{
			ID: userID,
		}
		result, err := usersSvc.Show(ctx, payload)
		if err != nil {
			switch err.(type) {
			case users.ErrorNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "Error")
			return
		}
		juser, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong.")
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(juser)
		return
	})
}

func listUsers(e *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Populate request context
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))
		e.Log.Debug("*** List Users ***")
		id := req.URL.Query().Get("id")
		var payload *users.ListPayload
		if id != "" {
			userID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Expected id as an integer.")
				return
			}
			payload = &users.ListPayload{
				ID: userID,
			}
		}
		userList, err := usersSvc.List(ctx, payload)
		if err != nil {
			switch err.(type) {
			case users.ErrorNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			fmt.Fprintf(w, "Error")
			return
		}
		jusers, err := json.Marshal(userList.Users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error")
			return
		}
		w.Write(jusers)
		return
	})
}

func main() {
	environment := envy.Get("ENVIRONMENT", "prod")
	if environment == "prod" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		// The TextFormatter is default, you don't actually have to do this.
		logger.SetFormatter(&logrus.TextFormatter{})
	}
	// db, err := sqlx.Open("sqlite3", ":memory:")
	// if err != nil {
	// 	logger.Fatal("Unable to open db connection.")
	// }
	// db.Ping() // We're not doing much else with the db yet
	// logger.Info("DB connected")

	env := &services.Env{
		Log: logger,
		//DB:  db,
	}

	// Read API key from command line flag if provided.
	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "API key")
	flag.Parse()

	usersService, err := users.NewUsersSvc(env.Log)
	if err != nil {
		logger.Fatal("Failed to create users service.")
	}

	r := mux.NewRouter()
	r.Handle("/", hello(env))
	r.Handle("/api/v1/users", listUsers(env, usersService))
	r.Handle("/api/v1/users/{id}", getUser(env, usersService))

	// The simple API key security is optional.
	// If a key is provided, we will protect all routes containing "/api/".
	if apiKey != "" {
		akm := auth.APIKeyMiddleware{Path: "/api/"}
		akm.InitializeKey(apiKey)
		r.Use(akm.Middleware)
	}

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	logger.Fatal(http.ListenAndServe(":8090", loggedRouter))
}
