package api

import (
	"context"
	"fmt"
	"time"

	"github.com/fgrzl/enumerators"
)

func pagedEnumerator[T any](ctx context.Context, fetch func(page int) ([]T, error)) enumerators.Enumerator[T] {
	page := 1

	return enumerators.PageItemEnumerator(func() ([]T, bool, error) {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		items, err := fetch(page)
		if err != nil {
			return nil, false, err
		}

		if len(items) < defaultPerPage {
			return items, false, nil
		}

		page++
		return items, true, nil
	})
}

func cursorEnumerator[T any](
	ctx context.Context,
	fetch func(cursor string) ([]T, string, error),
) enumerators.Enumerator[T] {
	cursor := ""

	return enumerators.PageItemEnumerator(func() ([]T, bool, error) {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}

		items, nextCursor, err := fetch(cursor)
		if err != nil {
			return nil, false, err
		}

		if nextCursor == "" {
			return items, false, nil
		}

		cursor = nextCursor
		return items, true, nil
	})
}

func (c *Client) MemberEnumerator(ctx context.Context) enumerators.Enumerator[Member] {
	return pagedEnumerator(ctx, func(page int) ([]Member, error) {
		url := fmt.Sprintf("%s/orgs/%s/members?per_page=%d&page=%d", baseURL, c.org, defaultPerPage, page)
		var members []Member
		if err := c.get(ctx, url, &members); err != nil {
			return nil, fmt.Errorf("list org members page %d: %w", page, err)
		}

		return members, nil
	})
}

func (c *Client) TeamEnumerator(ctx context.Context) enumerators.Enumerator[Team] {
	return pagedEnumerator(ctx, func(page int) ([]Team, error) {
		url := fmt.Sprintf("%s/orgs/%s/teams?per_page=%d&page=%d", baseURL, c.org, defaultPerPage, page)
		var teams []Team
		if err := c.get(ctx, url, &teams); err != nil {
			return nil, fmt.Errorf("list teams page %d: %w", page, err)
		}

		return teams, nil
	})
}

func (c *Client) TeamMemberEnumerator(ctx context.Context, teamSlug string) enumerators.Enumerator[TeamMember] {
	return pagedEnumerator(ctx, func(page int) ([]TeamMember, error) {
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

		return members, nil
	})
}

func (c *Client) RepositoryEnumerator(ctx context.Context) enumerators.Enumerator[Repository] {
	return pagedEnumerator(ctx, func(page int) ([]Repository, error) {
		url := fmt.Sprintf(
			"%s/orgs/%s/repos?type=all&per_page=%d&page=%d",
			baseURL,
			c.org,
			defaultPerPage,
			page,
		)
		var repositories []Repository
		if err := c.get(ctx, url, &repositories); err != nil {
			return nil, fmt.Errorf("list org repositories page %d: %w", page, err)
		}

		return repositories, nil
	})
}

func (c *Client) RepositoryCollaboratorEnumerator(
	ctx context.Context,
	owner string,
	repo string,
) enumerators.Enumerator[RepositoryCollaborator] {
	return pagedEnumerator(ctx, func(page int) ([]RepositoryCollaborator, error) {
		url := fmt.Sprintf(
			"%s/repos/%s/%s/collaborators?affiliation=all&per_page=%d&page=%d",
			baseURL,
			owner,
			repo,
			defaultPerPage,
			page,
		)
		var collaborators []RepositoryCollaborator
		if err := c.get(ctx, url, &collaborators); err != nil {
			return nil, fmt.Errorf("list repository collaborators for %s/%s page %d: %w", owner, repo, page, err)
		}

		return collaborators, nil
	})
}

func (c *Client) AuditLogEnumerator(
	ctx context.Context,
	after string,
	since time.Time,
) enumerators.Enumerator[AuditLogEntry] {
	return cursorEnumerator(ctx, func(cursor string) ([]AuditLogEntry, string, error) {
		resolvedCursor := cursor
		if resolvedCursor == "" {
			resolvedCursor = after
		}

		entries, nextCursor, err := c.getAuditLogPage(ctx, resolvedCursor, since)
		if err != nil {
			return nil, "", fmt.Errorf("list audit log: %w", err)
		}
		return entries, nextCursor, nil
	})
}
