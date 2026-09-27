# Vault list-view skill

The `vault-list-view` Bob skill helps generate and modernize Vault UI list views using the shared `Page::ListView` infrastructure. It turns a screenshot or Figma design into a typed `ListViewConfig`, generates the repeated route and view boilerplate, and leaves implementation-specific decisions for the developer to resolve.

The skill is defined in [`.agents/skills/vault-list-view/SKILL.md`](../../.agents/skills/vault-list-view/SKILL.md).

## When to use it

Use the skill when you need to:

- Generate a new Vault UI list or table page from a screenshot.
- Implement a list view from a Figma URL.
- Modernize an existing list route to use `Page::ListView`.
- Create the standard controller, route, template, list-view configuration, and acceptance-test files for a list page.

The skill supports both new routes and existing routes. It automatically determines the intent from the repository state when possible.

## Prerequisites

Before using the skill, have the following information available:

- A screenshot, Figma URL, or equivalent design reference.
- The route module, or enough context for it to be derived from the page and breadcrumbs.
- The API endpoint and response shape when they cannot be determined from the existing route.
- An expanded screenshot or explicit list of options for any `···` actions menu.

The skill does not infer actions menu items, API field names, or destructive API methods from incomplete information. Providing these details explicitly prevents incorrect generated code.

## Example prompt

Attach the list-view screenshot to the prompt. Include the route intent, column behavior, filtering rules, model transformations, and any API methods that should be used. If the screenshot includes an actions (`···`) menu, attach a second screenshot with the menu open or list every menu item explicitly.

```text
Use the vault-list-view skill to modify the existing secrets backends route.

The attached screenshot shows the secrets engines list view. Use it as the source of truth for the page title, breadcrumbs, actions, empty state, and overall table layout.

Configure the table as follows:

- Engine Type column: use `icon-text`. Read the icon from `engine.icon` and the label from `engineDisplayData(engine.engineType).displayName`.
- Engine path column: use a link with `routeKey: "backendLink"` and make the column sortable.
- Accessor column: use `header-tooltip` with this tooltip text: "An accessor is a stable, system-generated identifier used to represent a secret engine path.".
- Description and Version columns: use plain text.
- Promote `icon`, `backendLink`, `backendConfigurationLink`, `version`, `engineType`, and `displayName` to own properties in `model()`. These values are available as getters on `SecretsEngineResource`.
- Filter the list using `engine.shouldIncludeInList`.
- For the delete action, call `this.api.sys.mountsDisableSecretsEngine(item.id)`.

Inspect the existing route, controller, resource type, API response, and acceptance test before making changes. Preserve the existing route registration. Resolve every `TODO(scaffold)` marker, use the exact raw API field names for config keys and filtering, and run the relevant acceptance test when implementation is complete.
```

This prompt works well because it separates visual requirements from implementation requirements. The screenshot supplies the visual structure, while the explicit list supplies the field mappings, resource getters, filtering behavior, and destructive API method that cannot be safely inferred from pixels alone.

## Workflow

The skill follows a config-first workflow:

```text
Screenshot or Figma URL
        ↓
15-step extraction checklist
        ↓
Validated ListViewConfig
        ↓
Review and confirmation
        ↓
Deterministic scaffold generation
        ↓
Resolve TODO(scaffold) markers
        ↓
Static checks and acceptance tests
```

### 1. Provide the design input

Ask Bob to generate or modernize a list view and provide a screenshot or Figma URL. For Figma designs, the skill uses the design context to identify the page structure and controls. A screenshot can be used when Figma context is not available.

The skill collects the route module, parent route block, route segment, and URL path before generating files. Existing route files are inspected automatically when modernizing a page.

### 2. Complete the extraction checklist

The skill records the page requirements in a `ListViewConfig` by checking the design in order. The checklist covers:

- Header title, description, and breadcrumbs.
- Primary and secondary header actions.
- Filter input, filter dropdown, sorting, and bulk actions.
- Selection and nested child rows.
- Table columns, sorting support, and exact raw API field names.
- Declarative column value types, including links, icons, status badges, copy buttons, and date/time values.
- Actions-menu entries and row actions.
- Confirmation modals and empty states.
- Route query parameters, API calls, and client-side filtering.

Column keys must match the raw API field names, normally in `snake_case`. The skill should inspect the route model or API response types rather than inventing aliases.

If the design contains an actions (`···`) column, provide the opened menu or list every option. The skill intentionally stops instead of guessing the menu contents.

### 3. Review and validate the configuration

Before files are generated, review the assembled `ListViewConfig`. The skill validates that:

- The page has a title and valid breadcrumbs.
- The table has a popup-menu column only when row actions are configured.
- Route and action row actions contain their required route or handler data.
- Destructive actions have confirmation-modal configuration.
- The route uses the standard `page` query parameter.
- Nested rows do not use sortable columns, selection, or striped-table behavior.
- The resource key and filter key are present and correspond to the API response.

