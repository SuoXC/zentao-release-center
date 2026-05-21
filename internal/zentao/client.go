package zentao

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout int) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:12345/api"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *Client) get(path string, params map[string]string) (json.RawMessage, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			if v != "" && v != "0" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %w", u.String(), err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return result.Data, nil
}

func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/../health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zentao-mini health check failed: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) GetProducts() (json.RawMessage, error) {
	return c.get("/products", nil)
}

func (c *Client) GetProjects(productID int) (json.RawMessage, error) {
	params := map[string]string{}
	if productID > 0 {
		params["productId"] = strconv.Itoa(productID)
	}
	return c.get("/projects", params)
}

func (c *Client) GetExecutions(projectID int) (json.RawMessage, error) {
	params := map[string]string{}
	if projectID > 0 {
		params["projectId"] = strconv.Itoa(projectID)
	}
	return c.get("/executions", params)
}

func (c *Client) GetBugs(productID, projectID int, status string, page, pageSize int) (json.RawMessage, error) {
	params := map[string]string{}
	if productID > 0 {
		params["productId"] = strconv.Itoa(productID)
	}
	if projectID > 0 {
		params["projectId"] = strconv.Itoa(projectID)
	}
	if status != "" {
		params["status"] = status
	}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	return c.get("/bugs", params)
}

func (c *Client) GetTasks(executionID, productID int, status string, page, pageSize int) (json.RawMessage, error) {
	params := map[string]string{}
	if executionID > 0 {
		params["executionId"] = strconv.Itoa(executionID)
	}
	if productID > 0 {
		params["productId"] = strconv.Itoa(productID)
	}
	if status != "" {
		params["status"] = status
	}
	if page > 0 {
		params["page"] = strconv.Itoa(page)
	}
	if pageSize > 0 {
		params["pageSize"] = strconv.Itoa(pageSize)
	}
	return c.get("/tasks", params)
}

func (c *Client) GetBug(bugID int) (json.RawMessage, error) {
	return c.get("/bugs", map[string]string{
		"productId": "0",
		"page":      "1",
		"pageSize":  "1",
	})
}
