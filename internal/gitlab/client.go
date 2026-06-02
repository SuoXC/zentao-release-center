package gitlab

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

type Branch struct {
	Name               string `json:"name"`
	Commit             Commit `json:"commit"`
	Merged             bool   `json:"merged"`
	Protected          bool   `json:"protected"`
	Default            bool   `json:"default"`
	CanPush            bool   `json:"can_push"`
	WebURL             string `json:"web_url"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
}

type Commit struct {
	ID             string    `json:"id"`
	ShortID        string    `json:"short_id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	AuthorName     string    `json:"author_name"`
	AuthorEmail    string    `json:"author_email"`
	AuthoredDate   time.Time `json:"authored_date"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	CommittedDate  time.Time `json:"committed_date"`
	WebURL         string    `json:"web_url"`
}

type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
	HTTPURLToRepo    string `json:"http_url_to_repo"`
	SSHURLToRepo     string `json:"ssh_url_to_repo"`
	DefaultBranch     string `json:"default_branch"`
	Description       string `json:"description"`
}

type Pipeline struct {
	ID        int       `json:"id"`
	IID       int       `json:"iid"`
	ProjectID int       `json:"project_id"`
	Status    string    `json:"status"`
	Source    string    `json:"source"`
	Ref       string    `json:"ref"`
	SHA       string    `json:"sha"`
	WebURL    string    `json:"web_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PipelineEvent struct {
	ObjectKind       string            `json:"object_kind"`
	ObjectAttributes PipelineEventAttr `json:"object_attributes"`
	Project          PipelineProject   `json:"project"`
	Commit           PipelineCommit    `json:"commit"`
	MergeRequest     *MergeRequest     `json:"merge_request,omitempty"`
	Builds           []Build           `json:"builds,omitempty"`
	Variables        []PipelineVariable `json:"variables,omitempty"`
}

type PipelineVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PipelineEventAttr struct {
	ID         int       `json:"id"`
	IID        int       `json:"iid"`
	ProjectID  int       `json:"project_id"`
	Status     string    `json:"status"`
	Source     string    `json:"source"`
	Ref        string    `json:"ref"`
	SHA        string    `json:"sha"`
	BeforeSHA  string    `json:"before_sha"`
	WebURL     string    `json:"web_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration   int       `json:"duration"`
}

type PipelineProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
}

type PipelineCommit struct {
	ID             string `json:"id"`
	Message        string `json:"message"`
	Timestamp      string `json:"timestamp"`
	URL            string `json:"url"`
	Author         Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MergeRequest struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	Title     string `json:"title"`
	Source    string `json:"source_branch"`
	Target    string `json:"target_branch"`
	State     string `json:"state"`
	URL       string `json:"url"`
}

type Build struct {
	ID        int       `json:"id"`
	Stage     string    `json:"stage"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (c *Client) doRequest(path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + "/api/v4" + path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitlab api %s returned %d: %s", path, resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) SearchProjects(query string) ([]Project, error) {
	data, err := c.doRequest("/projects", map[string]string{
		"search": query,
		"per_page": "20",
	})
	if err != nil {
		return nil, err
	}
	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("unmarshal projects: %w", err)
	}
	return projects, nil
}

func (c *Client) GetProject(projectID int) (*Project, error) {
	data, err := c.doRequest(fmt.Sprintf("/projects/%d", projectID), nil)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("unmarshal project: %w", err)
	}
	return &project, nil
}

func (c *Client) ListBranches(projectID int) ([]Branch, error) {
	data, err := c.doRequest(fmt.Sprintf("/projects/%d/repository/branches", projectID), map[string]string{
		"per_page": "100",
	})
	if err != nil {
		return nil, err
	}
	var branches []Branch
	if err := json.Unmarshal(data, &branches); err != nil {
		return nil, fmt.Errorf("unmarshal branches: %w", err)
	}
	return branches, nil
}

func (c *Client) CreateBranch(projectID int, branchName, ref string) error {
	u := fmt.Sprintf("%s/api/v4/projects/%d/repository/branches", c.baseURL, projectID)
	form := url.Values{}
	form.Set("branch", branchName)
	form.Set("ref", ref)

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.URL.RawQuery = form.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create branch returned %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) GetBranch(projectID int, branchName string) (*Branch, error) {
	encoded := url.PathEscape(branchName)
	data, err := c.doRequest(fmt.Sprintf("/projects/%d/repository/branches/%s", projectID, encoded), nil)
	if err != nil {
		return nil, err
	}
	var branch Branch
	if err := json.Unmarshal(data, &branch); err != nil {
		return nil, fmt.Errorf("unmarshal branch: %w", err)
	}
	return &branch, nil
}

func (c *Client) ListPipelines(projectID int, ref string) ([]Pipeline, error) {
	params := map[string]string{
		"per_page": "20",
	}
	if ref != "" {
		params["ref"] = ref
	}
	data, err := c.doRequest(fmt.Sprintf("/projects/%d/pipelines", projectID), params)
	if err != nil {
		return nil, err
	}
	var pipelines []Pipeline
	if err := json.Unmarshal(data, &pipelines); err != nil {
		return nil, fmt.Errorf("unmarshal pipelines: %w", err)
	}
	return pipelines, nil
}

func (c *Client) GetPipeline(projectID, pipelineID int) (*Pipeline, error) {
	data, err := c.doRequest(fmt.Sprintf("/projects/%d/pipelines/%d", projectID, pipelineID), nil)
	if err != nil {
		return nil, err
	}
	var pipeline Pipeline
	if err := json.Unmarshal(data, &pipeline); err != nil {
		return nil, fmt.Errorf("unmarshal pipeline: %w", err)
	}
	return &pipeline, nil
}
