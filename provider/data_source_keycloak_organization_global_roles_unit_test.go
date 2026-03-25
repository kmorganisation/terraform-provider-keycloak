package provider

import (
	"testing"

	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestSortOrganizationGlobalRoles(t *testing.T) {
	roles := []*keycloak.OrganizationGlobalRole{
		{Id: "3", App: "reporting", Name: "z"},
		{Id: "2", App: "nexus", Name: "b"},
		{Id: "1", App: "nexus", Name: "a"},
		{Id: "0", App: "nexus", Name: "a"},
	}

	sortOrganizationGlobalRoles(roles)

	expected := []string{"0", "1", "2", "3"}
	for i, role := range roles {
		if role.Id != expected[i] {
			t.Fatalf("unexpected role order at index %d. expected %s, got %s", i, expected[i], role.Id)
		}
	}
}
