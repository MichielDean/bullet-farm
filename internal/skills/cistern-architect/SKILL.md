---
name: cistern-architect
description: Design brief creation for Cistern architect cataractae. Produces a DESIGN_BRIEF.md that constrains implementation to fit the existing codebase's patterns, conventions, and idioms. Use when the architect cataractae needs to investigate the codebase and produce a blueprint for the implementer.
---

# Cistern Architect — Design Brief Protocol

## Purpose

You produce DESIGN_BRIEF.md — a mandatory constraint document that the implementer must follow. The brief closes the gap between "working code" and "code that fits the codebase." Every item in the brief's API Surface Checklist must be verifiably addressed before implementation can pass.

## Investigation Steps

1. Read CONTEXT.md and all revision notes
2. Read the full requirements from the droplet description
3. For each area the requirements touch, investigate the codebase using the checklist below
4. Write DESIGN_BRIEF.md (see format below)
5. Commit (see cistern-git skill — exclude CONTEXT.md)
6. Signal outcome

## Investigation Checklist

Run every item. Missing an item means the brief is incomplete.

### Existing Patterns

| What to find | How to find it | What to write in the brief |
|---|---|---|
| ORM/query patterns | Search for existing column types, query builders, DAO methods in the target package | "Use Exposed's Exists/NotExists DSL — see OrganizationTable.kt:45 for precedent" |
| Naming conventions | Check file names, class names, constant names, column names | "Constants live in OrganizationPermissionNames object — see OrganizationPermissionNames.kt" |
| Error handling | Search for error types, requireNotNull, custom exceptions | "Use requireCatalogPermissionId pattern — see CatalogPermissionCache.kt:23" |
| Collection types | Search for Set vs List vs Map in similar contexts | "loadPermissionsForOrgs returns Map<Long, Map<String, Set<String>>> — Set because UNIQUE constraint means no duplicates" |
| Logging/observability | Search for logging patterns, metrics, tracing | Note any conventions found |

### Reusability and Abstraction Boundaries

| What to find | How to find it | What to write |
|---|---|---|
| Is the feature entity-specific or generic? | Ask: "Could another entity use the same pattern?" | If generic: "PermissionBooleanColumn must accept outerIdColumn as constructor parameter, not hardcode OrganizationTable.id" |
| Existing utilities or base classes | Search for abstract classes, interfaces, extension functions in the target package | "Extend ExistingColumn<Long> — see CustomColumnType.kt" |

### Migration Quality

| What to find | How to find it | What to write |
|---|---|---|
| Migration numbering | List existing migration files, note the next available number | "Use V137…V142 — split DDL and DML into separate migrations" |
| SQL identifier quoting | Read existing migrations — do they use backticks or double quotes? | "SQL identifiers must be backtick-quoted — see V135__add_organization_settings.kt" |
| DDL/DML separation | Check if reference data inserts are in separate migrations from CREATE TABLE | "Separate schema changes (V137) from reference data inserts (V138)" |
| Description quality | Read existing reference data inserts — are descriptions meaningful? | "Permission descriptions must explain the feature, not just repeat the flag name" |

### Repeated Patterns (DRY)

| What to find | How to find it | What to write |
|---|---|---|
| Repeated inline expressions | Search for patterns like `perms[X]?.contains("true") ?: false` — if repeated 3+ times, it's a DRY violation waiting to happen | "Extract boolPerm helper — signature: fun Map<String, Map<String, Set<String>>>.boolPerm(orgId: Long, perm: String): Boolean" |

### Test Coverage

| What to find | How to find it | What to write |
|---|---|---|
| Existing test patterns | Find test files for similar features | "Follow OrganizationDAOSearchTest pattern — see that file for integration test structure" |
| Integration test locations | Find where integration tests live | "Integration tests go in src/test/kotlin/.../integration/ — see existing files" |
| Missing coverage | Identify what the diff touches that needs real DB tests | "New DAO methods for permission search require integration tests in OrganizationDAOSearchTest" |

## Brief Format

```markdown
# Design Brief: <feature title>

## Requirements Summary
<One-paragraph summary of what needs to be built>

## Existing Patterns to Follow

### ORM / Query
<Specific patterns found — name the files and lines>

### Naming Conventions
<File names, class names, column names, constant names, migration numbering>

### Error Handling
<How the codebase handles errors — specific patterns>

### Collection Types
<Where Set vs List vs Map is used and why>

### Migrations
<Numbering, quoting, DDL/DML separation, description quality>

### Testing
<Test patterns, integration test locations, naming conventions>

## Reusability Requirements

<For each new class/utility that applies to more than one entity, specify
that it must accept its context (table, column, ID) as a constructor parameter
rather than hardcoding a reference to a specific entity.>

## DRY Requirements

<Any repeated inline expression pattern that the brief identifies must be
extracted into a helper. Name the helper and specify its signature.>

## Migration Requirements

<Specify: file naming, identifier quoting, DDL/DML separation, description
quality for reference data inserts.>

## Test Requirements

<Specify: which test files need new tests, what kind of tests (unit vs
integration), naming convention for new test functions, and any coverage
gaps to fill.>

## Forbidden Patterns

<Anti-patterns found in the codebase that the brief specifically excludes
from the implementation. Examples: hardcoded entity references in generic
classes, placeholder return values in query builders, mixing constants into
schema definition objects, misleading type names that don't match their
purpose.>

## API Surface Checklist

<Before the implementer can pass, every item in this list must be addressed
in the implementation. Each item is a verification gate.>

- [ ] <specific requirement>
- [ ] <specific requirement>
- [ ] ...
```

## What the Brief Is NOT

- It is NOT a full implementation. Do not write production code.
- It is NOT a test file. Do not write test cases.
- It is NOT a review. Do not review code that does not exist yet.
- It IS a constraint document that the implementer must satisfy.

## Quality Bar

A brief that says "follow existing patterns" without naming specific files and lines is incomplete. Every pattern reference must include a file path. Every constraint must be specific enough that the implementer can verify it with `grep` or by reading the named file.

The brief is complete when:
1. Every checklist item above has been investigated
2. Every pattern reference includes a specific file path
3. The API Surface Checklist covers every user-facing and internally-facing contract
4. There are no "TBD" or "determine during implementation" items