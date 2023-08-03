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

func resourceKeycloakOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakOrganizationCreate,
		ReadContext:   resourceKeycloakOrganizationRead,
		DeleteContext: resourceKeycloakOrganizationDelete,
		UpdateContext: resourceKeycloakOrganizationUpdate,
		// This resource can be imported using {{realmName}}/{{organizationId}}. The Organization ID is displayed in the URL when editing it from the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakOrganizationImport,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				Set:      schema.HashString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func mapFromDataToOrganization(data *schema.ResourceData) *keycloak.Organization {
	attributes := map[string][]string{}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = strings.Split(value.(string), MULTIVALUE_ATTRIBUTE_SEPARATOR)
		}
	}

	organization := &keycloak.Organization{
		Id:          data.Id(),
		RealmName:   data.Get("realm").(string),
		Name:        data.Get("name").(string),
		DisplayName: data.Get("display_name").(string),
		URL:         data.Get("url").(string),
		Domains:     data.Get("domains").([]string),
		Attributes:  attributes,
	}

	return organization
}

func mapFromOrganizationToData(data *schema.ResourceData, organization *keycloak.Organization) {
	attributes := map[string]string{}
	for k, v := range organization.Attributes {
		attributes[k] = strings.Join(v, MULTIVALUE_ATTRIBUTE_SEPARATOR)
	}
	data.SetId(organization.Id)
	data.Set("realm", organization.RealmName)
	data.Set("name", organization.Name)
	data.Set("display_name", organization.DisplayName)
	data.Set("url", organization.URL)
	data.Set("domains", organization.Domains)
	data.Set("attributes", attributes)
}

func resourceKeycloakOrganizationCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	organization := mapFromDataToOrganization(data)

	err := keycloakClient.NewOrganization(ctx, organization)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationToData(data, organization)

	return resourceKeycloakOrganizationRead(ctx, data, meta)
}

func resourceKeycloakOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	id := data.Id()

	organization, err := keycloakClient.GetOrganization(ctx, realm, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromOrganizationToData(data, organization)

	return nil
}

func resourceKeycloakOrganizationUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	organization := mapFromDataToOrganization(data)

	err := keycloakClient.UpdateOrganization(ctx, organization)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationToData(data, organization)

	return nil
}

func resourceKeycloakOrganizationDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	id := data.Id()

	return diag.FromErr(keycloakClient.DeleteOrganization(ctx, realm, id))
}

func resourceKeycloakOrganizationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmName}}/{{organizationId}}")
	}

	_, err := keycloakClient.GetOrganization(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakOrganizationRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
