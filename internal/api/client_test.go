package api

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestIsAuditLogUnavailableReturnsTrueForWrapped404(t *testing.T) {
	err := fmt.Errorf("list audit log: %w", newHTTPStatusError(http.StatusNotFound, []byte(`{"message":"Not Found"}`)))

	if !IsAuditLogUnavailable(err) {
		t.Fatal("expected wrapped 404 audit log error to be classified as unavailable")
	}
}

func TestIsAuditLogUnavailableReturnsFalseForNon404(t *testing.T) {
	err := fmt.Errorf("list audit log: %w", newHTTPStatusError(http.StatusForbidden, []byte(`{"message":"Forbidden"}`)))

	if IsAuditLogUnavailable(err) {
		t.Fatal("expected non-404 audit log error to remain fatal")
	}
}

func TestIsAuditLogUnavailableReturnsFalseForNonHTTPStatusError(t *testing.T) {
	err := fmt.Errorf("list audit log: %w", errors.New("boom"))

	if IsAuditLogUnavailable(err) {
		t.Fatal("expected non-http status error to remain unclassified")
	}
}
