// Package instacart implements the client for the instacart web API.
package instacart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

// Client is the HTTP client for the Instacart orders API.
type Client struct {
	SessionToken string
	httpClient   *http.Client
}

// NewClient constructs a Client.
func NewClient(sessionToken string) *Client {
	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: jar,
	}
	return &Client{
		SessionToken: sessionToken,
		httpClient:   httpClient,
	}
}

func (c *Client) FetchPage(page int) (*OrdersResponse, error) {
	url := "https://www.instacart.com/v3/orders?page=" + strconv.Itoa(page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Host", "www.instacart.com")
	req.Header.Set("User-Agent", "Instacart Orders To CSV Client")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Identifier", "web")

	sessionCookie := &http.Cookie{
		Name:  "_instacart_session_id",
		Value: c.SessionToken,
	}
	req.AddCookie(sessionCookie)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request status not OK, statusCode: %d, response: %+v", resp.StatusCode, resp)
	}

	defer resp.Body.Close()

	var ordersResp *OrdersResponse

	if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
		return nil, fmt.Errorf("failed to decode the response: %w", err)
	}

	return ordersResp, nil
}
