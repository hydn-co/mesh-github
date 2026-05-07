package collectors

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hydn-co/mesh-github/internal/api"
)

func TestNewAuditLogUnavailableErrorIncludesHint(t *testing.T) {
	baseErr := fmt.Errorf(
		"list audit log: %w",
		&api.HTTPStatusError{StatusCode: http.StatusNotFound, Body: `{"message":"Not Found"}`},
	)

	err := newAuditLogUnavailableError(baseErr)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), auditLogUnavailableHint) {
		t.Fatalf("expected hint in error, got %q", err.Error())
	}

	if !api.IsAuditLogUnavailable(err) {
		t.Fatal("expected wrapped error to remain classified as audit log unavailable")
	}
}
