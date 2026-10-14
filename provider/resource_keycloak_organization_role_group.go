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

var keycloakOrganizationTypes = []string{"broker", "shipper", "carrier"}

func resourceKeycloakOrganizationRoleGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationRoleGroupCreate,
		ReadContext:   resourceKeycloakOrganizationRoleGroupRead,
		UpdateContext: resourceKeycloakOrganizationRoleGroupUpdate,
		DeleteContext: resourceKeycloakOrganizationRoleGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationRoleGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// the server refuses to change it: members in organizations of the old type would keep
			// a group they can no longer see
			"organization_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(keycloakOrganizationTypes, false),
				Description:  "The organization type that sees the group. Leave empty for every type.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Organization global role ids granted to every member of the group.",
			},
		},
	}
}

func mapFromDataToOrganizationRoleGroup(data *schema.ResourceData) *keycloak.OrganizationRoleGroup {
	roleIds := make([]string, 0)
	if v, ok := data.GetOk("role_ids"); ok {
		for _, id := range v.(*schema.Set).List() {
			roleIds = append(roleIds, id.(string))
		}
	}

	return &keycloak.OrganizationRoleGroup{
		Id:               data.Id(),
		RealmId:          data.Get("realm_id").(string),
		Name:             data.Get("name").(string),
		Description:      data.Get("description").(string),
		OrganizationType: data.Get("organization_type").(string),
		RoleIds:          roleIds,
	}
}

func mapFromOrganizationRoleGroupToData(data *schema.ResourceData, group *keycloak.OrganizationRoleGroup) {
	data.SetId(group.Id)
	data.Set("group_id", group.Id)
	data.Set("realm_id", group.RealmId)
	data.Set("name", group.Name)
	data.Set("description", group.Description)
	data.Set("organization_type", group.OrganizationType)
	data.Set("role_ids", group.RoleIds)
}

func resourceKeycloakOrganizationRoleGroupCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	group := mapFromDataToOrganizationRoleGroup(data)

	err := keycloakClient.NewOrganizationRoleGroup(ctx, group.RealmId, group)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationRoleGroupToData(data, group)

	return resourceKeycloakOrganizationRoleGroupRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRoleGroupRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)

	group, err := keycloakClient.GetOrganizationRoleGroup(ctx, realmId, data.Id())
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromOrganizationRoleGroupToData(data, group)

	return nil
}

func resourceKeycloakOrganizationRoleGroupUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	group := mapFromDataToOrganizationRoleGroup(data)

	err := keycloakClient.UpdateOrganizationRoleGroup(ctx, group.RealmId, group)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeycloakOrganizationRoleGroupRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRoleGroupDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	realmId := data.Get("realm_id").(string)

	return diag.FromErr(keycloakClient.DeleteOrganizationRoleGroup(ctx, realmId, data.Id()))
}

func resourceKeycloakOrganizationRoleGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm_id}}/{{group_id}}")
	}

	_, err := keycloakClient.GetOrganizationRoleGroup(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakOrganizationRoleGroupRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
