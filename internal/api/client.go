package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL          = "https://api.github.com"
	apiVersion       = "2026-03-10"
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

// HTTPStatusError captures non-2xx responses from the GitHub REST API.
type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	if e == nil {
		return "HTTP error"
	}

	if e.Body == "" {
		return fmt.Sprintf("HTTP %d", e.StatusCode)
	}

	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

func newHTTPStatusError(statusCode int, body []byte) error {
	return &HTTPStatusError{
		StatusCode: statusCode,
		Body:       strings.TrimSpace(string(body)),
	}
}

func IsAuditLogUnavailable(err error) bool {
	var statusErr *HTTPStatusError
	if !errors.As(err, &statusErr) {
		return false
	}

	return statusErr.StatusCode == http.StatusNotFound
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

// GetTeamMembershipRole retrieves the membership role for a user in a team.
func (c *Client) GetTeamMembershipRole(ctx context.Context, teamSlug, username string) (string, error) {
	url := fmt.Sprintf("%s/orgs/%s/teams/%s/memberships/%s", baseURL, c.org, teamSlug, username)
	var result struct {
		Role string `json:"role"`
	}
	if err := c.get(ctx, url, &result); err != nil {
		return "", fmt.Errorf("get team membership for %s/%s: %w", teamSlug, username, err)
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

// ListAuditLog returns audit log entries for the organization.
func (c *Client) ListAuditLog(ctx context.Context, after string, since time.Time) ([]AuditLogEntry, error) {
	var all []AuditLogEntry
	cursor := after

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		entries, nextCursor, err := c.getAuditLogPage(ctx, cursor, since)
		if err != nil {
			return nil, fmt.Errorf("list audit log: %w", err)
		}

		if len(entries) == 0 {
			break
		}

		all = append(all, entries...)
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}

	return all, nil
}

func (c *Client) getAuditLogPage(ctx context.Context, after string, since time.Time) ([]AuditLogEntry, string, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/orgs/%s/audit-log", baseURL, c.org))
	if err != nil {
		return nil, "", fmt.Errorf("parse audit log url: %w", err)
	}

	query := endpoint.Query()
	query.Set("per_page", strconv.Itoa(defaultPerPage))
	query.Set("include", "all")
	query.Set("order", "asc")
	if !since.IsZero() {
		query.Set("phrase", fmt.Sprintf("created:>=%s", since.UTC().Format(time.RFC3339)))
	}
	if after != "" {
		query.Set("after", after)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("create audit log request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.doWithRetry(ctx, req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, "", newHTTPStatusError(resp.StatusCode, body)
	}

	var entries []AuditLogEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, "", fmt.Errorf("decode audit log response: %w", err)
	}

	return entries, cursorFromLinkHeader(resp.Header.Get("Link"), "next"), nil
}

func cursorFromLinkHeader(linkHeader, rel string) string {
	for _, segment := range strings.Split(linkHeader, ",") {
		parts := strings.Split(strings.TrimSpace(segment), ";")
		if len(parts) < 2 {
			continue
		}

		matchesRel := false
		for _, attribute := range parts[1:] {
			if strings.TrimSpace(attribute) == fmt.Sprintf("rel=%q", rel) {
				matchesRel = true
				break
			}
		}
		if !matchesRel {
			continue
		}

		parsedURL, err := url.Parse(strings.Trim(parts[0], "<> "))
		if err != nil {
			continue
		}

		if cursor := parsedURL.Query().Get("after"); cursor != "" {
			return cursor
		}
		if cursor := parsedURL.Query().Get("before"); cursor != "" {
			return cursor
		}
	}

	return ""
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
		return newHTTPStatusError(resp.StatusCode, body)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
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
