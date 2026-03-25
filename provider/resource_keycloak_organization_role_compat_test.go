package provider

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestWrapOrganizationRoleAmbiguousError(t *testing.T) {
	err := wrapOrganizationRoleAmbiguousError(&keycloak.ApiError{
		Code:    http.StatusConflict,
		Message: "conflict",
	})
	if err == nil {
		t.Fatal("expected wrapped error")
	}

	if !strings.Contains(err.Error(), "Role name is ambiguous. Use role_id-based resource.") {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestWrapOrganizationRoleAmbiguousErrorNonConflict(t *testing.T) {
	sourceError := &keycloak.ApiError{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	}

	err := wrapOrganizationRoleAmbiguousError(sourceError)
	if err == nil {
		t.Fatal("expected non-nil error")
	}

	if err.Error() != sourceError.Error() {
		t.Fatalf("expected original error message, got %s", err.Error())
	}
}
