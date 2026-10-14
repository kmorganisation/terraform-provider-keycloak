package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
)

// OrganizationRoleGroup is a predefined role group: a realm-wide set of organization roles that
// every organization of the group's type can grant to its members.
type OrganizationRoleGroup struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	OrganizationType string   `json:"organizationType,omitempty"`
	RoleIds          []string `json:"roleIds"`
}

func (keycloakClient *KeycloakClient) NewOrganizationRoleGroup(ctx context.Context, realmId string, group *OrganizationRoleGroup) error {
	body, location, err := keycloakClient.postRoot(ctx, fmt.Sprintf("/realms/%s/orgs/role-groups", realmId), group)
	if err != nil {
		return err
	}

	if location != "" {
		group.Id = getIdFromLocationHeader(location)
	}

	if group.Id == "" && len(body) != 0 {
		var created OrganizationRoleGroup
		if err := json.Unmarshal(body, &created); err == nil {
			group.Id = created.Id
		}
	}

	if group.Id == "" {
		return fmt.Errorf("unable to determine id of created organization role group %s", group.Name)
	}

	created, err := keycloakClient.GetOrganizationRoleGroup(ctx, realmId, group.Id)
	if err != nil {
		return err
	}

	*group = *created

	return nil
}

func (keycloakClient *KeycloakClient) GetOrganizationRoleGroups(ctx context.Context, realmId string) ([]*OrganizationRoleGroup, error) {
	var groups []*OrganizationRoleGroup

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/role-groups", realmId), &groups, nil)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		group.RealmId = realmId
	}

	return groups, nil
}

func (keycloakClient *KeycloakClient) GetOrganizationRoleGroup(ctx context.Context, realmId, groupId string) (*OrganizationRoleGroup, error) {
	var group *OrganizationRoleGroup

	err := keycloakClient.getRoot(ctx, fmt.Sprintf("/realms/%s/orgs/role-groups/%s", realmId, groupId), &group, nil)
	if err != nil {
		return nil, err
	}

	group.RealmId = realmId

	return group, nil
}

func (keycloakClient *KeycloakClient) UpdateOrganizationRoleGroup(ctx context.Context, realmId string, group *OrganizationRoleGroup) error {
	return keycloakClient.putRoot(ctx, fmt.Sprintf("/realms/%s/orgs/role-groups/%s", realmId, group.Id), group)
}

func (keycloakClient *KeycloakClient) DeleteOrganizationRoleGroup(ctx context.Context, realmId, groupId string) error {
	return keycloakClient.deleteRoot(ctx, fmt.Sprintf("/realms/%s/orgs/role-groups/%s", realmId, groupId), nil)
}
