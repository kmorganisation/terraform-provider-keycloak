package keycloak

import (
	"context"
	"fmt"
)

type OrganizationRole struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Realm       string `json:"-"`
	OrgId       string `json:"-"`
}

func (keycloakClient *KeycloakClient) NewOrganizationRole(ctx context.Context, role *OrganizationRole) error {
	_, _, err := keycloakClient.postRoot(ctx, fmt.Sprintf("/realms/%s/orgs/%s/roles", role.Realm, role.OrgId), role)
	if err != nil {
		return err
	}

	role, err = keycloakClient.GetOrganizationRole(ctx, role.Realm, role.OrgId, role.Name)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetOrganizationRole(ctx context.Context, realmName, orgId, roleName string) (*OrganizationRole, error) {
	var role *OrganizationRole

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/%s/roles/%s", realmName, orgId, roleName), &role, nil)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (keycloakClient *KeycloakClient) GetOrganizationRoles(ctx context.Context, realmName, orgId string) ([]*OrganizationRole, error) {
	var roles []*OrganizationRole

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/%s/roles", realmName, orgId), &roles, nil)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (keycloakClient *KeycloakClient) UpdateOrganizationRole(ctx context.Context, role *OrganizationRole) error {
	return keycloakClient.putRoot(ctx, fmt.Sprintf("/realms/%s/orgs/%s/roles/%s", role.Realm, role.OrgId, role.Name), role)
}

func (keycloakClient *KeycloakClient) DeleteOrganizationRole(ctx context.Context, realmName, orgId, roleName string) error {
	return keycloakClient.deleteRoot(ctx, fmt.Sprintf("/realms/%s/orgs/%s/roles/%s", realmName, orgId, roleName), nil)
}
