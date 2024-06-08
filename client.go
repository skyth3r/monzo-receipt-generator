package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var client *MonzoClient

type MonzoClient struct {
	clientID     string
	clientSecret string
	accessToken  string
	refreshToken string
	callbackCode string
	endpoints    map[string]string
	httpClient   *http.Client
}

func init() {
	if client == nil {
		client = NewClient()
	}

	client = &MonzoClient{
		clientID:     os.Getenv("MONZO_CLIENT_ID"),
		clientSecret: os.Getenv("MONZO_CLIENT_SECRET"),
		accessToken:  "",
		refreshToken: "",
		callbackCode: "",
		endpoints: map[string]string{
			"AuthURL":  "https://auth.monzo.com",
			"TokenURL": "https://api.monzo.com/oauth2/token",
			"APIURL":   "https://api.monzo.com",
		},
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func NewClient() *MonzoClient {
	return client
}

func (c *MonzoClient) Do(req *http.Request) (*http.Response, error) {
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// If the response status code is 401, refresh the token and retry the request
	if rsp.StatusCode == http.StatusUnauthorized {
		if err := refreshToken(c); err != nil {
			return nil, err
		}

		// Update the Authorization header with the new access token.
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

		rsp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
	}

	return rsp, nil
}
