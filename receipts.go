package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func createReceipt(c *MonzoClient, payload []byte) error {
	path := "transaction-receipts"
	requestURL := fmt.Sprintf("%s/%s", c.endpoints["APIURL"], path)

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Content-Type", "application/json")

	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return fmt.Errorf("response body is empty")
	}

	return nil
}

func createReceipts(c *MonzoClient, r []Receipt) error {
	for i, receipt := range r {
		rb, err := json.Marshal(receipt)
		if err != nil {
			return err
		}
		err = createReceipt(c, rb)
		if err != nil {
			return err
		}
		fmt.Printf("Uploaded %d/%d receipts\n", i+1, len(r))
	}
	return nil
}
