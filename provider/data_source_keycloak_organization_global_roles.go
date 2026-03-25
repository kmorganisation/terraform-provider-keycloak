package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakOrganizationGlobalRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakOrganizationGlobalRolesRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakOrganizationRoleApps, false),
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func sortOrganizationGlobalRoles(roles []*keycloak.OrganizationGlobalRole) {
	sort.Slice(roles, func(i, j int) bool {
		if roles[i].App != roles[j].App {
			return roles[i].App < roles[j].App
		}

		if roles[i].Name != roles[j].Name {
			return roles[i].Name < roles[j].Name
		}

		return roles[i].Id < roles[j].Id
	})
}

func dataSourceKeycloakOrganizationGlobalRolesRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	app := data.Get("app").(string)

	roles, err := keycloakClient.GetOrganizationGlobalRoles(ctx, realmId, app)
	if err != nil {
		return diag.FromErr(err)
	}

	sortOrganizationGlobalRoles(roles)

	roleViews := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		roleViews = append(roleViews, map[string]interface{}{
			"role_id":     role.Id,
			"app":         role.App,
			"name":        role.Name,
			"description": role.Description,
			"is_default":  role.IsDefault,
		})
	}

	if err := data.Set("roles", roleViews); err != nil {
		return diag.FromErr(fmt.Errorf("could not set roles: %w", err))
	}

	if app == "" {
		data.SetId(fmt.Sprintf("%s/all", realmId))
	} else {
		data.SetId(fmt.Sprintf("%s/%s", realmId, app))
	}

	return nil
}
