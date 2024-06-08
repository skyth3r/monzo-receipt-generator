package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := NewClient()

	if client.clientID == "" || client.clientSecret == "" {
		return fmt.Errorf("the Client ID and Client secret were not found in env vars")
	}

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "monzo-access-token",
	})
	if err != nil {
		return err
	}

	i, err := ring.Get("tokens")
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			if err := oauth(client); err != nil {
				return fmt.Errorf("failed to authenticate: %w", err)
			}
			// join the tokens and save them to the keychain
			tokens := client.accessToken + "::" + client.refreshToken
			if err := ring.Set(keyring.Item{
				Key:         "tokens",
				Data:        []byte(tokens),
				Label:       "Monzo Access Token",
				Description: "Access and refresh tokens for the Monzo API",
			}); err != nil {
				return fmt.Errorf("failed to set tokens in keychain: %w", err)
			}
		} else {
			return err
		}
	} else {
		// split the tokens from i
		tokens := string(i.Data)
		tokenSlice := strings.Split(tokens, "::")
		if len(tokenSlice) != 2 {
			return fmt.Errorf("unexpected token format: %s", tokens)
		}
		client.accessToken = tokenSlice[0]
		client.refreshToken = tokenSlice[1]
	}

	payload, err := os.ReadFile("payload.json")
	if err != nil {
		return err
	}

	var receipts []Receipt
	err = json.Unmarshal(payload, &receipts)
	if err != nil {
		return err
	}

	err = createReceipts(client, receipts)
	if err != nil {
		return err
	}

	return nil
}
