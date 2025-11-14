package clients

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/model"
)

type LoyaltyClient struct {
	baseURL string
	client  *http.Client
}

func NewLoyaltyClient(baseURL string) *LoyaltyClient {
	return &LoyaltyClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *LoyaltyClient) GetLoyalty(username string) (model.Loyalty, error) {
	url := fmt.Sprintf("%s/internal/loyalty/%s", c.baseURL, username)

	resp, err := c.client.Get(url)
	if err != nil {
		return model.Loyalty{}, fmt.Errorf("request loyalty: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Loyalty{}, fmt.Errorf("loyalty status %d", resp.StatusCode)
	}

	var lo model.Loyalty
	if err := json.NewDecoder(resp.Body).Decode(&lo); err != nil {
		return model.Loyalty{}, fmt.Errorf("decode loyalty: %w", err)
	}

	return lo, nil
}
