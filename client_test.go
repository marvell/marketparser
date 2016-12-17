package marketparser

import (
	"fmt"
	"os"
	"testing"
)

func getApiKey() (string, error) {
	apiKey := os.Getenv("APIKEY")
	if len(apiKey) == 0 {
		return "", fmt.Errorf("empty APIKEY")
	}

	return apiKey, nil
}

func createClient(t *testing.T) *client {
	apiKey, err := getApiKey()
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewClient(apiKey)
	if c == nil || err != nil {
		t.Fatalf("failed to create Client: %s", err)
	}

	c.DebugMode()

	return c
}

func TestNewClient(t *testing.T) {
	client, err := NewClient("")
	if client != nil || err == nil {
		t.Errorf("we shouldn't created Client if user don't pass an API-key")
	}

	createClient(t)
}

func TestSimpleRequest(t *testing.T) {
	body, err := createClient(t).get("/campaigns/123456789/price.json", 1)
	if len(body) > 0 || err == nil {
		t.Fatalf("should be a failed request, but something went wrong:\nerror=%s\nbody=%q", err, body)
	}

	body, err = createClient(t).get("/campaigns.json", 1)
	if err != nil {
		t.Fatalf("got error while doing request: %s", err)
	}

	if len(body) == 0 {
		t.Fatalf("got empty body")
	}
}
