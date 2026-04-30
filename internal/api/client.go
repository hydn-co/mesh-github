package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	baseURL          = "https://api.github.com"
	maxRateLimitWait = 60 * time.Second
	maxRetries       = 3
	defaultPerPage   = 100
)

// Client is a GitHub API client using raw net/http.
type Client struct {
	token      string
	httpClient *http.Client
	org        string
}

// NewClient creates a new GitHub API client.
func NewClient(token, org string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		org:        org,
	}
}

// Member represents a GitHub organization member.
type Member struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// Team represents a GitHub team.
type Team struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Privacy     string `json:"privacy"`
}

// TeamMember represents a member of a GitHub team.
type TeamMember struct {
	Login string `json:"login"`
	ID    int64  `json:"id"`
	Role  string `json:"role"`
}

// Repository represents a GitHub repository.
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Archived    bool   `json:"archived"`
}

// Collaborator represents a repository collaborator.
type Collaborator struct {
	Login       string      `json:"login"`
	ID          int64       `json:"id"`
	Permissions Permissions `json:"permissions"`
	RoleName    string      `json:"role_name"`
}

// Permissions represents the permission set for a collaborator.
type Permissions struct {
	Admin    bool `json:"admin"`
	Maintain bool `json:"maintain"`
	Push     bool `json:"push"`
	Triage   bool `json:"triage"`
	Pull     bool `json:"pull"`
}

// AuditLogEntry represents a GitHub audit log entry.
type AuditLogEntry struct {
	Timestamp  int64          `json:"@timestamp"`
	Action     string         `json:"action"`
	Actor      string         `json:"actor"`
	ActorID    int64          `json:"actor_id"`
	Org        string         `json:"org"`
	CreatedAt  int64          `json:"created_at"`
	DocumentID string         `json:"_document_id"`
	User       string         `json:"user"`
	Team       string         `json:"team"`
	Repo       string         `json:"repo"`
	Data       map[string]any `json:"data"`
}

// ListOrgMembers returns all members of the organization.
func (c *Client) ListOrgMembers(ctx context.Context) ([]Member, error) {
	var all []Member
	page := 1

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s/orgs/%s/members?per_page=%d&page=%d", baseURL, c.org, defaultPerPage, page)
		var members []Member
		if err := c.get(ctx, url, &members); err != nil {
			return nil, fmt.Errorf("list org members page %d: %w", page, err)
		}

		if len(members) == 0 {
			break
		}

		all = append(all, members...)
		if len(members) < defaultPerPage {
			break
		}
		page++
	}

	return all, nil
}

// GetUser retrieves detailed user information.
func (c *Client) GetUser(ctx context.Context, login string) (*Member, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, login)
	var user Member
	if err := c.get(ctx, url, &user); err != nil {
		return nil, fmt.Errorf("get user %s: %w", login, err)
	}
	return &user, nil
}

// GetOrgMembership retrieves the membership role for a user in the org.
func (c *Client) GetOrgMembership(ctx context.Context, username string) (string, error) {
	url := fmt.Sprintf("%s/orgs/%s/memberships/%s", baseURL, c.org, username)
	var result struct {
		Role string `json:"role"`
	}
	if err := c.get(ctx, url, &result); err != nil {
		return "", fmt.Errorf("get org membership for %s: %w", username, err)
	}
	return result.Role, nil
}

// ListTeams returns all teams in the organization.
func (c *Client) ListTeams(ctx context.Context) ([]Team, error) {
	var all []Team
	page := 1

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s/orgs/%s/teams?per_page=%d&page=%d", baseURL, c.org, defaultPerPage, page)
		var teams []Team
		if err := c.get(ctx, url, &teams); err != nil {
			return nil, fmt.Errorf("list teams page %d: %w", page, err)
		}

		if len(teams) == 0 {
			break
		}

		all = append(all, teams...)
		if len(teams) < defaultPerPage {
			break
		}
		page++
	}

	return all, nil
}

// ListTeamMembers returns all members of a team.
func (c *Client) ListTeamMembers(ctx context.Context, teamSlug string) ([]TeamMember, error) {
	var all []TeamMember
	page := 1

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf(
			"%s/orgs/%s/teams/%s/members?per_page=%d&page=%d",
			baseURL,
			c.org,
			teamSlug,
			defaultPerPage,
			page,
		)
		var members []TeamMember
		if err := c.get(ctx, url, &members); err != nil {
			return nil, fmt.Errorf("list team members for %s page %d: %w", teamSlug, page, err)
		}

		if len(members) == 0 {
			break
		}

		all = append(all, members...)
		if len(members) < defaultPerPage {
			break
		}
		page++
	}

	return all, nil
}

