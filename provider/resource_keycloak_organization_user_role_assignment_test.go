package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOrganizationUserRoleAssignment_basic(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_ORGANIZATIONS")
	t.Parallel()

	organizationName := acctest.RandomWithPrefix("tf-acc-org")
	username := strings.ToLower(acctest.RandomWithPrefix("tf-acc-user"))
	roleName := acctest.RandomWithPrefix("tf-acc-org-role")
	resourceName := "keycloak_organization_user_role_assignment.assignment"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOrganizationUserRoleAssignmentDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationUserRoleAssignmentBasic(organizationName, username, roleName),
				Check:  testAccCheckKeycloakOrganizationUserRoleAssignmentExists(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckKeycloakOrganizationUserRoleAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getOrganizationUserRoleAssignmentFromState(state, resourceName)
		return err
	}
}

func testAccCheckKeycloakOrganizationUserRoleAssignmentDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for name, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_organization_user_role_assignment" || strings.HasPrefix(name, "data") {
				continue
			}

			realmId := rs.Primary.Attributes["realm_id"]
			organizationId := rs.Primary.Attributes["organization_id"]
			userId := rs.Primary.Attributes["user_id"]
			roleId := rs.Primary.Attributes["role_id"]

			assignment, err := keycloakClient.GetOrganizationUserRoleAssignment(testCtx, realmId, organizationId, roleId, userId)
			if err != nil && !keycloak.ErrorIs404(err) {
				return fmt.Errorf("error getting organization user role assignment: %w", err)
			}

			if assignment != nil && err == nil {
				return fmt.Errorf("%s still exists", name)
			}
		}

		return nil
	}
}

func getOrganizationUserRoleAssignmentFromState(state *terraform.State, resourceName string) (*keycloak.OrganizationUserRoleAssignment, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("failed to find %s in state", resourceName)
	}

	realmId := rs.Primary.Attributes["realm_id"]
	organizationId := rs.Primary.Attributes["organization_id"]
	userId := rs.Primary.Attributes["user_id"]
	roleId := rs.Primary.Attributes["role_id"]

	assignment, err := keycloakClient.GetOrganizationUserRoleAssignment(testCtx, realmId, organizationId, roleId, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting organization user role assignment: %w", err)
	}

	return assignment, nil
}

func testKeycloakOrganizationUserRoleAssignmentBasic(organizationName, username, roleName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_organization" "organization" {
	realm = data.keycloak_realm.realm.id
	name  = "%s"
}

resource "keycloak_user" "user" {
	realm_id = data.keycloak_realm.realm.id
	username = "%s"
}

resource "keycloak_organization_global_role" "role" {
	realm_id = data.keycloak_realm.realm.id
	app      = "nexus"
	name     = "%s"
}

resource "keycloak_organization_user_role_assignment" "assignment" {
	realm_id        = data.keycloak_realm.realm.id
	organization_id = keycloak_organization.organization.id
	user_id         = keycloak_user.user.id
	role_id         = keycloak_organization_global_role.role.role_id
}
`, testAccRealm.Realm, organizationName, username, roleName)
}
