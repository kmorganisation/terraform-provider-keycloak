package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceKeycloakOrganizationGlobalRoles_basic(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_ORGANIZATIONS")
	t.Parallel()

	roleOneName := fmt.Sprintf("zeta-%s", acctest.RandString(6))
	roleTwoName := fmt.Sprintf("alpha-%s", acctest.RandString(6))
	roleThreeName := fmt.Sprintf("beta-%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakOrganizationGlobalRoles(roleOneName, roleTwoName, roleThreeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.#", "3"),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.nexus_roles", "roles.#", "2"),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.0.app", "nexus"),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.0.name", roleTwoName),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.1.app", "nexus"),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.1.name", roleOneName),
					resource.TestCheckResourceAttr("data.keycloak_organization_global_roles.all_roles", "roles.2.app", "reporting"),
				),
			},
		},
	})
}

func testDataSourceKeycloakOrganizationGlobalRoles(roleOneName, roleTwoName, roleThreeName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_organization_global_role" "role_one" {
	realm_id = data.keycloak_realm.realm.id
	app      = "nexus"
	name     = "%s"
}

resource "keycloak_organization_global_role" "role_two" {
	realm_id = data.keycloak_realm.realm.id
	app      = "nexus"
	name     = "%s"
}

resource "keycloak_organization_global_role" "role_three" {
	realm_id = data.keycloak_realm.realm.id
	app      = "reporting"
	name     = "%s"
}

data "keycloak_organization_global_roles" "all_roles" {
	realm_id   = data.keycloak_realm.realm.id
	depends_on = [keycloak_organization_global_role.role_one, keycloak_organization_global_role.role_two, keycloak_organization_global_role.role_three]
}

data "keycloak_organization_global_roles" "nexus_roles" {
	realm_id   = data.keycloak_realm.realm.id
	app        = "nexus"
	depends_on = [keycloak_organization_global_role.role_one, keycloak_organization_global_role.role_two]
}
`, testAccRealm.Realm, roleOneName, roleTwoName, roleThreeName)
}
