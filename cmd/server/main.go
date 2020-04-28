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
	"github.com/google/gops/agent"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kkeuning/go-api-example/pkg/auth"
	"github.com/kkeuning/go-api-example/pkg/models"
	"github.com/kkeuning/go-api-example/pkg/services"
	"github.com/kkeuning/go-api-example/pkg/services/users"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func hello(e *services.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "<h1>Hello, World!</h1>\n")
	})
}

func getUser(e *services.Env, usersSvc users.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), services.AcceptTypeKey, req.Header.Get("Accept"))
		e.Log.Debug().Msg("*** Get User ***")
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
		e.Log.Debug().Msg("*** List Users ***")
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
		jusers, err := json.Marshal(&userList)
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
	logger := &log.Logger
	environment := envy.Get("ENVIRONMENT", "dev")
	if environment == "prod" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger.Info().Msg("*** Production Configuration ***")
		logger.Debug().Msg("*** Debug Logging Disabled ***") // Won't display
	} else {
		zl := logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		logger = &zl
		// Start gops agent
		if err := agent.Listen(agent.Options{}); err != nil {
			logger.Fatal().Err(err)
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Info().Msg("*** Non-production Configuration ***")
		logger.Debug().Msg("*** Debug Logging Enabled ***")
	}
	env := &services.Env{
		// DB:  db, // Shared database connection goes here
		Log: logger,
	}

	var apiKey string
	flag.StringVar(&apiKey, "apikey", "", "API key")
	flag.Parse()

	logger.Debug().Msg("*** Creating Users Service ***")
	uc := models.UserDB
	uc.Log = env.Log
	usersService, err := users.NewUsersSvc(env.Log, &uc)
	if err != nil {
		logger.Fatal().Msg("Failed to create users service.")
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
	logger.Fatal().Msg(http.ListenAndServe(":8090", loggedRouter).Error())
}
