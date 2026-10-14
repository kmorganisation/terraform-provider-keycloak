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

func TestAccKeycloakOrganizationRoleGroup_basic(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_ORGANIZATIONS")
	t.Parallel()

	prefix := acctest.RandomWithPrefix("tf-acc-org-role-group")
	resourceName := "keycloak_organization_role_group.group"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOrganizationRoleGroupDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationRoleGroup(prefix, "broker", "keycloak_organization_global_role.view.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationRoleGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", prefix),
					resource.TestCheckResourceAttr(resourceName, "organization_type", "broker"),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "1"),
				),
			},
			{
				Config: testKeycloakOrganizationRoleGroup(prefix, "broker", "keycloak_organization_global_role.view.id, keycloak_organization_global_role.manage.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationRoleGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "role_ids.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource %s not found in state", resourceName)
					}

					return fmt.Sprintf("%s/%s", rs.Primary.Attributes["realm_id"], rs.Primary.ID), nil
				},
			},
		},
	})
}

func testAccCheckKeycloakOrganizationRoleGroupExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("failed to find %s in state", resourceName)
		}

		_, err := keycloakClient.GetOrganizationRoleGroup(testCtx, rs.Primary.Attributes["realm_id"], rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting organization role group with id %s: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckKeycloakOrganizationRoleGroupDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for name, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_organization_role_group" || strings.HasPrefix(name, "data") {
				continue
			}

			group, err := keycloakClient.GetOrganizationRoleGroup(testCtx, rs.Primary.Attributes["realm_id"], rs.Primary.ID)
			if err != nil && !keycloak.ErrorIs404(err) {
				return fmt.Errorf("error getting organization role group with id %s: %w", rs.Primary.ID, err)
			}

			if group != nil && err == nil {
				return fmt.Errorf("%s with id %s still exists", name, rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOrganizationRoleGroup(name, organizationType, roleIds string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%[1]s"
}

resource "keycloak_organization_global_role" "view" {
	realm_id = data.keycloak_realm.realm.id
	app      = "nexus"
	name     = "%[2]s-view"
}

resource "keycloak_organization_global_role" "manage" {
	realm_id = data.keycloak_realm.realm.id
	app      = "nexus"
	name     = "%[2]s-manage"
}

resource "keycloak_organization_role_group" "group" {
	realm_id          = data.keycloak_realm.realm.id
	name              = "%[2]s"
	description       = "Day-to-day operations"
	organization_type = "%[3]s"
	role_ids          = [%[4]s]
}
`, testAccRealm.Realm, name, organizationType, roleIds)
}
