package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakOrganizationRead,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
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
				Set:      schema.HashString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				Default:  "broker",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "broker" && v != "shipper" && v != "carrier" {
						errs = append(errs, fmt.Errorf("%q must be either 'broker', 'shipper' or 'carrier', got: %q", key, v))
					}
					return warns, errs
				},
			},
		},
	}
}

func dataSourceKeycloakOrganizationRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmName := data.Get("realm").(string)
	organizationId := data.Get("id").(string)

	organization, err := keycloakClient.GetOrganization(ctx, realmName, organizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromOrganizationToData(data, organization)

	return nil
}
