---
page_title: "keycloak_organization_user_role_assignment Resource"
---

# keycloak\_organization\_user\_role\_assignment Resource

Assigns a global organization role to a user within an organization using `role_id`.

## Example Usage

```hcl
resource "keycloak_organization_user_role_assignment" "assign_manage_load" {
  realm_id        = "my-realm"
  organization_id = "org-id"
  user_id         = "user-id"
  role_id         = "role-id"
}
```

## Argument Reference

- `realm_id` - (Required) The realm ID.
- `organization_id` - (Required) The organization ID.
- `user_id` - (Required) The user ID.
- `role_id` - (Required) The organization global role ID.

## Import

Organization user role assignments can be imported using `{{realm_id}}/{{organization_id}}/{{user_id}}/{{role_id}}`.

Example:

```bash
$ terraform import keycloak_organization_user_role_assignment.assign_manage_load my-realm/org-id/user-id/role-id
```
