package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakOrganizationGlobalRole_basic(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_ORGANIZATIONS")
	t.Parallel()

	roleName := acctest.RandomWithPrefix("tf-acc-org-role")
	description := acctest.RandomWithPrefix("tf-acc-description")
	resourceName := "keycloak_organization_global_role.role"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOrganizationGlobalRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationGlobalRoleBasic(roleName, description, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationGlobalRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app", "nexus"),
					resource.TestCheckResourceAttr(resourceName, "name", roleName),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "is_default", "false"),
					resource.TestMatchResourceAttr(resourceName, "role_id", regexp.MustCompile(".+")),
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

func TestAccKeycloakOrganizationGlobalRole_update(t *testing.T) {
	skipIfEnvNotSet(t, "KEYCLOAK_TEST_ORGANIZATIONS")
	t.Parallel()

	roleName := acctest.RandomWithPrefix("tf-acc-org-role")
	descriptionOne := acctest.RandomWithPrefix("tf-acc-description")
	descriptionTwo := acctest.RandomWithPrefix("tf-acc-description")
	resourceName := "keycloak_organization_global_role.role"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakOrganizationGlobalRoleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOrganizationGlobalRoleBasic(roleName, descriptionOne, false),
				Check:  testAccCheckKeycloakOrganizationGlobalRoleExists(resourceName),
			},
			{
				Config: testKeycloakOrganizationGlobalRoleBasic(roleName, descriptionTwo, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOrganizationGlobalRoleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionTwo),
					resource.TestCheckResourceAttr(resourceName, "is_default", "true"),
				),
			},
		},
	})
}

func testAccCheckKeycloakOrganizationGlobalRoleExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getOrganizationGlobalRoleFromState(state, resourceName)
		return err
	}
}

func testAccCheckKeycloakOrganizationGlobalRoleDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for name, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_organization_global_role" || strings.HasPrefix(name, "data") {
				continue
			}

			id := rs.Primary.ID
			realmId := rs.Primary.Attributes["realm_id"]

			role, err := keycloakClient.GetOrganizationGlobalRole(testCtx, realmId, id)
			if err != nil && !keycloak.ErrorIs404(err) {
				return fmt.Errorf("error getting organization global role with id %s: %w", id, err)
			}

			if role != nil && err == nil {
				return fmt.Errorf("%s with id %s still exists", name, id)
			}
		}

		return nil
	}
}

func getOrganizationGlobalRoleFromState(state *terraform.State, resourceName string) (*keycloak.OrganizationGlobalRole, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("failed to find %s in state", resourceName)
	}

	id := rs.Primary.ID
	realmId := rs.Primary.Attributes["realm_id"]

	role, err := keycloakClient.GetOrganizationGlobalRole(testCtx, realmId, id)
	if err != nil {
		return nil, fmt.Errorf("error getting organization global role with id %s: %w", id, err)
	}

	return role, nil
}

func testKeycloakOrganizationGlobalRoleBasic(name, description string, isDefault bool) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_organization_global_role" "role" {
	realm_id    = data.keycloak_realm.realm.id
	app         = "nexus"
	name        = "%s"
	description = "%s"
	is_default  = %t
}
`, testAccRealm.Realm, name, description, isDefault)
}
