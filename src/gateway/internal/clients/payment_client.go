package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gazizov-ai/lab2-rsoi/src/gateway/internal/model"
)

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *PaymentClient) CreatePayment(username string, price int) (model.Payment, error) {
	body := map[string]interface{}{
		"username": username,
		"price":    price,
	}
	data, _ := json.Marshal(body)

	url := fmt.Sprintf("%s/internal/payments", c.baseURL)

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return model.Payment{}, fmt.Errorf("create payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Payment{}, fmt.Errorf("payment status %d", resp.StatusCode)
	}

	var p model.Payment
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return model.Payment{}, fmt.Errorf("decode payment: %w", err)
	}

	return p, nil
}

func (c *PaymentClient) GetPayment(uid string) (model.Payment, error) {
	url := fmt.Sprintf("%s/internal/payments/%s", c.baseURL, uid)

	resp, err := c.client.Get(url)
	if err != nil {
		return model.Payment{}, fmt.Errorf("get payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Payment{}, fmt.Errorf("payment status %d", resp.StatusCode)
	}

	var p model.Payment
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return model.Payment{}, fmt.Errorf("decode payment: %w", err)
	}

	return p, nil
}

func (c *PaymentClient) CancelPayment(uid string) error {
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/internal/payments/%s", c.baseURL, uid),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("cancel payment: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cancel payment status %d", resp.StatusCode)
	}

	return nil
}
