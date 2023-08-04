package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakOrganizationRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakOrganizationRoleRead,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceKeycloakOrganizationRoleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmName := data.Get("realm").(string)
	organizationId := data.Get("organization_id").(string)
	roleName := data.Get("name").(string)

	role, err := keycloakClient.GetOrganizationRole(ctx, realmName, organizationId, roleName)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationRoleToData(data, role)

	return nil
}
