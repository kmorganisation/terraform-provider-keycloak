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

func resourceKeycloakOrganizationUserRoleAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationUserRoleAssignmentCreate,
		ReadContext:   resourceKeycloakOrganizationUserRoleAssignmentRead,
		DeleteContext: resourceKeycloakOrganizationUserRoleAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationUserRoleAssignmentImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func organizationUserRoleAssignmentId(realmId, organizationId, userId, roleId string) string {
	return fmt.Sprintf("%s/%s/%s/%s", realmId, organizationId, userId, roleId)
}

func resourceKeycloakOrganizationUserRoleAssignmentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	organizationId := data.Get("organization_id").(string)
	userId := data.Get("user_id").(string)
	roleId := data.Get("role_id").(string)

	err := keycloakClient.NewOrganizationUserRoleAssignment(ctx, realmId, organizationId, roleId, userId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(organizationUserRoleAssignmentId(realmId, organizationId, userId, roleId))

	return resourceKeycloakOrganizationUserRoleAssignmentRead(ctx, data, meta)
}

func resourceKeycloakOrganizationUserRoleAssignmentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	organizationId := data.Get("organization_id").(string)
	userId := data.Get("user_id").(string)
	roleId := data.Get("role_id").(string)

	_, err := keycloakClient.GetOrganizationUserRoleAssignment(ctx, realmId, organizationId, roleId, userId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	data.SetId(organizationUserRoleAssignmentId(realmId, organizationId, userId, roleId))

	return nil
}

func resourceKeycloakOrganizationUserRoleAssignmentDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	organizationId := data.Get("organization_id").(string)
	userId := data.Get("user_id").(string)
	roleId := data.Get("role_id").(string)

	err := keycloakClient.DeleteOrganizationUserRoleAssignment(ctx, realmId, organizationId, roleId, userId)
	if err != nil && !keycloak.ErrorIs404(err) {
		return diag.FromErr(err)
	}

	return nil
}

func resourceKeycloakOrganizationUserRoleAssignmentImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm_id}}/{{organization_id}}/{{user_id}}/{{role_id}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("organization_id", parts[1])
	d.Set("user_id", parts[2])
	d.Set("role_id", parts[3])
	d.SetId(organizationUserRoleAssignmentId(parts[0], parts[1], parts[2], parts[3]))

	diagnostics := resourceKeycloakOrganizationUserRoleAssignmentRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
