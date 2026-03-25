# Terraform Provider Migration Brief: Keycloak Orgs Global Roles

## Purpose
This document describes required Terraform provider changes after backend role architecture moved from per-organization role definitions to realm-scoped global role definitions.

Use this in the Terraform provider repository as the implementation spec for a Codex session.

## Backend Change Summary

### Old model
- Role definitions were effectively organization-scoped.
- Provider code likely used org-specific role endpoints and role-name based assignment flows.

### New model
- Role definitions are global per realm, app-scoped.
- Role assignments remain organization-scoped.
- Role assignment and invitation flows should be role-ID first.

## New API Endpoints (source of truth)

### Global role definitions
1. `GET /realms/{realm}/orgs/roles?app={optional}`
2. `POST /realms/{realm}/orgs/roles`
3. `GET /realms/{realm}/orgs/roles/{roleId}`
4. `PUT /realms/{realm}/orgs/roles/{roleId}`
5. `DELETE /realms/{realm}/orgs/roles/{roleId}`

### Org role assignment by role ID
1. `GET /realms/{realm}/orgs/{orgId}/roles/id/{roleId}`
2. `PUT /realms/{realm}/orgs/{orgId}/roles/id/{roleId}/users/{userId}`
3. `DELETE /realms/{realm}/orgs/{orgId}/roles/id/{roleId}/users/{userId}`
4. `GET /realms/{realm}/orgs/{orgId}/roles/id/{roleId}/users/{userId}`

### Compatibility endpoint still exists
1. `GET /orgs/{orgId}/roles`
- Returns legacy shape (`id`, `name`, `description`) but filtered by app access.
- Name-based endpoints may return `409` when role name is ambiguous.

## Domain Rules You Must Preserve
1. Role apps: `admin-portal`, `itracker`, `nexus`, `reporting`, `lighthouse`.
2. `admin-portal` roles are always allowed (not gated by org app flags).
3. Assign role must fail for disabled app roles.
4. Revoke role should still work even if app currently disabled.
5. Existing assignments must reappear when app is re-enabled.

## Provider Changes Required

## 1) API Client Layer
- Add global role-definition client methods for CRUD/list by realm and role ID.
- Add role-ID assignment client methods for org/user role binding.
- Keep name-based methods only for backward compatibility if still required by existing resources.

## 2) Resources / Data Sources

### Role definition resource (realm-scoped)
Create or migrate to resource that represents global role definition.

Suggested schema:
1. `realm_id` (required, ForceNew)
2. `role_id` (computed, used as primary API identity)
3. `app` (required)
4. `name` (required)
5. `description` (optional)
6. `is_default` (optional/computed depending on API behavior)

Suggested import formats (support at least one):
1. `{realm_id}/{role_id}`
2. optional helper: `{realm_id}/{app}/{name}` resolving to role ID

### Role definitions data source
- Add optional `app` filter.
- Return stable deterministic ordering for Terraform diff stability.

### Org user role assignment resource
- Switch identity to `role_id` instead of role name.
- If existing resource is name-based, add a new role-ID resource and deprecate old one.

Suggested schema:
1. `realm_id` (required)
2. `organization_id` (required)
3. `user_id` (required)
4. `role_id` (required)

## 3) Backward Compatibility Strategy
- Do not silently map ambiguous names.
- If old name-based path returns `409`, surface actionable error:
  - “Role name is ambiguous. Use role_id-based resource.”
- If provider has legacy role-name resources:
1. Mark deprecated.
2. Keep read behavior if uniquely resolvable.
3. Fail fast on ambiguity.

## 4) State Migration Expectations
- If renaming/changing resource identity, provide migration guidance in docs.
- If provider supports StateUpgrader equivalents, migrate state to role ID where possible.
- Never change server-side role assignment semantics during migration.

## 5) Acceptance Test Plan

### Global roles CRUD
1. Create role definition (`app=nexus`) and verify read/list/update/delete.
2. Validate `is_default` persisted/read correctly.

### Assignment behavior
1. Assign role by role ID succeeds when app enabled.
2. Assign role by role ID fails when app disabled.
3. Revoke still succeeds when app disabled.

### Compatibility behavior
1. Legacy name-based flow works for unique names.
2. Ambiguous names produce `409` and clear provider error.

### Drift and import
1. Import role definition by `{realm}/{role_id}` works.
2. Plan after import is empty (no drift).

## Example HCL (target style)

```hcl
resource "keycloak_org_global_role" "nexus_manage_load" {
  realm_id    = var.realm_id
  app         = "nexus"
  name        = "manage-load"
  description = "Can manage loads in Nexus"
  is_default  = false
}

resource "keycloak_org_user_role_assignment" "assign_manage_load" {
  realm_id        = var.realm_id
  organization_id = var.organization_id
  user_id         = var.user_id
  role_id         = keycloak_org_global_role.nexus_manage_load.role_id
}
```

## Implementation Notes for Codex
1. Discover existing provider resource names and keep naming conventions already used in that repo.
2. Prefer minimal breaking changes, but prioritize role-ID correctness.
3. Keep API retries and error wrapping consistent with provider style.
4. Preserve deterministic set/list ordering to reduce Terraform plan churn.

---

## Initial Prompt for Codex (for Terraform provider repo)

```text
Implement provider support for Keycloak org global roles (realm-scoped, app-scoped) and role-ID based assignment.

Context:
- Backend role definitions moved from org-scoped to global endpoints:
  - GET/POST /realms/{realm}/orgs/roles
  - GET/PUT/DELETE /realms/{realm}/orgs/roles/{roleId}
- Role assignment is role-ID based:
  - GET /realms/{realm}/orgs/{orgId}/roles/id/{roleId}
  - PUT/DELETE/GET /realms/{realm}/orgs/{orgId}/roles/id/{roleId}/users/{userId}
- Legacy /orgs/{orgId}/roles exists for compatibility and may return ambiguous-name 409 in name-based flows.

Goals:
1. Add/upgrade API client methods for global role definitions and role-ID assignment.
2. Add or migrate Terraform resource(s) for global role definitions with fields: realm_id, role_id, app, name, description, is_default.
3. Add or migrate Terraform assignment resource(s) to role_id-based identity.
4. Keep backward compatibility where feasible, but fail clearly on ambiguous name-based operations.
5. Add acceptance tests for CRUD, assignment gating behavior, import, and ambiguous-name error handling.

Required behavior:
- Apps: admin-portal, itracker, nexus, reporting, lighthouse.
- Do not silently resolve ambiguous role names.
- Prefer role-ID as canonical identity.
- Keep deterministic ordering in data sources for stable plans.

Process:
- First inspect current resources/data sources and API client wrappers.
- Then implement in small commits: client layer -> resources/data sources -> tests -> docs.
- Summarize changed files and migration guidance in final output.
```
