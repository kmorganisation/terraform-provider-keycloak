package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func dataSourceKeycloakWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKeycloakWebhookRead,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"url": {
				Type: schema.TypeString,
				Required: true,
			},
			"secret": {
				Type: schema.TypeString,
				Required: true,
			},
			"event_types": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func dataSourceKeycloakWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmName := data.Get("realm").(string)
	webhookId := data.Get("id").(string)

	webhook, err := keycloakClient.GetWebhook(ctx, realmName, webhookId)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromWebhookToData(data, webhook)

	return nil
}
