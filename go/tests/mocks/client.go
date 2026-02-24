/*
© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package mocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// BugsClient handles API communication for l8bugs tests
type BugsClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewBugsClient creates a new BugsClient with the given base URL and HTTP client
func NewBugsClient(baseURL string, httpClient *http.Client) *BugsClient {
	return &BugsClient{baseURL: baseURL, client: httpClient}
}

// BaseURL returns the client's base URL for constructing custom endpoints.
func (c *BugsClient) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying HTTP client for raw requests.
func (c *BugsClient) HTTPClient() *http.Client {
	return c.client
}

// L8QueryText builds a JSON-encoded L8Query with the text field
func L8QueryText(queryText string) string {
	q := map[string]interface{}{
		"text": queryText,
	}
	data, _ := json.Marshal(q)
	return string(data)
}

func (c *BugsClient) Authenticate(user, password string) error {
	authData := map[string]string{
		"user": user,
		"pass": password,
	}
	body, _ := json.Marshal(authData)

	resp, err := c.client.Post(c.baseURL+"/auth", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var authResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	token, ok := authResp["token"].(string)
	if !ok {
		return fmt.Errorf("token not found in auth response")
	}
	c.token = token
	return nil
}

func (c *BugsClient) Post(endpoint string, data interface{}) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return string(respBody), fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (c *BugsClient) Get(endpoint string, queryJSON string) (string, error) {
	fullURL := c.baseURL + endpoint + "?body=" + url.QueryEscape(queryJSON)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return string(respBody), fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (c *BugsClient) Put(endpoint string, data interface{}) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return string(respBody), fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

func (c *BugsClient) Delete(endpoint string, queryJSON string) (string, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+endpoint, bytes.NewReader([]byte(queryJSON)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return string(respBody), fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
