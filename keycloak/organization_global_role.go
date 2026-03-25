package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
)

type OrganizationGlobalRole struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	App         string `json:"app"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsDefault   bool   `json:"isDefault"`
	Default     bool   `json:"default,omitempty"`
}

func (keycloakClient *KeycloakClient) NewOrganizationGlobalRole(ctx context.Context, realmId string, role *OrganizationGlobalRole) error {
	role.Default = role.IsDefault

	body, location, err := keycloakClient.postRoot(ctx, fmt.Sprintf("/realms/%s/orgs/roles", realmId), role)
	if err != nil {
		return err
	}

	if location != "" {
		role.Id = getIdFromLocationHeader(location)
	}

	if role.Id == "" && len(body) != 0 {
		var createdRole OrganizationGlobalRole
		if err := json.Unmarshal(body, &createdRole); err == nil {
			role.Id = createdRole.Id
		}
	}

	if role.Id == "" {
		createdRole, err := keycloakClient.getOrganizationGlobalRoleByAppAndName(ctx, realmId, role.App, role.Name)
		if err != nil {
			return err
		}

		role.Id = createdRole.Id
	}

	createdRole, err := keycloakClient.GetOrganizationGlobalRole(ctx, realmId, role.Id)
	if err != nil {
		return err
	}

	*role = *createdRole

	return nil
}

func (keycloakClient *KeycloakClient) getOrganizationGlobalRoleByAppAndName(ctx context.Context, realmId, app, name string) (*OrganizationGlobalRole, error) {
	roles, err := keycloakClient.GetOrganizationGlobalRoles(ctx, realmId, app)
	if err != nil {
		return nil, err
	}

	var matchedRoles []*OrganizationGlobalRole
	for _, role := range roles {
		if role.Name == name {
			matchedRoles = append(matchedRoles, role)
		}
	}

	if len(matchedRoles) == 0 {
		return nil, fmt.Errorf("unable to find created organization global role with app %s and name %s", app, name)
	}

	if len(matchedRoles) > 1 {
		return nil, fmt.Errorf("multiple organization global roles found with app %s and name %s", app, name)
	}

	return matchedRoles[0], nil
}

func (keycloakClient *KeycloakClient) GetOrganizationGlobalRoles(ctx context.Context, realmId, app string) ([]*OrganizationGlobalRole, error) {
	var roles []*OrganizationGlobalRole
	var params map[string]string

	if app != "" {
		params = map[string]string{
			"app": app,
		}
	}

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/roles", realmId), &roles, params)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.Default {
			role.IsDefault = true
		}

		role.RealmId = realmId
	}

	return roles, nil
}

func (keycloakClient *KeycloakClient) GetOrganizationGlobalRole(ctx context.Context, realmId, roleId string) (*OrganizationGlobalRole, error) {
	var role *OrganizationGlobalRole

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/roles/%s", realmId, roleId), &role, nil)
	if err != nil {
		return nil, err
	}

	role.RealmId = realmId
	if role.Default {
		role.IsDefault = true
	}

	return role, nil
}

func (keycloakClient *KeycloakClient) UpdateOrganizationGlobalRole(ctx context.Context, realmId string, role *OrganizationGlobalRole) error {
	role.Default = role.IsDefault

	err := keycloakClient.putRoot(ctx, fmt.Sprintf("/realms/%s/orgs/roles/%s", realmId, role.Id), role)
	if err != nil {
		return err
	}

	role.RealmId = realmId

	return nil
}

func (keycloakClient *KeycloakClient) DeleteOrganizationGlobalRole(ctx context.Context, realmId, roleId string) error {
	return keycloakClient.deleteRoot(ctx, fmt.Sprintf("/realms/%s/orgs/roles/%s", realmId, roleId), nil)
}
