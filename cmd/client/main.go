package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kkeuning/go-api-example/pkg/models"
)

// Request is a standard request that will close the body and client
func Request(method string, url string, requestBody []byte, options ...func(*http.Request)) []byte {
	client := &http.Client{}
	var req *http.Request
	if requestBody != nil {
		r, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
		if err != nil {
			os.Exit(1)
		}
		req = r
	} else {
		r, err := http.NewRequest(method, url, nil)
		if err != nil {
			os.Exit(1)
		}
		req = r
	}
	for i := range options {
		options[i](req)
	}
	req.Close = true
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		os.Exit(1)
	}
	// Read Response Body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}
	return responseBody
}

func showUserByID(id int, apiKey *string) {
	url := fmt.Sprintf("http://localhost:8090/api/v1/users/%d", id)
	opts := func(r *http.Request) {
		if apiKey != nil {
			r.Header.Add("Authorization", *apiKey)
		}
	}
	respBody := Request(http.MethodGet, url, nil, opts)
	var user models.User
	if err := json.Unmarshal(respBody, &user); err != nil {
		os.Exit(1)
	}
	fmt.Printf("%s, %s %s\n", user.LastName, user.FirstName, user.MiddleInitial)
}

func listUsers(apiKey *string) {
	url := fmt.Sprintf("http://localhost:8090/api/v1/users")
	opts := func(r *http.Request) {
		if apiKey != nil {
			r.Header.Add("Authorization", *apiKey)
		}
	}
	respBody := Request(http.MethodGet, url, nil, opts)
	var out bytes.Buffer
	if err := json.Indent(&out, respBody, "", "    "); err != nil {
		os.Exit(1)
	}
	out.WriteTo(os.Stdout)
	fmt.Println()
	var users []models.User
	if err := json.Unmarshal(respBody, &users); err != nil {
		os.Exit(1)
	}
	for _, y := range users {
		fmt.Printf("%s, %s %s\n", y.LastName, y.FirstName, y.MiddleInitial)
	}
}

func main() {
	var apiKey string
	listCmd := flag.NewFlagSet("list-users", flag.ExitOnError)
	listCmd.StringVar(&apiKey, "apikey", "", "api key")
	getCmd := flag.NewFlagSet("get-user", flag.ExitOnError)
	userID := getCmd.Int("id", 0, "user id")
	getCmd.StringVar(&apiKey, "apikey", "", "api key")

	if len(os.Args) < 2 {
		fmt.Println("expected 'list-users' or 'get-user' subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "list-users":
		listCmd.Parse(os.Args[2:])
		listUsers(&apiKey)
	case "get-user":
		getCmd.Parse(os.Args[2:])
		if userID != nil {
			showUserByID(*userID, &apiKey)
		}
	default:
		fmt.Println("expected 'list-users' or 'get-user' subcommands")
		os.Exit(1)
	}
}
