---
page_title: "keycloak_organization_global_role Resource"
---

# keycloak\_organization\_global\_role Resource

Allows for creating and managing organization global roles in a realm.

Role definitions are realm-scoped and app-scoped.

## Example Usage

```hcl
resource "keycloak_organization_global_role" "nexus_manage_load" {
  realm_id    = "my-realm"
  app         = "nexus"
  name        = "manage-load"
  description = "Can manage loads in Nexus"
  is_default  = false
}
```

## Argument Reference

- `realm_id` - (Required) The realm this organization global role exists within.
- `app` - (Required) The app this role belongs to. Must be one of `admin-portal`, `itracker`, `nexus`, `reporting`, `lighthouse`.
- `name` - (Required) The role name.
- `description` - (Optional) The role description.
- `is_default` - (Optional) Whether this role is default.

## Attributes Reference

- `id` - (Computed) The role ID.
- `role_id` - (Computed) The role ID (same value as `id`).

## Import

Organization global roles can be imported using `{{realm_id}}/{{role_id}}`.

Example:

```bash
$ terraform import keycloak_organization_global_role.nexus_manage_load my-realm/7e8cf32a-8acb-4d34-89c4-04fb1d10ccad
```

## Migration from `keycloak_organization_role`

If you currently manage roles using `keycloak_organization_role`, migrate with manual reimport:

1. Add a `keycloak_organization_global_role` resource.
2. Remove legacy state:
   `terraform state rm keycloak_organization_role.<name>`
3. Import the new resource with role ID:
   `terraform import keycloak_organization_global_role.<name> <realm_id>/<role_id>`
