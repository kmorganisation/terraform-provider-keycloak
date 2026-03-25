---
page_title: "keycloak_organization_global_roles Data Source"
---

# keycloak\_organization\_global\_roles Data Source

Fetches organization global roles for a realm, optionally filtered by app.

Returned roles are sorted deterministically by `app`, then `name`, then `role_id`.

## Example Usage

```hcl
data "keycloak_organization_global_roles" "nexus_roles" {
  realm_id = "my-realm"
  app      = "nexus"
}
```

## Argument Reference

- `realm_id` - (Required) The realm ID.
- `app` - (Optional) Filter by app. Must be one of `admin-portal`, `itracker`, `nexus`, `reporting`, `lighthouse`.

## Attributes Reference

- `roles` - A list of roles with:
  - `role_id` - Role ID.
  - `app` - App name.
  - `name` - Role name.
  - `description` - Role description.
  - `is_default` - Whether the role is default.