Resolve validation failures before generating the scaffold.

### 4. Generate the scaffold

After the configuration is confirmed, write it to a temporary JSON file and run the generator from the repository root:

```bash
node .agents/skills/vault-list-view/scripts/generate-scaffold.mjs \
  /tmp/<resource>-list-view-config.json
```

The generator creates the following files at the paths in `routePaths`:

| File | Purpose |
| --- | --- |
| List-view config `.ts` | Defines the title, columns, actions, filter, and empty-state display configuration. |
| Controller `.js` | Imports the config and owns client-side filtering, pagination, and action handlers. |
| Route `.ts` | Loads the API data, exposes the typed model, and sets breadcrumbs. |
| Template `.hbs` | Invokes `Page::ListView` directly. |
| Acceptance test `.js` | Provides config-driven page, table, filter, action, empty-state, and pagination coverage. |

For a new route, register the route in `router.js` manually. The generator does not modify the router. For an existing JavaScript route, the generator creates the TypeScript replacement; delete the old route only after the replacement is complete and verified.

### 5. Resolve generated TODOs

The generated files contain `TODO(scaffold)` markers where judgment or repository-specific information is required. Resolve all markers before committing.

In the list-view config:

- Add the complete `columns` array with exact API field names.
- Add declarative value-type properties for non-text content.
- Add the confirmed `rowActions` and modal configuration.
- Add the filter, primary action, child-row key, and empty-state description when applicable.

In the controller:

- Confirm the config import and resource getter name.
- Confirm the `filterKey` matches the raw API field name.
- Replace any delete API placeholder with the real API method when a destructive action is configured.

In the route:

- Replace the placeholder response type with the actual API item type.
- Add the real API call and response normalization.
- Add the complete breadcrumbs array.
- Stamp `_isChild: true` on nested child rows when `childrenKey` is configured.

In the template and acceptance test:

- Confirm the resource getter and route URL.
- Add stub data for every configured data column.
- Stub the correct API response shape.
- Add assertions for all breadcrumbs, columns, actions, and applicable controls.
- Remove tests for controls that are not present on the page.

## Implementation conventions

The skill follows these list-view conventions:

- Use `Page::ListView` directly; do not create a wrapper page component.
- Keep display behavior in a stable imported `ListViewConfig` object.
- Keep filtering and pagination client-side in the controller getter.
- Use `page` and `pageFilter` as route query parameters. The controller holds only these two tracked properties plus `listViewConfig` — no filtered getter and no action handlers.
- `Page::ListView` receives `@model` (raw array) and computes filtering and pagination internally using `config.filter.filterKey`, `@pageFilter`, and `@page`.
- `Page::ListView` handles the router transition when the filter input changes and owns the try/catch, flash message, and `router.refresh()` for modal deletions. Supply the API call as `deleteAction: (item) => ...` in the modal config.
- Use `.ts` for routes and `.js` for controllers and acceptance tests.
- Use raw API field names for `filterKey` (inside `config.filter`) and column keys.
- Configure popup menus and confirmation modals declaratively.
- For nested rows, omit sortable columns, selection, and striped-table behavior.

## Static checks

Run the static checks after resolving the scaffold markers. The skill checks for common errors including:

- Invalid or unverified HDS icon names.
- Incorrect `@model` and `@models` usage.
- Invalid `concat` expressions in Glimmer attributes.
- Classic `(action ...)` modifiers in Octane templates.
- Inline `@config={{hash ...}}` configuration.
- Incorrect test-file extensions.
- Incorrect filter bindings.
- Missing `@actions={{this}}` when action handlers are configured.
- Invalid modal focus-trap structure.
- Unresolved delete API placeholders.

The full check list is available in [`.agents/skills/vault-list-view/references/static-checks.md`](../../.agents/skills/vault-list-view/references/static-checks.md).

## Testing

Run the generated acceptance test and any related list-view tests after completing the scaffold:

```bash
cd ui
ember test --filter "<ResourceName>"
```

The generated test should verify the page title, breadcrumbs, configured column headers, applicable filter behavior, empty states, row actions, confirmation modals, and pagination.

## Supporting references

- [Skill instructions](../../.agents/skills/vault-list-view/SKILL.md)
- [Configuration schema](../../.agents/skills/vault-list-view/references/config-schema.md)
- [Cardinal rules and pre-flight checks](../../.agents/skills/vault-list-view/references/rules.md)
- [Extraction checklist](../../.agents/skills/vault-list-view/references/checklist.md)
- [Static checks](../../.agents/skills/vault-list-view/references/static-checks.md)
- [Scaffold generator](../../.agents/skills/vault-list-view/scripts/generate-scaffold.mjs)
