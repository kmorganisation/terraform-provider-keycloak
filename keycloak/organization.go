package keycloak

import (
	"context"
	"fmt"
)

type Organization struct {
	Id          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	DisplayName string              `json:"displayName"`
	URL         string              `json:"url"`
	RealmName   string              `json:"realm"`
	Domains     []string            `json:"domains"`
	Attributes  map[string][]string `json:"attributes"`
}

func (keycloakClient *KeycloakClient) NewOrganization(ctx context.Context, organization *Organization) error {
	_, location, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/orgs", organization.RealmName), organization)
	if err != nil {
		return err
	}

	organization.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetOrganization(ctx context.Context, realmName string, orgId string) (*Organization, error) {
	var organization *Organization

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/orgs/%s", realmName, orgId), &organization, nil)
	if err != nil {
		return nil, err
	}

	return organization, nil
}

func (keycloakClient *KeycloakClient) GetOrganizations(ctx context.Context, realmName string) ([]*Organization, error) {
	var organizations []*Organization

	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/orgs", realmName), &organizations, nil)
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (keycloakClient *KeycloakClient) UpdateOrganization(ctx context.Context, organization *Organization) error {
	return keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/orgs/%s", organization.RealmName, organization.Id), organization)
}

func (keycloakClient *KeycloakClient) DeleteOrganization(ctx context.Context, realmName, id string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/orgs/%s", realmName, id), nil)
}
