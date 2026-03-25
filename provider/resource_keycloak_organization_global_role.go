package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOrganizationGlobalRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationGlobalRoleCreate,
		ReadContext:   resourceKeycloakOrganizationGlobalRoleRead,
		DeleteContext: resourceKeycloakOrganizationGlobalRoleDelete,
		UpdateContext: resourceKeycloakOrganizationGlobalRoleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationGlobalRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakOrganizationRoleApps, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func mapFromDataToOrganizationGlobalRole(data *schema.ResourceData) *keycloak.OrganizationGlobalRole {
	return &keycloak.OrganizationGlobalRole{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		App:         data.Get("app").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		IsDefault:   data.Get("is_default").(bool),
	}
}

func mapFromOrganizationGlobalRoleToData(data *schema.ResourceData, role *keycloak.OrganizationGlobalRole) {
	data.SetId(role.Id)
	data.Set("role_id", role.Id)
	data.Set("realm_id", role.RealmId)
	data.Set("app", role.App)
	data.Set("name", role.Name)
	data.Set("description", role.Description)
	data.Set("is_default", role.IsDefault)
}

func resourceKeycloakOrganizationGlobalRoleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	role := mapFromDataToOrganizationGlobalRole(data)

	err := keycloakClient.NewOrganizationGlobalRole(ctx, role.RealmId, role)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationGlobalRoleToData(data, role)

	return resourceKeycloakOrganizationGlobalRoleRead(ctx, data, meta)
}

func resourceKeycloakOrganizationGlobalRoleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)
	roleId := data.Id()

	if roleId == "" {
		roleId = data.Get("role_id").(string)
	}

	role, err := keycloakClient.GetOrganizationGlobalRole(ctx, realmId, roleId)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromOrganizationGlobalRoleToData(data, role)

	return nil
}

func resourceKeycloakOrganizationGlobalRoleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	role := mapFromDataToOrganizationGlobalRole(data)

	err := keycloakClient.UpdateOrganizationGlobalRole(ctx, role.RealmId, role)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationGlobalRoleToData(data, role)

	return resourceKeycloakOrganizationGlobalRoleRead(ctx, data, meta)
}

func resourceKeycloakOrganizationGlobalRoleDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)

	return diag.FromErr(keycloakClient.DeleteOrganizationGlobalRole(ctx, realmId, data.Id()))
}

func resourceKeycloakOrganizationGlobalRoleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm_id}}/{{role_id}}")
	}

	_, err := keycloakClient.GetOrganizationGlobalRole(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakOrganizationGlobalRoleRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
