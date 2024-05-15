package keycloak

import (
	"context"
	"fmt"
)

type Webhook struct {
	Id          string              `json:"id,omitempty"`
	Enabled		bool 				`json:"enabled"`
	URL         string              `json:"url"`
	Secret      string              `json:"secret"`
	EventTypes  []string            `json:"eventTypes"`
	RealmName   string              `json:"realm"`
}

func (keycloakClient *KeycloakClient) NewWebhook(ctx context.Context, webhook *Webhook) error {
	_, location, err := keycloakClient.postRoot(ctx, fmt.Sprintf("/realms/%s/webhooks", webhook.RealmName), webhook)
	if err != nil {
		return err
	}

	webhook.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetWebhook(ctx context.Context, realmName string, webhookId string) (*Webhook, error) {
	var webhook *Webhook

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/webhooks/%s", realmName, webhookId), &webhook, nil)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

func (keycloakClient *KeycloakClient) GetWebhooks(ctx context.Context, realmName string) ([]*Webhook, error) {
	var webhooks []*Webhook

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/webhooks", realmName), &webhooks, nil)
	if err != nil {
		return nil, err
	}

	return webhooks, nil
}

func (keycloakClient *KeycloakClient) UpdateWebhook(ctx context.Context, webhook *Webhook) error {
	return keycloakClient.putRoot(ctx, fmt.Sprintf("/realms/%s/webhooks/%s", webhook.RealmName, webhook.Id), webhook)
}

func (keycloakClient *KeycloakClient) DeleteWebhook(ctx context.Context, realmName, id string) error {
	return keycloakClient.deleteRoot(ctx, fmt.Sprintf("/realms/%s/webhooks/%s", realmName, id), nil)
}
