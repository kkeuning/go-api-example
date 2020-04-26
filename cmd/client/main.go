package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func listUsers(apiKey *string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8090/api/v1/users", nil)
	if err != nil {
		os.Exit(1)
	}
	req.Close = true
	if apiKey != nil {
		req.Header.Add("Authorization", *apiKey)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		os.Exit(1)
	}
	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(respBody))
}

func showUserByID(id int, apiKey *string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8090/api/v1/users/%d", id), nil)
	if err != nil {
		os.Exit(1)
	}
	req.Close = true
	if apiKey != nil {
		req.Header.Add("Authorization", *apiKey)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		os.Exit(1)
	}
	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(string(respBody))
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
