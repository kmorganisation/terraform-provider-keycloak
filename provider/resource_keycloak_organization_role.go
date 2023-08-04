package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOrganizationRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationRoleCreate,
		ReadContext:   resourceKeycloakOrganizationRoleRead,
		DeleteContext: resourceKeycloakOrganizationRoleDelete,
		UpdateContext: resourceKeycloakOrganizationRoleUpdate,
		// This resource can be imported using {{realmName}}/{{organizationId}}/{{roleName}}. The Organization ID is displayed in the URL when editing it from the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organisation_id": {
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

func mapFromDataToOrganizationRole(data *schema.ResourceData) *keycloak.OrganizationRole {
	return &keycloak.OrganizationRole{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Realm:       data.Get("realm").(string),
		OrgId:       data.Get("organization_id").(string),
	}
}

func mapFromOrganizationRoleToData(data *schema.ResourceData, role *keycloak.OrganizationRole) {
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("realm", role.Realm)
	data.Set("organization_id", role.OrgId)
}

func resourceKeycloakOrganizationRoleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToOrganizationRole(data)

	err := keycloakClient.NewOrganizationRole(ctx, role)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationRoleToData(data, role)

	return resourceKeycloakOrganizationRoleRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRoleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	orgId := data.Get("organisation_id").(string)
	roleName := data.Get("name").(string)

	role, err := keycloakClient.GetOrganizationRole(ctx, realm, orgId, roleName)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromOrganizationRoleToData(data, role)

	return nil
}

func resourceKeycloakOrganizationRoleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToOrganizationRole(data)

	err := keycloakClient.UpdateOrganizationRole(ctx, role)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationRoleToData(data, role)

	return nil
}

func resourceKeycloakOrganizationRoleDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	role := mapFromDataToOrganizationRole(data)

	return diag.FromErr(keycloakClient.DeleteOrganizationRole(ctx, role.Realm, role.OrgId, role.Name))
}

func resourceKeycloakOrganizationRoleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmName}}/{{organizationId}}/{{roleName}}")
	}

	_, err := keycloakClient.GetOrganizationRole(ctx, parts[0], parts[1], parts[2])
	if err != nil {
		return nil, err
	}

	d.Set("realm", parts[0])
	d.Set("organization_id", parts[1])
	d.Set("name", parts[2])

	diagnostics := resourceKeycloakOrganizationRoleRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
