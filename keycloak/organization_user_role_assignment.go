package keycloak

import (
	"context"
	"fmt"
)

type OrganizationUserRoleAssignment struct {
	RealmId        string `json:"-"`
	OrganizationId string `json:"-"`
	UserId         string `json:"-"`
	RoleId         string `json:"-"`
}

func organizationRoleByIdPath(realmId, organizationId, roleId string) string {
	return fmt.Sprintf("/realms/%s/orgs/%s/roles/id/%s", realmId, organizationId, roleId)
}

func organizationUserRoleAssignmentPath(realmId, organizationId, roleId, userId string) string {
	return fmt.Sprintf("%s/users/%s", organizationRoleByIdPath(realmId, organizationId, roleId), userId)
}

func (keycloakClient *KeycloakClient) GetOrganizationRoleById(ctx context.Context, realmId, organizationId, roleId string) (*OrganizationGlobalRole, error) {
	var role *OrganizationGlobalRole

	err := keycloakClient.getRoot(ctx, organizationRoleByIdPath(realmId, organizationId, roleId), &role, nil)
	if err != nil {
		return nil, err
	}

	role.RealmId = realmId
	if role.Default {
		role.IsDefault = true
	}

	return role, nil
}

func (keycloakClient *KeycloakClient) NewOrganizationUserRoleAssignment(ctx context.Context, realmId, organizationId, roleId, userId string) error {
	return keycloakClient.putRootWithNoBody(ctx, organizationUserRoleAssignmentPath(realmId, organizationId, roleId, userId))
}

func (keycloakClient *KeycloakClient) GetOrganizationUserRoleAssignment(ctx context.Context, realmId, organizationId, roleId, userId string) (*OrganizationUserRoleAssignment, error) {
	_, err := keycloakClient.getRawRoot(ctx, organizationUserRoleAssignmentPath(realmId, organizationId, roleId, userId), nil)
	if err != nil {
		return nil, err
	}

	return &OrganizationUserRoleAssignment{
		RealmId:        realmId,
		OrganizationId: organizationId,
		RoleId:         roleId,
		UserId:         userId,
	}, nil
}

func (keycloakClient *KeycloakClient) DeleteOrganizationUserRoleAssignment(ctx context.Context, realmId, organizationId, roleId, userId string) error {
	return keycloakClient.deleteRoot(ctx, organizationUserRoleAssignmentPath(realmId, organizationId, roleId, userId), nil)
}
