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

func resourceKeycloakWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeycloakWebhookCreate,
		ReadContext:   resourceKeycloakWebhookRead,
		DeleteContext: resourceKeycloakWebhookDelete,
		UpdateContext: resourceKeycloakWebhookUpdate,
		// This resource can be imported using {{realmName}}/{{webhookId}}. The Webhook ID is displayed in the URL when editing it from the GUI
		Importer: &schema.ResourceImporter{
			StateContext: resourceKeycloakWebhookImport,
		},
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type: schema.TypeBool,
				Optional: true,
				Default: true,
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
				MinItems: 1,
				Required: true,
			},
		},
	}
}

func mapFromDataToWebhook(data *schema.ResourceData) *keycloak.Webhook {
	var eventTypes []string
	if v, ok := data.GetOk("event_types"); ok {
		for _, eventType := range v.(*schema.Set).List() {
			eventTypes = append(eventTypes, eventType.(string))
		}
	}

	webhook := &keycloak.Webhook{
		Id:         data.Id(),
		Enabled: 	data.Get("enabled").(bool),
		URL:        data.Get("url").(string),
		Secret:   	data.Get("secret").(string),
		EventTypes: eventTypes,
		RealmName:  data.Get("realm").(string),
	}

	return webhook
}

func mapFromWebhookToData(data *schema.ResourceData, webhook *keycloak.Webhook) {
	data.SetId(webhook.Id)
	data.Set("enabled", webhook.Enabled)
	data.Set("url", webhook.URL)
	data.Set("secret", webhook.Secret)
	data.Set("event_types", webhook.EventTypes)
	data.Set("realm", webhook.RealmName)
}

func resourceKeycloakWebhookCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	webhook := mapFromDataToWebhook(data)

	err := keycloakClient.NewWebhook(ctx, webhook)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromWebhookToData(data, webhook)

	return resourceKeycloakWebhookRead(ctx, data, meta)
}

func resourceKeycloakWebhookRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	id := data.Id()

	webhook, err := keycloakClient.GetWebhook(ctx, realm, id)
	if err != nil {
		return handleNotFoundError(ctx, err, data)
	}

	mapFromWebhookToData(data, webhook)

	return nil
}

func resourceKeycloakWebhookUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	webhook := mapFromDataToWebhook(data)

	err := keycloakClient.UpdateWebhook(ctx, webhook)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromWebhookToData(data, webhook)

	return nil
}

func resourceKeycloakWebhookDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realm := data.Get("realm").(string)
	id := data.Id()

	return diag.FromErr(keycloakClient.DeleteWebhook(ctx, realm, id))
}

func resourceKeycloakWebhookImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmName}}/{{webhookId}}")
	}

	_, err := keycloakClient.GetWebhook(ctx, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	d.Set("realm", parts[0])
	d.SetId(parts[1])

	diagnostics := resourceKeycloakWebhookRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
