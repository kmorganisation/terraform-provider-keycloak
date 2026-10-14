---
page_title: "keycloak_organization_role_group Resource"
---

# keycloak\_organization\_role\_group Resource

Manages a predefined organization role group: a named set of organization global roles that organization admins can assign to their members.

Every organization of the group's `organization_type` sees the group. Organization admins can assign it but not change it; to customise one, they duplicate it into a group owned by their organization. A member holds every role in the group, so changing `role_ids` changes what every member holds, in every organization.

## Example Usage

```hcl
resource "keycloak_organization_role_group" "broker_operations" {
  realm_id          = "my-realm"
  name              = "Operations"
  description       = "Dispatchers and carrier reps"
  organization_type = "broker"
  role_ids = [
    keycloak_organization_global_role.nexus_view_load.role_id,
    keycloak_organization_global_role.nexus_manage_load.role_id,
  ]
}
```

## Argument Reference

- `realm_id` - (Required) The realm this group exists within.
- `name` - (Required) The group name. It must be unique among predefined groups.
- `description` - (Optional) The group description.
- `organization_type` - (Optional) `broker`, `shipper` or `carrier`. If set, only organizations of that type see the group; leave it empty for every type. Changing it forces a new group, which drops every existing member.
- `role_ids` - (Optional) Ids of `keycloak_organization_global_role`s the group grants. A member only holds the roles of apps enabled for their organization.

## Attributes Reference

- `id` - (Computed) The group ID.
- `group_id` - (Computed) The group ID (same value as `id`).

## Import

Organization role groups can be imported using `{{realm_id}}/{{group_id}}`.

```bash
$ terraform import keycloak_organization_role_group.broker_operations my-realm/7e8cf32a-8acb-4d34-89c4-04fb1d10ccad
```