// ListRepositories returns all repositories in the organization.
func (c *Client) ListRepositories(ctx context.Context) ([]Repository, error) {
	var all []Repository
	page := 1

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s/orgs/%s/repos?per_page=%d&page=%d&type=all", baseURL, c.org, defaultPerPage, page)
		var repos []Repository
		if err := c.get(ctx, url, &repos); err != nil {
			return nil, fmt.Errorf("list repositories page %d: %w", page, err)
		}

		if len(repos) == 0 {
			break
		}

		all = append(all, repos...)
		if len(repos) < defaultPerPage {
			break
		}
		page++
	}

	return all, nil
}

// ListCollaborators returns all collaborators for a repository.
func (c *Client) ListCollaborators(ctx context.Context, repoName string) ([]Collaborator, error) {
	var all []Collaborator
	page := 1

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s/repos/%s/%s/collaborators?per_page=%d&page=%d&affiliation=all",
			baseURL, c.org, repoName, defaultPerPage, page)
		var collabs []Collaborator
		if err := c.get(ctx, url, &collabs); err != nil {
			return nil, fmt.Errorf("list collaborators for %s page %d: %w", repoName, page, err)
		}

		if len(collabs) == 0 {
			break
		}

		all = append(all, collabs...)
		if len(collabs) < defaultPerPage {
			break
		}
		page++
	}

	return all, nil
}

// ListAuditLog returns audit log entries for the organization.
func (c *Client) ListAuditLog(ctx context.Context, after string, since time.Time) ([]AuditLogEntry, error) {
	var all []AuditLogEntry
	cursor := after

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		url := fmt.Sprintf("%s/orgs/%s/audit-log?per_page=%d&include=all", baseURL, c.org, defaultPerPage)
		if !since.IsZero() {
			url += fmt.Sprintf("&created_at=>=%s", since.UTC().Format(time.RFC3339))
		}
		if cursor != "" {
			url += fmt.Sprintf("&after=%s", cursor)
		}

		var entries []AuditLogEntry
		if err := c.get(ctx, url, &entries); err != nil {
			return nil, fmt.Errorf("list audit log: %w", err)
		}

		if len(entries) == 0 {
			break
		}

		all = append(all, entries...)

		// Use the last entry's document ID as cursor for next page
		if len(entries) < defaultPerPage {
			break
		}
		cursor = entries[len(entries)-1].DocumentID
	}

	return all, nil
}

// AddTeamMember adds a user to a team.
func (c *Client) AddTeamMember(ctx context.Context, teamSlug, username, role string) error {
	url := fmt.Sprintf("%s/orgs/%s/teams/%s/memberships/%s", baseURL, c.org, teamSlug, username)

	body := fmt.Sprintf(`{"role":"%s"}`, role)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, io.NopCloser(
		io.LimitReader(jsonReader(body), 1024),
	))
	if err != nil {
		return fmt.Errorf("create add team member request: %w", err)
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("add team member %s to %s: %w", username, teamSlug, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("add team member %s to %s: HTTP %d", username, teamSlug, resp.StatusCode)
	}

	return nil
}

// RemoveTeamMember removes a user from a team.
func (c *Client) RemoveTeamMember(ctx context.Context, teamSlug, username string) error {
	url := fmt.Sprintf("%s/orgs/%s/teams/%s/memberships/%s", baseURL, c.org, teamSlug, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("create remove team member request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return fmt.Errorf("remove team member %s from %s: %w", username, teamSlug, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("remove team member %s from %s: HTTP %d", username, teamSlug, resp.StatusCode)
	}

	return nil
}

func (c *Client) get(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	for attempt := range maxRetries {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
			if wait := rateLimitWait(resp); wait > 0 && wait <= maxRateLimitWait {
				_ = resp.Body.Close()
				slog.WarnContext(ctx, "rate limited, waiting", "wait", wait, "attempt", attempt+1)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(wait):
					continue
				}
			}
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts", maxRetries)
}

func rateLimitWait(resp *http.Response) time.Duration {
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		if ts, err := strconv.ParseInt(reset, 10, 64); err == nil {
			wait := time.Until(time.Unix(ts, 0))
			if wait > 0 {
				return wait
			}
		}
	}
	if retry := resp.Header.Get("Retry-After"); retry != "" {
		if seconds, err := strconv.Atoi(retry); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return 5 * time.Second
}

type stringReader struct {
	s string
	i int
}

func (r *stringReader) Read(p []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	n := copy(p, r.s[r.i:])
	r.i += n
	return n, nil
}

func jsonReader(s string) io.Reader {
	return &stringReader{s: s}
}
