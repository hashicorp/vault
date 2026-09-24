# Accessibility violations — UI test run

> Extracted from 14 CI test result artifacts (QUnit / axe-core 4.11).
> Each entry lists all **unique offending DOM patterns**, the **source files** to fix,
> and the specific **fix required**.

## Summary

| Severity | Unique rules | Affected tests |
|----------|:------------:|:--------------:|
| 🔴 Critical | 4 | 13 |
| 🟠 Serious | 10 | 656 |
| **Total** | **14** | |

---

## 🔴 Critical

### ARIA attributes must conform to valid values

- **Rule:** `aria-valid-attr-value`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/aria-valid-attr-value?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/aria-valid-attr-value?application=axeAPI)
- **Affected tests:** 2

#### Source files and fix

- `ui/app/components/tools/random.hbs`
- `ui/app/components/tools/hash.hbs`
- `ui/app/components/identity/lookup-input.hbs`
- `ui/app/components/form-field-from-model.hbs`
- `ui/app/components/keymgmt/distribute.hbs`
- `ui/app/components/transit-key-action/datakey.hbs`
- `ui/app/components/transit-key-action/decrypt.hbs`
- `ui/app/components/transit-key-action/verify.hbs`
- `ui/app/components/transit-key-action/encrypt.hbs`
- `ui/app/components/transit-key-action/rewrap.hbs`
- `ui/app/components/transit-key-action/sign.hbs`
- `ui/app/components/transit-key-action/hmac.hbs`
- `ui/app/components/transit-key-action/export.hbs`
- `ui/app/components/key-version-select.hbs`
- `ui/app/components/transit-form-edit.hbs`
- `ui/app/components/mount-accessor-select.hbs`

**Fix:** A `<div class="select is-fullwidth">` wrapper has both `aria-controls` and `aria-describedby` pointing to the same ID, which makes `aria-controls` invalid (the referenced element is not a controlled widget). Remove `aria-controls` from the wrapper `<div>` — it should only be on the `<select>` element itself if needed, and `aria-describedby` should reference a description element, not the same element that `aria-controls` points to.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<div class="select is-fullwidth" aria-controls="container-ember53857" aria-describedby="container-ember53857">
```

</details>

<details>
<summary>Affected tests (2)</summary>

- Integration | Component | pki key form: it generates a key type=exported
- Integration | Component | pki key form: it generates a key type=internal

</details>

---

### Buttons must have discernible text

- **Rule:** `button-name`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/button-name?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/button-name?application=axeAPI)
- **Affected tests:** 1

#### Source files and fix

- `ui/app/components/namespace-picker.hbs`

**Fix:** The `<D.ToggleButton>` for the namespace picker renders a button that axe reports as having no discernible text in one scenario (the enterprise landing page dashboard test). `@text={{this.namespace.currentNamespace}}` may be empty/undefined when the current namespace is root (empty string). Add a fallback: `@text={{or this.namespace.currentNamespace "root"}}` so the button always has visible text, or add an `aria-label` fallback via `@aria-label`.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<button class="hds-dropdown-toggle-..." id="toggle-button-ember1..." data-test-button="namespace-picker" aria-expanded="false" type="button" aria-controls="ember1871" popovertarget="ember1871">
```

</details>

<details>
<summary>Affected tests (1)</summary>

- Acceptance | landing page dashboard: hides the configuration details card on a non-root namespace enterprise version

</details>

---

### Certain ARIA roles must contain particular children

- **Rule:** `aria-required-children`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/aria-required-children?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/aria-required-children?application=axeAPI)
- **Affected tests:** 3

#### Source files and fix

- `ui/app/components/namespace-picker.hbs`
- `ui/app/components/clients/date-range.hbs`
- `ui/app/components/clients/filter-toolbar.hbs`

**Fix:** An `<Hds::Dropdown>` list renders with `role="listbox"` but its children (`<D.Description>`, `<D.Generic>`, `<D.Separator>`) are not `role="option"` items. `role="listbox"` requires all direct children to be `role="option"`. This is an HDS component issue: the HDS `Hds::Dropdown` uses `role="listbox"` on the list but places non-option elements (descriptions, separators) inside it. File an upstream HDS issue or, as a workaround, avoid mixing `<D.Description>` or `<D.Separator>` inside dropdowns that use `role="listbox"`. Switch to `role="list"` or use `<D.Interactive>` items only.

<details>
<summary>All unique offending DOM nodes (3)</summary>

```html
<ul class="hds-dropdown__list" role="listbox" aria-labelledby="toggle-button-ember17517">
<ul class="hds-dropdown__list" role="listbox" aria-labelledby="toggle-button-ember2350">
<ul class="hds-dropdown__list" role="listbox" aria-labelledby="toggle-button-ember3108">
```

</details>

<details>
<summary>Affected tests (3)</summary>

- Acceptance | clients | counts > manual refresh: enterprise: it refreshes the client-list route and preserves query params
- Acceptance | clients | counts | client list: it renders error message if export has no data
- Integration | Component | clients/filter-toolbar: it searches and renders no matches found message

</details>

---

### Form elements must have labels

- **Rule:** `label`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/label?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/label?application=axeAPI)
- **Affected tests:** 7

#### Source files and fix

- `ui/app/components/pgp-file.hbs`

**Fix:** `<Hds::Form::Textarea::Field name="pgp-key">` is rendered without a `@label` argument. The HDS `Textarea::Field` component requires either `@label` or a `<F.Label>` block to produce an associated `<label>` element. Add `@label="PGP Key"` (or an appropriate label string) to the `Textarea::Field` invocation. The outer `<label>` on line 8 labels the toggle checkbox, not the textarea.

- `ui/app/components/tools/random.hbs`
- `ui/app/components/tools/hash.hbs`
- `ui/app/components/form-field-from-model.hbs`
- `ui/app/components/mount-accessor-select.hbs`

**Fix:** A plain `<input>` is rendered with an empty `id=""`, making it impossible for a `<label for="">` to associate with it. The input's `id` must be a non-empty unique value that matches a `for` attribute on its `<label>`. Use Ember's `{{unique-id}}` helper or `this.elementId` to generate a stable ID and set it on both the `<label for=...>` and the `<input id=...>`.

<details>
<summary>All unique offending DOM nodes (6)</summary>

```html
<input id="" class="ember-text-field ember-view input" autocomplete="off" spellcheck="false" type="text">
<textarea id="ember17818" class="ember-text-area ember-view hds-form-textarea hds-typography-body-200 hds-font-weight-regular" rows="4" aria-describedby="helper-text-ember17818" name="pgp-key"></textarea>
<textarea id="ember17824" class="ember-text-area ember-view hds-form-textarea hds-typography-body-200 hds-font-weight-regular" rows="4" aria-describedby="helper-text-ember17824" name="pgp-key"></textarea>
<textarea id="ember25063" class="ember-text-area ember-view hds-form-textarea hds-typography-body-200 hds-font-weight-regular" rows="4" aria-describedby="helper-text-ember25063" name="pgp-key"></textarea>
<textarea id="ember25069" class="ember-text-area ember-view hds-form-textarea hds-typography-body-200 hds-font-weight-regular" rows="4" aria-describedby="helper-text-ember25069" name="pgp-key"></textarea>
<textarea id="ember25075" class="ember-text-area ember-view hds-form-textarea hds-typography-body-200 hds-font-weight-regular" rows="4" aria-describedby="helper-text-ember25075" name="pgp-key"></textarea>
```

</details>

<details>
<summary>Affected tests (7)</summary>

- Integration | Component | autocomplete-input: it should render label
- Integration | Component | autocomplete-input: it should trigger dropdown
- Integration | Component | choose-pgp-key-form: it calls cancel on cancel
- Integration | Component | choose-pgp-key-form: it calls onSubmit correctly
- Integration | Component | choose-pgp-key-form: it renders correctly
- Integration | Component | pgp file: it allows for text entry
- Integration | Component | pgp file: toggling back and forth

</details>

---

## 🟠 Serious

### <li> elements must be contained in a <ul> or <ol>

- **Rule:** `listitem`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/listitem?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/listitem?application=axeAPI)
- **Affected tests:** 5

#### Source files and fix

- `ui/app/components/clients/date-range.hbs`
- `ui/app/components/clients/filter-toolbar.hbs`
- `ui/app/components/secret-engine/list.hbs`
- `ui/app/components/namespace-picker.hbs`

**Fix:** `<D.Description>` items rendered by `<Hds::Dropdown>` produce `<li>` elements that appear outside a `<ul>` or `<ol>` in certain layouts. This is the inverse of the `list` rule: the `<li>` items lack a proper list container. This is linked to the same root cause as `aria-required-children` — the HDS Dropdown component's internal DOM structure around description/generic slots. Audit each affected dropdown and either remove `<D.Description>` items or wrap them appropriately. Check for the 'No matching namespaces' and 'Current period' / 'Historical periods' labels specifically.

<details>
<summary>All unique offending DOM nodes (3)</summary>

```html
<li class="hds-dropdown-list-item hds-dropdown-list-item--variant-interactive hds-dropdown-list-item--color-critical">
<li id="ember17535" class="ember-view hds-text hds-typography-body-100 hds-font-weight-regular hds-foreground-faint hds-dropdown-list-item hds-dropdown-list-item--variant-description has-top-padding-xs">
<li id="ember2360" class="ember-view hds-text hds-typography-body-100 hds-font-weight-regular hds-foreground-faint hds-dropdown-list-item hds-dropdown-list-item--variant-description">
```

</details>

<details>
<summary>Affected tests (5)</summary>

- Acceptance | clients | counts > manual refresh: enterprise: it refreshes the client-list route and preserves query params
- Acceptance | clients | counts | client list: it renders error message if export has no data
- Acceptance | raft storage: it should remove raft peer
- Integration | Component | clients/filter-toolbar: it searches and renders no matches found message
- Integration | Component | confirm-action: it renders isInDropdown defaults and calls onConfirmAction

</details>

---

### <svg> elements with an img role must have alternative text

- **Rule:** `svg-img-alt`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/svg-img-alt?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/svg-img-alt?application=axeAPI)
- **Affected tests:** 24

#### Source files and fix

- `ui/app/components/clients/charts/carbon-chart.hbs`
- `ui/app/templates/vault/cluster/clients/config.hbs`

**Fix:** Carbon Charts renders `<path>` elements with `role="graphics-symbol"` and `aria-roledescription="bar"` inside an SVG that itself has `role="img"`. Axe requires that elements with `role="img"` (including SVG child paths with that role) have an accessible name via `aria-label` or `aria-labelledby`. The Carbon Charts library auto-generates these paths — the fix is to ensure the wrapping chart container has `aria-label` set (or `title` inside the SVG). In `carbon-chart.hbs`, pass an `aria-label` via `...attributes` from the call site, or configure the Carbon Chart options to include a chart title accessible name. Alternatively, add `role="presentation"` to the SVG if the chart already has a separate text description.

<details>
<summary>All unique offending DOM nodes (13)</summary>

```html
<path opacity="0.001826131999994478..." class="bar fill-1-3-1" width="20" d="M176.4453125,198.979..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.001826131999994478..." class="bar fill-1-3-1" width="20" d="M81.3984375,198.9799..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.002634011999992950..." class="bar fill-1-3-1" width="20" d="M107.005859375,198.9..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.002634011999992950..." class="bar fill-1-3-1" width="20" d="M445.404296875,198.9..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.01618005214819544" class="bar fill-1-3-1" width="20" d="M107.005859375,198.9..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.01618005214819544" class="bar fill-1-3-1" width="20" d="M445.404296875,198.9..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.06375835185185184" class="bar fill-1-3-1" width="20" d="M292.40625,198.97999..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.2131955094814485" class="bar fill-1-3-1" width="20" d="M123.98883928571428,..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.2131955094814485" class="bar fill-1-3-1" width="20" d="M63.91294642857143,1..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.340736" class="bar fill-1-3-1" width="20" d="M123.98883928571428,..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.340736" class="bar fill-1-3-1" width="20" d="M63.91294642857143,1..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.49900066651875097" class="bar fill-1-3-1" width="20" d="M129.97991071428572,..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
<path opacity="0.49900066651875097" class="bar fill-1-3-1" width="20" d="M70.99330357142857,1..." role="graphics-symbol" aria-roledescription="bar" style="fill: var(--clients-...">
```

</details>

<details>
<summary>Affected tests (24)</summary>

- Acceptance | clients | counts > manual refresh: enterprise: it refreshes the overview route and preserves query params
- Acceptance | clients | counts | client list: enterprise: it hides client list tab on HVD managed clusters
- Acceptance | clients | counts | client list: it navigates to client list tab
- Acceptance | clients | counts: it should redirect to counts overview route for transitions to parent
- Acceptance | clients | overview > static data: it filters attribution table when filters are applied
- Acceptance | clients | overview: it should hide secrets sync stats when feature is NOT on license
- Acceptance | clients | overview: it should render charts
- Integration | Component | clients/page/overview: it filters data if @filterQueryParams specify a month
- Integration | Component | clients/page/overview: it filters data if @filterQueryParams specify a mount_path
- Integration | Component | clients/page/overview: it filters data if @filterQueryParams specify a mount_type
- Integration | Component | clients/page/overview: it filters data if @filterQueryParams specify a multiple filters
- Integration | Component | clients/page/overview: it filters data if @filterQueryParams specify a namespace_path
- Integration | Component | clients/page/overview: it initially renders attribution with by_namespace data
- Integration | Component | clients/page/overview: it renders NEW monthly clients for self-managed clusters instead of total clients
- Integration | Component | clients/page/overview: it renders TOTAL monthly clients for HVD instead of new clients
- Integration | Component | clients/page/overview: it renders dropdown lists from activity response to filter table data
- Integration | Component | clients/page/overview: it renders empty state message when filter selections yield no results
- Integration | Component | clients/page/overview: it shows correct empty state message when selected month has no data
- Integration | Component | clients/running-total: it hides secret sync totals when feature is not activated
- Integration | Component | clients/running-total: it renders text for HVD managed versions
- Integration | Component | clients/running-total: it renders text for ent versions
- Integration | Component | clients/running-total: it renders with full monthly activity data
- Integration | Component | clients/running-total: it text for community versions
- Integration | Component | clients/running-total: it toggles to split chart by client type

</details>

---

### <ul> and <ol> must only directly contain <li>, <script> or <template> elements

- **Rule:** `list`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/list?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/list?application=axeAPI)
- **Affected tests:** 33

#### Source files and fix

- `ui/app/components/mfa/nav.hbs`

**Fix:** `<LinkTo>` components are direct children of a `<ul>`, but only `<li>`, `<script>`, or `<template>` are valid direct children of `<ul>`. Wrap each `<LinkTo>` in an `<li>`: `<li><LinkTo @route="...">Methods</LinkTo></li>`.

- `ui/lib/replication/addon/templates/index.hbs`
- `ui/lib/pki/addon/components/page/pki-configure-create.hbs`

**Fix:** `<h3 class="box-label-header">` is placed inside a `<ul>` as a direct child, which is invalid. Move the heading outside the list or restructure so list items are always `<li>` elements.

- `ui/app/components/tools/wrap.hbs`
- `ui/app/components/tools/random.hbs`
- `ui/app/components/tools/hash.hbs`
- `ui/app/components/tools/rewrap.hbs`
- `ui/app/components/tools/unwrap.hbs`

**Fix:** A bare `<ul>` in the flash/result area contains non-`<li>` direct children. Ensure that every direct child of `<ul>` is an `<li>` element. Check for interpolated text nodes or wrapper `<div>` elements rendered directly inside the list.

<details>
<summary>All unique offending DOM nodes (5)</summary>

```html
<a id="ember1386" class="ember-view active" data-test-tab="methods" href="/ui/vault/access/mfa/methods">
<a id="ember1387" class="ember-view" data-test-tab="enforcements" href="/ui/vault/access/mfa/enforcements">
<ol class="has-left-margin-m has-bottom-margin-s">
<ul class="hds-dropdown__list">
<ul>
```

</details>

<details>
<summary>Affected tests (33)</summary>

- Acceptance |  oidc-config providers and scopes: it creates a scope, and creates a provider with that scope
- Acceptance |  oidc-config providers and scopes: it hides delete and edit for a provider when no permission
- Acceptance |  oidc-config providers and scopes: it lists default provider and navigates to details
- Acceptance |  oidc-config providers and scopes: it navigates to scopes list view and renders empty state when no scopes are configured
- Acceptance |  oidc-config providers and scopes: it renders scope list when scopes exist
- Acceptance |  oidc-config providers and scopes: it throws error when trying to delete when scope is currently being associated with any provider
- Acceptance | /access/identity/entities: it should render correct flash message on entity edit success
- Acceptance | Create groups and entities alias test: entities: it allows create, list, delete of an entity alias
- Acceptance | Create groups and entities alias test: groups: it allows create, list, delete of an entity alias
- Acceptance | Enterprise | replication modes: replication page both primary
- Acceptance | Enterprise | replication modes: replication page when perf primary only
- Acceptance | Enterprise | replication modes: replication page when perf secondary only
- Acceptance | Enterprise | replication navigation: navigate between replication types updates page
- Acceptance | mfa-login-enforcement: it should create login enforcement
- Acceptance | mfa-login-enforcement: it should display login enforcement
- Acceptance | mfa-login-enforcement: it should edit login enforcement
- Acceptance | mfa-login-enforcement: it should list login enforcements
- Acceptance | mfa-login-enforcement: it should send the correct data when creating an enforcement
- Acceptance | mfa-method: it should create method and add it to existing enforcement
- Acceptance | mfa-method: it should create method with new enforcement
- Acceptance | mfa-method: it should create methods
- Acceptance | mfa-method: it should delete method that is not associated with any login enforcements
- Acceptance | mfa-method: it should display method details
- Acceptance | mfa-method: it should edit methods
- Acceptance | mfa-method: it should list methods
- Acceptance | mfa-method: it should navigate to enforcements create route from method enforcement tab
- Acceptance | mfa-method: it should not display for the root namespace
- Acceptance | pki tidy: it configures a manual tidy operation
- Acceptance | pki tidy: it configures a manual tidy operation and shows its details and tidy states
- Acceptance | pki tidy: it configures an auto tidy operation and shows its details
- Acceptance | pki tidy: it opens a tidy modal when the user clicks on the tidy toolbar action
- Acceptance | pki tidy: it should show correct toolbar action depending on whether auto tidy is enabled
- Integration | Component | code-generator/policy/flyout: it yields custom trigger component

</details>

---

### All touch targets must be 24px large, or leave sufficient space

- **Rule:** `target-size`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/target-size?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/target-size?application=axeAPI)
- **Affected tests:** 45

#### Source files and fix

- `ui/lib/core/addon/components/search-select.hbs`

**Fix:** The `ember-basic-dropdown` trigger `<div>` that acts as the `SearchSelect` combobox trigger is below 24×24 px, the minimum touch target size required by WCAG 2.5.5 (AAA) and WCAG 2.5.8 (AA, axe 4.11). The trigger is styled in `ui/app/styles/components/search-select.scss`. Add a `min-height: 36px` (or `min-height: var(--token-component-input-height-medium)`) to the trigger element to match HDS input sizing. Also affects the `hds-copy-button` and `hds-tooltip-button` — these HDS small-size components inherit this constraint and may need a size upgrade from `@size="small"` to `@size="medium"` at call sites, or CSS `min-height`/`padding` overrides in context.

<details>
<summary>All unique offending DOM nodes (44)</summary>

```html
<button class="hds-button hds-button--color-secondary hds-button--is-icon-only hds-button--size-small hds-copy-button hds-button--size-small hds-copy-button--status-idle hds-code-block__copy-button" aria-describedby="pre-code-ember27710" aria-label="Copy" type="button">
<button class="hds-tooltip-button hds-tooltip-button--is-inline" aria-label="More information" data-test-tooltip="search-select" type="button" aria-controls="container-ember29903" aria-describedby="container-ember29903" tabindex="0">
<div id="ember11696-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Filter by type" aria-autocomplete="list" role="combobox" data-ebd-id="ember11697-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember14875-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Enforcement" aria-autocomplete="list" role="combobox" data-ebd-id="ember14876-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember14882-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Entities" aria-autocomplete="list" role="combobox" data-ebd-id="ember14883-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember18304-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Role" aria-autocomplete="list" role="combobox" data-ebd-id="ember18305-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember18539-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Filter by type" aria-autocomplete="list" role="combobox" data-ebd-id="ember18540-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember18782-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Oidc key" aria-autocomplete="list" role="combobox" data-ebd-id="ember18783-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember21520-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Search input role" aria-autocomplete="list" role="combobox" data-ebd-id="ember21521-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." ...>
<div id="ember25787-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Oidc client form key..." aria-autocomplete="list" role="combobox" data-ebd-id="ember25788-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." ...>
<div id="ember25827-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Assignment name" aria-autocomplete="list" role="combobox" data-ebd-id="ember25828-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember2829-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Key name" aria-autocomplete="list" role="combobox" data-ebd-id="ember2830-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember28385-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Application name" aria-autocomplete="list" role="combobox" data-ebd-id="ember28386-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." ...>
<div id="ember29677-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29678-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29684-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29685-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29693-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29694-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29701-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29702-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29708-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29709-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29715-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29716-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29722-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29723-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29729-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29730-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29737-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29738-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29745-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29746-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29768-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29769-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29775-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29776-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29782-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29783-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29793-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29794-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29800-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29801-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29807-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29808-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29820-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29821-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29827-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29828-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29834-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29835-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29841-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29842-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29848-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29849-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29855-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29856-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29862-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29863-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember29869-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="foo" aria-autocomplete="list" role="combobox" data-ebd-id="ember29870-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember39515-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Provider" aria-autocomplete="list" role="combobox" data-ebd-id="ember39516-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember39522-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Key name" aria-autocomplete="list" role="combobox" data-ebd-id="ember39523-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember39533-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Key name" aria-autocomplete="list" role="combobox" data-ebd-id="ember39534-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember39540-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Provider" aria-autocomplete="list" role="combobox" data-ebd-id="ember39541-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember40322-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Role" aria-autocomplete="list" role="combobox" data-ebd-id="ember40323-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember40524-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Methods" aria-autocomplete="list" role="combobox" data-ebd-id="ember40525-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." aria-controls="ember-basic-dropdown...">
<div id="ember40558-trigger" class="ember-view ember-bas..." tabindex="0" aria-label="Existing targets" aria-autocomplete="list" role="combobox" data-ebd-id="ember40559-trigger" aria-expanded="true" aria-disabled="false" aria-owns="ember-basic-dropdown..." ...>
```

</details>

<details>
<summary>Affected tests (45)</summary>

- Acceptance | Enterprise | keymgmt: it should add new key and distribute to provider
- Acceptance | sync | destinations (plural): it should filter destinations list
- Integration | Component | control group success: it unwraps data on submit
- Integration | Component | get-credentials-card: it shows button that can be clicked to credentials route when an item is selected
- Integration | Component | keymgmt/distribute: it does not allow operation selection until valid key/provider combo selected
- Integration | Component | keymgmt/distribute: it hides the key field if passed from the parent
- Integration | Component | keymgmt/distribute: it hides the provider field if passed from the parent
- Integration | Component | keymgmt/distribute: it shows key type select field if new key created
- Integration | Component | kubernetes | Page::Overview: it should show options for SearchSelect
- Integration | Component | ldap | Page::Overview: it should render overview cards
- Integration | Component | mfa-login-enforcement-form: it should add and remove targets
- Integration | Component | mfa-login-enforcement-form: it should save new enforcement
- Integration | Component | mfa-login-enforcement-header: it renders inline
- Integration | Component | mount backend form > auth method: it should render identity_token_key field for WIF engine type
- Integration | Component | oidc/assignment-form: it should save new assignment
- Integration | Component | oidc/client-form: it should save new client
- Integration | Component | oidc/client-form: it should show create assignment modal
- Integration | Component | oidc/key-form: it should update key and limit access to selected applications
- Integration | Component | search select: it adds created item to list items on create and removes without adding back to options on delete
- Integration | Component | search select: it adds discarded list items back into select
- Integration | Component | search select: it behaves correctly if new items not allowed
- Integration | Component | search select: it counts options when wildcard is used and displays the count
- Integration | Component | search select: it does not show name and smaller id for non-identity endpoints
- Integration | Component | search select: it filters options and adds option to create new item when text is entered
- Integration | Component | search select: it moves option from drop down to list when clicked
- Integration | Component | search select: it pre-populates list with passed in selectedOptions
- Integration | Component | search select: it preserves parentManageSelected objects for rendering and removes them from dropdown options
- Integration | Component | search select: it queries multiple models
- Integration | Component | search select: it renders a tooltip beside selection if does not match a record returned from query when passObject=false and idKey=id
- Integration | Component | search select: it renders a tooltip beside selection if does not match a record returned from query when passObject=true and idKey=id
- Integration | Component | search select: it renders correctly when model keys are not standardized
- Integration | Component | search select: it renders ids if model does not have the passed objectKeys as an attribute
- Integration | Component | search select: it renders when passObject=true and model does not have the passed objectKeys as an attr
- Integration | Component | search select: it renders when passed multiple models, passObject=true and one model does not have the attr in objectKeys
- Integration | Component | search select: it renders when passed multiple models, passedObject=false and one model does not have the attr in objectKeys
- Integration | Component | search select: it returns array with objects instead of strings if passObject=true
- Integration | Component | search select: it returns custom object and renders name if passObject=true and multiple objectKeys
- Integration | Component | search select: it returns custom object if passObject=true and multiple objectKeys with objectKeys[0]='id'
- Integration | Component | search select: it shows add suggestion if there are no models
- Integration | Component | search select: it shows both name and smaller id for identity endpoints
- Integration | Component | search select: it shows no results if endpoint 404s
- Integration | Component | search select: it shows options when trigger is clicked
- Integration | Component | search select: it shows passed in options when trigger is clicked
- Integration | Component | search select: it shows selected items not in the returned response and if one model 404s
- Integration | Component | sync | Page::Destinations: it should render toolbar filters and actions

</details>

---

### Elements must meet minimum color contrast ratio thresholds

- **Rule:** `color-contrast`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/color-contrast?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/color-contrast?application=axeAPI)
- **Affected tests:** 487

#### Source files and fix

- `ui/app/components/console/log-help.hbs`
- `ui/app/components/console/ui-panel.hbs`

**Fix:** `<Hds::Link::Inline @color="primary">` renders with the HDS primary link color (`--token-color-foreground-action`). In dark mode this token's resolved value does not meet the 4.5:1 contrast ratio against the panel background. Override the link color in dark mode in `ui/app/styles/theme/dark-mode/_overrides.scss`, or use `@color="secondary"` which has a higher-contrast dark-mode value.

- `ui/app/styles/components/box-label.scss`
- `ui/lib/replication/addon/templates/index.hbs`
- `ui/lib/pki/addon/components/page/pki-configure-create.hbs`

**Fix:** `.box-label-header` uses `color: var(--token-color-palette-neutral-400)` in light mode and `color: var(--token-color-foreground-disabled)` in dark mode (see `_overrides.scss` line 207). Both values fail contrast against the `.box-label` card background. Replace the neutral-400/foreground-disabled colors with a higher-contrast token: `var(--token-color-foreground-strong)` or `var(--token-color-foreground-primary)` in both modes. The dark-mode override in `_overrides.scss` already uses `foreground-strong` for the selected state — apply that same token to the default (unselected) state as well.

- `ui/lib/core/addon/components/replication-action-promote.hbs`
- `ui/lib/core/addon/components/replication-action-update-primary.hbs`
- `ui/lib/replication/addon/components/enable-replication-form.hbs`

**Fix:** `<em class="is-optional">(optional)</em>` uses `color: var(--token-color-palette-neutral-400)` (defined in `ui/app/styles/helper-classes/general.scss` line 76). This color is intentionally muted but fails 4.5:1 contrast in dark mode. Update the `.is-optional` rule to use `var(--token-color-foreground-faint)` which has dark-mode awareness, or switch to the HDS `<Hds::Badge>` / optional indicator pattern that uses accessible tokens by default.

- `ui/app/styles/components/search-select.scss`

**Fix:** The highlighted `ember-power-select-option[aria-current='true']` state uses `background-color: var(--token-color-surface-faint)` with `color: var(--token-color-foreground-primary)`. In dark mode the surface-faint value is close to the foreground-primary value, causing contrast failures on dynamic content rendered inside the option (e.g., KV path spans, role names). Increase the contrast of the highlighted state in dark mode in `_overrides.scss`, or use `var(--token-color-surface-interactive-active)` for the background.

<details>
<summary>All unique offending DOM nodes (23)</summary>

```html
<a target="_blank" rel="noopener noreferrer" class="hds-link-inline hds-link-inline--color-primary hds-link-inline--icon-trailing" href="https://developer.hashicorp.com/vault/docs/command/web"><!---->HashiCorp Developer site<!----></a>
<em class="is-optional">(optional)</em>
<h3 class="box-label-header title is-6">
<li class="ember-power-select-option" id="ember19626-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19633-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19640-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19647-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19654-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19661-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19680-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19687-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19694-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember19701-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41330-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41357-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41384-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41427-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41442-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<li class="ember-power-select-option" id="ember41469-0" aria-selected="false" aria-current="true" data-option-index="0" role="option">
<span>$1</span>
<span>$2</span>
<span>$3</span>
<span>$last</span>
```

</details>

<details>
<summary>Affected tests (487)</summary>

- Acceptance | /access/identity/entities: it renders popup menu for entities
- Acceptance | /access/identity/entities: it renders popup menu for external groups
- Acceptance | /access/identity/entities: it renders popup menu for internal groups
- Acceptance | Azure | configuration > Community > Error handling: it prevents transition and shows api error if config errored on save
- Acceptance | Azure | configuration > Community > create: it should save azure account options
- Acceptance | Azure | configuration > Community > details: it should show configuration with Azure account options configured
- Acceptance | Azure | configuration > Community > edit: it should not save client secret if it has NOT been changed
- Acceptance | Azure | configuration > Community > edit: it should save client secret if it HAS been changed
- Acceptance | Azure | configuration > Enterprise > create: it should NOT transition and show error if config errors
- Acceptance | Azure | configuration > Enterprise > create: it should transition and save issuer if config was not changed but issuer was
- Acceptance | Azure | configuration > Enterprise > create: it should transition and show issuer error if config saved but issuer encountered an error
- Acceptance | Azure | configuration > Enterprise > create: it should transition on save if config was changed but issuer was not
- Acceptance | Azure | configuration > Enterprise > details: it should not show issuer if no WIF configuration data is returned
- Acceptance | Azure | configuration > Enterprise > details: it should save WIF configuration options
- Acceptance | Azure | configuration > Enterprise > edit: it should update WIF attributes
- Acceptance | Azure | configuration: (configurable): it navigates from the list view when NOT configured
- Acceptance | Azure | configuration: (configurable): it navigates from the list view when configured
- Acceptance | Azure | configuration: (configurable): it navigates to the appropriate page when "Exit configuration" is clicked in the plugin settings route
- Acceptance | Azure | configuration: (configurable): it renders tabs when NOT configured
- Acceptance | Azure | configuration: (configurable): it renders tabs when configured
- Acceptance | Azure | configuration: it should transition to configure edit page once engine is mounted
- Acceptance | Enterprise | /access/namespaces: the route should show "delete" option menu for each namespace
- Acceptance | Enterprise | /access/namespaces: the route should switch to the selected namespace on click "Switch to namespace"
- Acceptance | Enterprise | /access/namespaces: the route should update namespace list after create/delete WITH manual refresh in the CLI
- Acceptance | Enterprise | KMIP secrets > kmip role edit form: it submits individually selected operations
- Acceptance | Enterprise | KMIP secrets > kmip role edit form: it submits when operation_all is unchecked
- Acceptance | Enterprise | KMIP secrets > kmip role edit form: it submits when operation_none is toggled off
- Acceptance | Enterprise | KMIP secrets > kmip role edit form: it submits when operation_none is toggled on
- Acceptance | Enterprise | KMIP secrets: it can configure a KMIP secrets engine
- Acceptance | Enterprise | KMIP secrets: it can create a credential
- Acceptance | Enterprise | KMIP secrets: it can create a role
- Acceptance | Enterprise | KMIP secrets: it can create a scope
- Acceptance | Enterprise | KMIP secrets: it can delete a role from the detail page
- Acceptance | Enterprise | KMIP secrets: it can delete a role from the list
- Acceptance | Enterprise | KMIP secrets: it can delete a scope from the list
- Acceptance | Enterprise | KMIP secrets: it can revoke a credential from the generate view
- Acceptance | Enterprise | KMIP secrets: it can revoke a credential from the list
- Acceptance | Enterprise | KMIP secrets: it can revoke from the credentials show page
- Acceptance | Enterprise | KMIP secrets: it navigates to kmip roles view using breadcrumbs
- Acceptance | Enterprise | KMIP secrets: it navigates to kmip scopes view using breadcrumbs
- Acceptance | Enterprise | KMIP secrets: it should enable KMIP & transitions to addon engine route after mount success
- Acceptance | Enterprise | config-ui/message: authenticated it should create, edit, view, and delete a message
- Acceptance | Enterprise | config-ui/message: authenticated it should show multiple messages modal
- Acceptance | Enterprise | config-ui/message: it should clear filter params when switching between tabs
- Acceptance | Enterprise | config-ui/message: it should filter by type and status
- Acceptance | Enterprise | config-ui/message: unauthenticated it should create, edit, view, and delete a message
- Acceptance | Enterprise | config-ui/message: unauthenticated it should show multiple messages modal
- Acceptance | Enterprise | control groups: for v2 secrets it redirects you if you try to navigate to a Control Group restricted path
- Acceptance | Enterprise | control groups: it allows the full flow to work with a saved token
- Acceptance | Enterprise | control groups: it allows the full flow to work without a saved token
- Acceptance | Enterprise | control groups: it displays the warning in the console when making a request to a Control Group path
- Acceptance | Enterprise | keymgmt-configuration-workflow: it navigates to the list-root page when "Exit configuration" is clicked in general-settings
- Acceptance | Enterprise | keymgmt-configuration-workflow: it should display keymgmt configuration and tune keymgmt in the general settings form
- Acceptance | Enterprise | keymgmt-configuration-workflow: it should hide plugin settings tab
- Acceptance | Enterprise | keymgmt: it transitions to list route after mount success
- Acceptance | Enterprise | kv-v2 workflow | edge cases > admin persona: namespace: it can create a secret and new secret version
- Acceptance | Enterprise | kv-v2 workflow | edge cases > admin persona: namespace: it manages state throughout delete, destroy and undelete operations
- Acceptance | Enterprise | namespaces: it clears namespaces when you log out
- Acceptance | Enterprise | namespaces: it displays namespaces whether you log in with a namespace prefixed with / or not
- Acceptance | Enterprise | namespaces: it navigates to the matching namespace when Enter is pressed
- Acceptance | Enterprise | namespaces: it should allow the user to delete a namespace
- Acceptance | Enterprise | oidc auth namespace test: oidc: request is made to auth_url when a namespace is inputted
- Acceptance | Enterprise | replication modes: replication page when disabled
- Acceptance | Enterprise | replication modes: replication page when dr primary only
- Acceptance | Enterprise | replication: DR primary: enables primary and adds secondary
- Acceptance | Enterprise | replication: DR primary: redirects to cluster route when navigating to the secondary details page
- Acceptance | Enterprise | replication: DR primary: runs analytics service when enabled
- Acceptance | Enterprise | replication: DR primary: shows demotion warning when Performance replication is active
- Acceptance | Enterprise | replication: Performance primary: add secondary and delete config
- Acceptance | Enterprise | replication: Performance primary: demotes primary to secondary and displays correct status
- Acceptance | Enterprise | replication: Performance primary: manages secondary token generation and TTL configuration
- Acceptance | Enterprise | replication: Replication Dashboard: displays summary cards for both Performance and DR primaries
- Acceptance | Enterprise | sidebar navigation: it should navigate to Resilience and recovery level > Replication (enterprise)
- Acceptance | GCP | configuration > Community > Error handling: it prevents transition and shows api error if config errored on save
- Acceptance | GCP | configuration > Community > create: it should save gcp account accessType options
- Acceptance | GCP | configuration > Community > details: it should show configuration details with GCP account options configured
- Acceptance | GCP | configuration > Community > edit: it should not save credentials if it has NOT been changed
- Acceptance | GCP | configuration > Community > edit: it should save credentials
- Acceptance | GCP | configuration > Enterprise > Error handling: it shows API error if user previously set credentials but tries to edit the configuration with wif fields
- Acceptance | GCP | configuration > Enterprise: it should show configuration details with WIF options configured
- Acceptance | GCP | configuration: (configurable): it navigates from the list view when NOT configured
- Acceptance | GCP | configuration: (configurable): it navigates from the list view when configured
- Acceptance | GCP | configuration: (configurable): it navigates to the appropriate page when "Exit configuration" is clicked in the plugin settings route
- Acceptance | GCP | configuration: (configurable): it renders tabs when NOT configured
- Acceptance | GCP | configuration: (configurable): it renders tabs when configured
- Acceptance | GCP | configuration: it should transition to configure page on click "Configure" from toolbar
- Acceptance | auth backend list > auth methods are linkable and link to correct view: alicloud auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: approle auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: aws auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: azure auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: cert auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: gcp auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: github auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: jwt auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: kubernetes auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: ldap auth method
- Acceptance | auth backend list > auth methods are linkable and link to correct view: oidc auth method
- Acceptance | auth backend list > enterprise: ent-only auth methods are linkable and link to correct view
- Acceptance | auth backend list > enterprise: token config within namespace
- Acceptance | auth backend list: userpass secret backend
- Acceptance | auth config form > azure: it renders mount fields
- Acceptance | auth config form > azure: it renders tune fields
- Acceptance | auth config form > jwt: it renders mount fields
- Acceptance | auth config form > jwt: it renders tune fields
- Acceptance | auth config form > ldap: it renders mount fields
- Acceptance | auth config form > ldap: it renders tune fields
- Acceptance | auth config form > oidc: it renders mount fields
- Acceptance | auth config form > oidc: it renders tune fields
- Acceptance | auth config form > okta: it renders mount fields
- Acceptance | auth config form > okta: it renders tune fields
- Acceptance | auth login > Enterprise: it sets namespace when renewing token
- Acceptance | auth-methods list view: it filters by auth type
- Acceptance | auth-methods list view: it filters by name
- Acceptance | auth-methods list view: it should disable an auth method
- Acceptance | aws secret backend: aws credentials - type assumed_role
- Acceptance | aws secret backend: aws credentials - type federation_token
- Acceptance | aws secret backend: aws credentials - type iam_user
- Acceptance | aws secret backend: aws credentials - type session_token
- Acceptance | aws secret backend: aws credentials without role read access
- Acceptance | aws | configuration > Community > Error handling: it does not try to save lease configuration if root configuration errored on save
- Acceptance | aws | configuration > Community > Error handling: it prevents transition and shows api error if root config errored on save
- Acceptance | aws | configuration > Community > Error handling: it shows a flash message error and transitions if lease configuration errored on save
- Acceptance | aws | configuration > Community: it does not show access type option and iam fields are shown
- Acceptance | aws | configuration > Enterprise: it saves lease configuration if root configuration was not changed
- Acceptance | aws | configuration > Enterprise: it should not show issuer if no root WIF configuration data is returned
- Acceptance | aws | configuration > Enterprise: it should save root AWS—with IAM options—configuration
- Acceptance | aws | configuration > Enterprise: it should save root AWS—with WIF options—configuration
- Acceptance | aws | configuration > Enterprise: it should show error if old url is entered
- Acceptance | aws | configuration > Enterprise: it should show identity_token_ttl or maxRetries even if they have not been set
- Acceptance | aws | configuration > Enterprise: it should transition to configure page on click "Configure" from toolbar
- Acceptance | aws | configuration > Enterprise: it should update AWS configuration details after editing
- Acceptance | aws | configuration > Enterprise: it shows AWS mount configuration details
- Acceptance | aws | configuration: (configurable): it navigates from the list view when NOT configured
- Acceptance | aws | configuration: (configurable): it navigates from the list view when configured
- Acceptance | aws | configuration: (configurable): it navigates to the appropriate page when "Exit configuration" is clicked in the plugin settings route
- Acceptance | aws | configuration: (configurable): it navigates when NOT configured via dropdown
- Acceptance | aws | configuration: (configurable): it navigates when configured via dropdown
- Acceptance | aws | configuration: (configurable): it renders tabs when NOT configured
- Acceptance | aws | configuration: (configurable): it renders tabs when configured
- Acceptance | billing/overview: should redirect to cluster dashboard when user switches namespace while on billing/overview route on enterprise
- Acceptance | chroot-namespace enterprise ui: a user with default policy should see nav items
- Acceptance | chroot-namespace enterprise ui: a user with read policy should see nav items
- Acceptance | chroot-namespace enterprise ui: it works within a child namespace
- Acceptance | chroot-namespace enterprise ui: root-only nav items are unavailable
- Acceptance | cluster: enterprise nav item links to first route that user has access to
- Acceptance | cluster: hides nav item if user does not have permission
- Acceptance | cluster: it hides mfa setup if user does not have entityId (ex: is a root user)
- Acceptance | cluster: redirects to /secrets-engines/kv/kv/list from legacy /secrets/kv/kv/list path
- Acceptance | cluster: shows error banner if resultant-acl check fails
- Acceptance | console: array output is correctly formatted
- Acceptance | console: boolean output is correctly formatted
- Acceptance | console: fullscreen command expands the cli panel
- Acceptance | console: it should open and close console panel
- Acceptance | console: number output is correctly formatted
- Acceptance | console: refresh reloads the current route's data
- Acceptance | database workflow > connections: create connection with rotate failure
- Acceptance | database workflow > connections: create failure
- Acceptance | database workflow > connections: create with rotate
- Acceptance | database workflow > connections: create without rotate
- Acceptance | database workflow > dynamic roles: it creates a dynamic role attached to the current connection
- Acceptance | database workflow > static roles: set parent db to not rotate static roles immediately, verify static role reflects that default
- Acceptance | database workflow > static roles: set parent db to rotate static roles immediately, verify static role reflects that default
- Acceptance | enterprise vault-reporting: it hides the nav item if policy does not allow access to sys/utilization-report
- Acceptance | enterprise vault-reporting: it visits the usage reporting dashboard and renders the header
- Acceptance | enterprise | pki | external | acme-accounts route: it displays 403 permission denied error for list response
- Acceptance | enterprise | pki | external | acme-accounts route: it displays 500 internal server error
- Acceptance | enterprise | pki | external | acme-accounts route: it fetches ACME account details for each account
- Acceptance | enterprise | pki | external | acme-accounts route: it handles a 404
- Acceptance | enterprise | pki | external | acme-accounts route: it handles partial failures when reading individual accounts
- Acceptance | enterprise | pki | external | acme-accounts route: it navigates to acme-accounts route
- Acceptance | enterprise | pki | external | dns-providers route: it displays 403 permission denied error
- Acceptance | enterprise | pki | external | dns-providers route: it displays 500 internal server error
- Acceptance | enterprise | pki | external | dns-providers route: it fetches DNS provider details for each provider type
- Acceptance | enterprise | pki | external | dns-providers route: it handles a 404
- Acceptance | enterprise | pki | external | dns-providers route: it handles partial failures when reading individual providers
- Acceptance | enterprise | pki | external | dns-providers route: it navigates to dns-providers route
- Acceptance | enterprise | pki | external | dns-providers route: it throws error for unsupported DNS provider type
- Acceptance | enterprise | pki | external | orders route: it handles 403 permission denied error
- Acceptance | enterprise | pki | external | orders route: it handles 500 internal server error
- Acceptance | enterprise | pki | external | orders route: it handles empty orders list (404)
- Acceptance | enterprise | pki | external | orders route: it navigates to recent orders
- Acceptance | enterprise | pki | external | orders route: it respects provided query param
- Acceptance | enterprise | pki | external | orders route: it sets default query param when not provided
- Acceptance | enterprise | pki | external | orders route: it updates query param when time period selected from dropdown
- Acceptance | enterprise | pki | external | overview route: it catches 403 permissions errors and hides cards
- Acceptance | enterprise | pki | external | overview route: it catches 404 errors
- Acceptance | enterprise | pki | external | overview route: it catches and displays non-404/non-403 error messages
- Acceptance | enterprise | pki | external | overview route: it navigates to external overview
- Acceptance | enterprise | pki | external | overview route: only "Overview" tab renders when no resources exist but user has permission to list everything
- Acceptance | enterprise | pki | external | roles | index route: it displays 403 permission denied error
- Acceptance | enterprise | pki | external | roles | index route: it displays 500 internal server error
- Acceptance | enterprise | pki | external | roles | index route: it handles a 404
- Acceptance | enterprise | pki | external | roles | index route: it navigates to roles index
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it fetches and displays active orders
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it handles 403 permission denied error
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it handles 500 internal server error
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it handles empty orders list (404)
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it navigates to individual order details
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it redirects to parent error route if role read 403s
- Acceptance | enterprise | pki | external | roles | role | active-orders route: it renders breadcrumbs for role active orders
- Acceptance | enterprise | pki | external | roles | role | order route: it catches cert fetch 403 error
- Acceptance | enterprise | pki | external | roles | role | order route: it catches cert fetch 404 error
- Acceptance | enterprise | pki | external | roles | role | order route: it catches order status 403 error
- Acceptance | enterprise | pki | external | roles | role | order route: it redirects to parent error route if role read 403s
- Acceptance | enterprise | pki | external | roles | role | order route: it renders breadcrumbs and header without tabs for role order
- Acceptance | enterprise | pki | external | roles | role | order route: it requests order details
- Acceptance | enterprise | pki | external | roles | role | order route: it throws if both endpoints error
- Acceptance | enterprise | pki | external | roles | role | order route: it throws order status 404
- Acceptance | enterprise | pki | external | roles | role | overview route: it catches active-order 403 error
- Acceptance | enterprise | pki | external | roles | role | overview route: it catches active-order 404 error
- Acceptance | enterprise | pki | external | roles | role | overview route: it fetches and displays role details
- Acceptance | enterprise | pki | external | roles | role | overview route: it handles role read 403 permission denied error
- Acceptance | enterprise | pki | external | roles | role | overview route: it handles role read 404 error
- Acceptance | enterprise | pki | external | roles | role | overview route: it navigates to role overview route
- Acceptance | kv-v2 workflow | delete, undelete, destroy > admin persona: can delete and undelete the latest secret version (a)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > admin persona: can destroy a secret version (a)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > admin persona: can permanently delete all secret versions (a)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > admin persona: can soft delete and undelete an older secret version (a)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-list-reader persona: can delete and cannot undelete the latest secret version (dlr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-list-reader persona: can soft delete and undelete an older secret version (dlr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-list-reader persona: cannot destroy a secret version (dlr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-list-reader persona: cannot permanently delete all secret versions (dlr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-reader persona: cannot delete and undelete the latest secret version (dr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-reader persona: cannot destroy a secret version (dr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-reader persona: cannot permanently delete all secret versions (dr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > data-reader persona: cannot soft delete and undelete an older secret version (dr)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > metadata-maintainer persona: can destroy a secret version (mm)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > metadata-maintainer persona: can soft delete and undelete an older secret version (mm)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > metadata-maintainer persona: cannot delete but can undelete the latest secret version (mm)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > metadata-maintainer persona: cannot permanently delete all secret versions (mm)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-creator persona: can permanently delete all secret versions (sc)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-creator persona: cannot delete and undelete the latest secret version (sc)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-creator persona: cannot destroy a secret version (sc)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-creator persona: cannot soft delete and undelete an older secret version (sc)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-nested-creator persona: can delete all secret versions from the nested list view (snc)
- Acceptance | kv-v2 workflow | delete, undelete, destroy > secret-nested-creator persona: can not delete all secret versions from root list view (snc)
- Acceptance | kv-v2 workflow | edge cases > destruction without read: it hides destroy option without version number
- Acceptance | kv-v2 workflow | edge cases > destruction without read: it renders the delete action and disables delete latest version option
- Acceptance | kv-v2 workflow | edge cases > destruction without read: it renders the delete action and disables delete this version option
- Acceptance | kv-v2 workflow | edge cases > destruction without read: it renders the destroy metadata action and expected modal copy
- Acceptance | kv-v2 workflow | edge cases > patch-persona: it adds and deletes a key
- Acceptance | kv-v2 workflow | edge cases > patch-persona: it patches a secret from the overview page
- Acceptance | kv-v2 workflow | edge cases > patch-persona: it patches a secret from the secret details
- Acceptance | kv-v2 workflow | edge cases > patch-persona: patching a destroyed secret is not allowed
- Acceptance | kv-v2 workflow | edge cases > persona with glob (*) read and list access on the secret level: it can navigate to secrets within a secret directory
- Acceptance | kv-v2 workflow | edge cases > persona with glob (*) read and list access on the secret level: it handles errors when attempting to view details of a secret that is a directory
- Acceptance | kv-v2 workflow | edge cases > persona with glob (*) read and list access on the secret level: it navigates back to engine index route via breadcrumbs from secret details
- Acceptance | kv-v2 workflow | edge cases > persona with glob (*) read and list access on the secret level: it navigates to secret if policy check fails for the subkeys endpoint
- Acceptance | kv-v2 workflow | edge cases > persona with list access on a secret path: it lists secrets within the root directory from the kv engine list
- Acceptance | kv-v2 workflow | edge cases > persona with list access on a secret path: it lists secrets within the root directory from the quick actions card
- Acceptance | kv-v2 workflow | edge cases: advanced secret values default to JSON display
- Acceptance | kv-v2 workflow | edge cases: does not register as advanced when value includes {
- Acceptance | kv-v2 workflow | edge cases: no ghost item after editing metadata
- Acceptance | kv-v2 workflow | edge cases: viewing advanced secret data versions displays the correct version data
- Acceptance | kv-v2 workflow | navigation > admin persona: breadcrumbs, tabs & page titles are correct (a)
- Acceptance | kv-v2 workflow | navigation > admin persona: can access nested secret (a)
- Acceptance | kv-v2 workflow | navigation > admin persona: empty backend - breadcrumbs, title, tabs, emptyState (a)
- Acceptance | kv-v2 workflow | navigation > admin persona: enterprise: patch route does not redirect for users with permissions (a)
- Acceptance | kv-v2 workflow | navigation > admin persona: it redirects from LIST, SHOW and EDIT views using old non-engine url to ember engine url (a)
- Acceptance | kv-v2 workflow | navigation > admin persona: versioned secret nav, tabs (a)
- Acceptance | kv-v2 workflow | navigation > data-list-reader persona: breadcrumbs & page titles are correct (dlr)
- Acceptance | kv-v2 workflow | navigation > data-list-reader persona: can access nested secret (dlr)
- Acceptance | kv-v2 workflow | navigation > data-list-reader persona: empty backend - breadcrumbs, title, tabs, emptyState (dlr)
- Acceptance | kv-v2 workflow | navigation > data-list-reader persona: enterprise: patch route redirects for users without permissions (dlr)
- Acceptance | kv-v2 workflow | navigation > data-list-reader persona: versioned secret nav, tabs, breadcrumbs (dlr)
- Acceptance | kv-v2 workflow | navigation > data-reader persona: breadcrumbs & page titles are correct (dr)
- Acceptance | kv-v2 workflow | navigation > data-reader persona: can access nested secret (dr)
- Acceptance | kv-v2 workflow | navigation > data-reader persona: enterprise: patch route redirects for users without permissions (dr)
- Acceptance | kv-v2 workflow | navigation > data-reader persona: versioned secret nav, tabs, breadcrumbs (dr)
- Acceptance | kv-v2 workflow | navigation > enterprise controlled access persona: breadcrumbs & page titles are correct (cg)
- Acceptance | kv-v2 workflow | navigation > enterprise controlled access persona: can access nested secret (cg)
- Acceptance | kv-v2 workflow | navigation > enterprise controlled access persona: can patch a secret (cg)
- Acceptance | kv-v2 workflow | navigation > enterprise controlled access persona: can read custom_metadata from data endpoint (cg)
- Acceptance | kv-v2 workflow | navigation > enterprise controlled access persona: can request custom_metadata from data endpoint (cg)
- Acceptance | kv-v2 workflow | navigation > metadata-maintainer persona: breadcrumbs & page titles are correct (mm)
- Acceptance | kv-v2 workflow | navigation > metadata-maintainer persona: can access nested secret (mm)
- Acceptance | kv-v2 workflow | navigation > metadata-maintainer persona: empty backend - breadcrumbs, title, tabs, emptyState (mm)
- Acceptance | kv-v2 workflow | navigation > metadata-maintainer persona: enterprise: patch route redirects for users without permissions (mm)
- Acceptance | kv-v2 workflow | navigation > metadata-maintainer persona: versioned secret nav, tabs, breadcrumbs (mm)
- Acceptance | kv-v2 workflow | navigation > patch-persona: it does not redirect for ent
- Acceptance | kv-v2 workflow | navigation > patch-persona: it navigates to patch a secret from overview
- Acceptance | kv-v2 workflow | navigation > patch-persona: it redirects for community edition
- Acceptance | kv-v2 workflow | navigation > patch-persona: overview subkeys card is hidden for community edition
- Acceptance | kv-v2 workflow | navigation > secret-creator persona: breadcrumbs & page titles are correct (sc)
- Acceptance | kv-v2 workflow | navigation > secret-creator persona: can access nested secret (sc)
- Acceptance | kv-v2 workflow | navigation > secret-creator persona: empty backend - breadcrumbs, title, tabs, emptyState (sc)
- Acceptance | kv-v2 workflow | navigation > secret-creator persona: enterprise: patch route redirects for users without permissions (sc)
- Acceptance | kv-v2 workflow | navigation > secret-creator persona: versioned secret nav, tabs, breadcrumbs (sc)
- Acceptance | kv-v2 workflow | navigation: KVv2 handles nested secret with % and space in path correctly
- Acceptance | kv-v2 workflow | navigation: KVv2 handles nested secret with a percent-encoded data octet in path correctly
- Acceptance | kv-v2 workflow | navigation: KVv2 handles secret with % and space in path correctly
- Acceptance | kv-v2 workflow | navigation: enterprise: it renders policy generator on each page header
- Acceptance | kv-v2 workflow | navigation: it does not render policy generator on community
- Acceptance | kv-v2 workflow | version history, paths > admin persona: can navigate to the paths page (a)
- Acceptance | kv-v2 workflow | version history, paths > admin persona: can navigate to the version history page (a)
- Acceptance | kv-v2 workflow | version history, paths > data-list-reader persona: can navigate to the paths page (dlr)
- Acceptance | kv-v2 workflow | version history, paths > data-list-reader persona: cannot navigate to the version history page (dlr)
- Acceptance | kv-v2 workflow | version history, paths > data-reader persona: can navigate to the paths page (dr)
- Acceptance | kv-v2 workflow | version history, paths > data-reader persona: cannot navigate to the version history page (dr)
- Acceptance | kv-v2 workflow | version history, paths > enterprise controlled access persona: can navigate to the paths page (cg)
- Acceptance | kv-v2 workflow | version history, paths > enterprise controlled access persona: can navigate to the version history page (cg)
- Acceptance | kv-v2 workflow | version history, paths > metadata-maintainer persona: can navigate to the paths page (mm)
- Acceptance | kv-v2 workflow | version history, paths > metadata-maintainer persona: can navigate to the version history page (mm)
- Acceptance | kv-v2 workflow | version history, paths > secret-creator persona: can navigate to the paths page (sc)
- Acceptance | kv-v2 workflow | version history, paths > secret-creator persona: cannot navigate to the version history page (sc)
- Acceptance | ldap | libraries: it should show libraries on overview page
- Acceptance | ldap | libraries: it should transition to create library route on toolbar link click
- Acceptance | ldap | libraries: it should transition to details routes from tab links
- Acceptance | ldap | libraries: it should transition to library details for hierarchical list items
- Acceptance | ldap | libraries: it should transition to library details route on list item click
- Acceptance | ldap | libraries: it should transition to routes from library details page header dropdown
- Acceptance | ldap | libraries: it should transition to routes from list item action menu
- Acceptance | ldap | overview: (configurable): it navigates from the list view when NOT configured
- Acceptance | ldap | overview: (configurable): it navigates from the list view when configured
- Acceptance | ldap | overview: (configurable): it navigates to the appropriate page when "Exit configuration" is clicked in the plugin settings route
- Acceptance | ldap | overview: (configurable): it navigates when NOT configured via dropdown
- Acceptance | ldap | overview: (configurable): it navigates when configured via dropdown
- Acceptance | ldap | overview: (configurable): it renders tabs when NOT configured
- Acceptance | ldap | overview: (configurable): it renders tabs when configured
- Acceptance | ldap | overview: it should delete the ldap engine on delete action
- Acceptance | ldap | overview: it should transition to configuration edit on empty state click
- Acceptance | ldap | overview: it should transition to configuration route when engine is not configured
- Acceptance | ldap | overview: it should transition to create library route on card action link click
- Acceptance | ldap | overview: it should transition to create role route on card action link click
- Acceptance | ldap | overview: it should transition to ldap overview on mount success
- Acceptance | ldap | overview: it should transition to role credentials route on generate credentials action
- Acceptance | ldap | overview: it should transition to routes on tab link click
- Acceptance | mfa-setup: it closes the dropdown after navigating
- Acceptance | mfa-setup: it should login through MFA and post to generate and be able to restart the setup
- Acceptance | mfa-setup: it should show a warning if you enter in the same UUID without restarting the setup
- Acceptance | oidc provider: OIDC Provider logs in and redirects correctly
- Acceptance | oidc provider: OIDC Provider redirects if authorization request throws a permission denied error
- Acceptance | oidc provider: OIDC Provider redirects to auth if current token and prompt = login
- Acceptance | oidc provider: OIDC Provider shows consent form when prompt = consent
- Acceptance | oidc provider: prompt=none with no session and registered redirect_uri: after login the backend redirects to the registered URI
- Acceptance | oidc provider: prompt=none with no session redirects to Vault auth, not directly to redirect_uri
- Acceptance | pki action forms test > generate CSR: happy path
- Acceptance | pki action forms test > generate CSR: type = exported
- Acceptance | pki action forms test > generate root: happy path
- Acceptance | pki action forms test > generate root: type=exported
- Acceptance | pki action forms test > import: happy path
- Acceptance | pki action forms test > import: shows None for imported items if nothing new imported
- Acceptance | pki action forms test > import: shows imported items when keys is empty
- Acceptance | pki action forms test > import: with many imports
- Acceptance | pki overview: hides roles and certificates card if user does not have permissions
- Acceptance | pki overview: navigates to certificate details page for View Certificates card
- Acceptance | pki overview: navigates to generate certificate page for Issue Certificates card
- Acceptance | pki overview: navigates to issuer details page for View Issuer card
- Acceptance | pki overview: navigates to view issuers when link is clicked on issuer card
- Acceptance | pki overview: navigates to view roles when link is clicked on roles card
- Acceptance | pki workflow > config: it updates config when user only has permission to some endpoints
- Acceptance | pki workflow > issuers: details view renders correct number of info items
- Acceptance | pki workflow > issuers: issuer list items link to the correct details route when issuer count exceeds 10
- Acceptance | pki workflow > issuers: issuer list page loads and renders linked items when issuer count exceeds 10
- Acceptance | pki workflow > issuers: lists the correct issuer metadata info
- Acceptance | pki workflow > issuers: lists the correct issuer metadata info when user has only read permission
- Acceptance | pki workflow > issuers: toolbar links navigate to expected routes
- Acceptance | pki workflow > keys: it hides correct actions for user with read policy
- Acceptance | pki workflow > keys: it shows correct toolbar items for the user with update policy
- Acceptance | pki workflow > keys: shows correct items if user has all permissions
- Acceptance | pki workflow > not configured: empty state messages are correct when PKI not configured
- Acceptance | pki workflow > roles: create role happy path
- Acceptance | pki workflow > roles: it does not show toolbar items the user does not have permission to see
- Acceptance | pki workflow > roles: it navigates between tabs when user only has permission to read roles
- Acceptance | pki workflow > roles: it shows correct toolbar items for the user policy
- Acceptance | pki workflow > roles: shows correct items if user has all permissions
- Acceptance | pki workflow > rotate: it renders a warning banner when parent issuer has unsupported OIDs
- Acceptance | pki/pki cross sign: it cross-signs an issuer
- Acceptance | reduced disclosure test > enterprise: it works for user accessing child namespace
- Acceptance | reduced disclosure test > enterprise: login works when reduced disclosure enabled (ent)
- Acceptance | reduced disclosure test: login works when reduced disclosure enabled
- Acceptance | reset password: allows password reset for userpass users logged in via dropdown
- Acceptance | reset password: does not allow password reset for non-userpass users
- Acceptance | reset password: renders error if auth data is unavailable
- Acceptance | reset password: renders error template when user lacks update permission
- Acceptance | secrets-engines/enable > WIF secret engines: it sets identity_token_key on mount config using search select list, resets after
- Acceptance | secrets-engines/enable: enable alicloud
- Acceptance | secrets-engines/enable: enable gcpkms
- Acceptance | secrets-engines/enable: it should transition to different locations for kv v1 and v2
- Acceptance | secrets-engines/enable: it should transition to general settings configuration page for unsupported backends
- Acceptance | secrets-engines/enable: it should transition to mountable addon engine after mount success
- Acceptance | secrets-engines/enable: it should transition to mountable non-addon engine after mount success
- Acceptance | secrets/database/*: Can create and delete a connection
- Acceptance | secrets/database/*: Role create form
- Acceptance | secrets/database/*: buttons show up for managing connection
- Acceptance | secrets/database/*: can enable the database secrets engine
- Acceptance | secrets/database/*: connection_url is decoded
- Acceptance | secrets/database/*: database connection create and edit: elasticsearch-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mongodb-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mssql-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mysql-aurora-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mysql-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mysql-legacy-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: mysql-rds-database-plugin
- Acceptance | secrets/database/*: database connection create and edit: postgresql-database-plugin
- Acceptance | secrets/database/*: database connection create: vault-plugin-database-oracle
- Acceptance | secrets/database/*: database connection edit: vault-plugin-database-oracle
- Acceptance | secrets/database/*: root and limited access
- Acceptance | secrets/secret/create, read, delete > kv v1: KVv1 handles secret with % in path correctly
- Acceptance | secrets/secret/create, read, delete > kv v1: creating a secret with a single or double quote works properly
- Acceptance | secrets/secret/create, read, delete > kv v1: filter clears on nav
- Acceptance | secrets/secret/create, read, delete > kv v1: first level secrets redirect properly upon deletion
- Acceptance | secrets/secret/create, read, delete > kv v1: it can edit via the JSON input
- Acceptance | secrets/secret/create, read, delete > kv v1: paths are properly encoded
- Acceptance | secrets/secret/create, read, delete > kv v1: version 1 performs the correct capabilities lookup
- Acceptance | secrets/secret/create, read, delete > kv v1: version 1 token without read permissions can create and update a secret
- Acceptance | secrets/secret/create, read, delete > kv v1: version 1: nested paths creation maintains ability to navigate the tree
- Acceptance | secrets/secret/create, read, delete > kv v2: it can create a secret when check-and-set is required
- Acceptance | secrets/secret/create, read, delete > kv v2: it navigates to version history and to a specific version
- Acceptance | secrets/secret/create, read, delete > mount and configure: it can mount a KV 2 secret engine with config metadata
- Acceptance | secrets/secret/create, read, delete > mount and configure: v1 key named keys
- Acceptance | settings/auth/configure > configure route does not require read on per-mount sys/auth/* path: it loads the configure route without a 403
- Acceptance | settings/auth/configure > configure route does not require read on per-mount sys/auth/* path: it routes to the error page for a non-existent auth mount path
- Acceptance | settings/auth/configure/section: it shows tabs for auth method: aws
- Acceptance | settings/auth/configure/section: it shows tabs for auth method: azure
- Acceptance | settings/auth/configure/section: it shows tabs for auth method: gcp
- Acceptance | settings/auth/configure/section: it shows tabs for auth method: github
- Acceptance | settings/auth/configure/section: it shows tabs for auth method: kubernetes
- Acceptance | settings/auth/enable: it mounts and redirects
- Acceptance | settings/auth/enable: it renders default config details
- Acceptance | settings/auth/enable: it renders direct login link for supported method
- Acceptance | sidebar navigation: collapsed sidebar does not overlap web repl when console is open
- Acceptance | ssh | configuration: it displays error if generate Signing key is not checked and no public and private keys
- Acceptance | ssh | configuration: it should show a public key after saving default configuration and allows you to delete public key
- Acceptance | ssh | configuration: it should show error if old url is entered
- Acceptance | ssh | roles > Acceptance | ssh | otp role: it deletes a role from list view
- Acceptance | ssh | roles > Acceptance | ssh | otp role: it generates an OTP
- Acceptance | ssh | roles: it creates roles, generates keys and deletes roles
- Acceptance | sync | overview > when feature is not activated > enterprise with namespaces: it should make activation-flag requests to correct namespace
- Acceptance | sync | overview > when feature is not activated > enterprise with namespaces: it should make activation-flag requests to correct namespace when managed
- Acceptance | transit: create form renders supported options for each key type
- Acceptance | transit: it generates a key
- Acceptance | transit: it rotates, encrypts and decrypts key type chacha20-poly1305
- Acceptance | transit: transit backend: aes128-gcm96
- Acceptance | transit: transit backend: aes128-gcm96 
- Acceptance | transit: transit backend: aes256-gcm96
- Acceptance | transit: transit backend: aes256-gcm96 
- Acceptance | transit: transit backend: chacha20-poly1305
- Acceptance | transit: transit backend: chacha20-poly1305 
- Acceptance | transit: transit backend: ecdsa-p256
- Acceptance | transit: transit backend: ecdsa-p384
- Acceptance | transit: transit backend: ecdsa-p521
- Acceptance | transit: transit backend: ed25519
- Acceptance | transit: transit backend: rsa-2048
- Acceptance | transit: transit backend: rsa-3072
- Acceptance | transit: transit backend: rsa-4096
- Acceptance | wrapped_token query param functionality: it authenticates when used with the with=token query param
- Acceptance | wrapped_token query param functionality: it authenticates you if the query param is present
- Acceptance | wrapped_token query param functionality: it makes request to authentication service with expected args
- Acceptance | wrapped_token query param functionality: it should authenticate when hitting logout url with wrapped_token when logged out
- Acceptance | wrapped_token query param functionality: it shows error if unwrap fails and goes back to login form
- Integration | Component | console/ui panel: it renders
- Integration | Component | enable-replication-form > only DR replication in features: attempting to enable performance replication
- Integration | Component | enable-replication-form > shows API errors: dr primary
- Integration | Component | enable-replication-form > shows API errors: dr secondary
- Integration | Component | enable-replication-form > shows API errors: performance primary
- Integration | Component | enable-replication-form > shows API errors: performance secondary
- Integration | Component | enable-replication-form > successful enable: dr primary
- Integration | Component | enable-replication-form > successful enable: dr secondary
- Integration | Component | enable-replication-form > successful enable: performance primary
- Integration | Component | enable-replication-form > successful enable: performance secondary
- Integration | Component | enable-replication-form: enable DR when cluster is perf primary
- Integration | Component | enable-replication-form: it renders correct form inputs when dr replication mode
- Integration | Component | enable-replication-form: it renders correct form inputs when performance replication mode
- Integration | Component | page/pki-configure-create: it renders
- Integration | Component | replication page/mode-index > DR mode: it renders correctly when replication disabled
- Integration | Component | replication page/mode-index > DR mode: it shows enable button if has permissions
- Integration | Component | replication page/mode-index > Performance mode: it renders correctly when replication disabled
- Integration | Component | replication page/mode-index > Performance mode: it shows enable button if has permissions
- Integration | Component | sidebar-frame: it should render logo and actions in app header
- Integration | Component | suggestion-input > Database type: it should fetch database static roles for initial mount
- Integration | Component | suggestion-input > Database type: it should filter current result set
- Integration | Component | suggestion-input > Database type: it should only render dropdown when suggestions exist
- Integration | Component | suggestion-input > Database type: it should set selected role as value
- Integration | Component | suggestion-input > KV type: it should fetch secrets and update suggestions on mountPath change
- Integration | Component | suggestion-input > KV type: it should fetch secrets at nested paths
- Integration | Component | suggestion-input > KV type: it should fetch suggestions for initial mount path
- Integration | Component | suggestion-input > KV type: it should filter current result set
- Integration | Component | suggestion-input > KV type: it should only render dropdown when suggestions exist
- Integration | Component | suggestion-input > KV type: it should replace filter terms with full path to secret
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should allow manual mount path input if kv mounts are not returned
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should render alert banner on sync error
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should render secret suggestions for nested paths
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should render secret suggestions for selected mount
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should sync database role
- Integration | Component | sync | Secrets::Page::Destinations::Destination::Sync: it should sync secret
- Integration | Component | transform-advanced-templating: it should render

</details>

---

### Elements must only use permitted ARIA attributes

- **Rule:** `aria-prohibited-attr`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/aria-prohibited-attr?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/aria-prohibited-attr?application=axeAPI)
- **Affected tests:** 52

#### Source files and fix

- `ui/lib/core/addon/components/sidebar/frame.hbs`

**Fix:** `<Hds::AppHeader::HomeLink>` is called with `@text="HashiCorp Vault Home Menu"` and `@isIconOnly={{true}}`. The HDS `HomeLink` component renders an `<a>` element and internally applies `aria-label` derived from `@text`. However, axe flags `aria-label` as a prohibited attribute on `<a>` elements that also have a visible text child or when `@isIconOnly` causes the label to conflict with another accessible-name source. The fix is to remove `@isIconOnly={{true}}` and rely on the visually-hidden text mechanism HDS provides, or update to a version of HDS that uses `title` instead of `aria-label` for icon-only links. This element appears on every authenticated page so it has the highest test coverage of any single violation.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<a id="ember17546" class="ember-view hds-app-header__home-link" data-test-app-header-logo="" aria-label="HashiCorp Vault Home Menu">
```

</details>

<details>
<summary>Affected tests (52)</summary>

- Integration | Component | sidebar-frame: it does not render the telemetry consent banner when no prompt is needed
- Integration | Component | sidebar-frame: it renders the telemetry consent banner only when a consent prompt is needed
- Integration | Component | sidebar-frame: it should hide and show app header
- Integration | Component | sidebar-frame: it should hide and show app sidebar
- Integration | Component | sidebar-frame: it should render a submit feedback button in the sidebar that persists regardless of edition
- Integration | Component | sidebar-frame: it should render a submit feedback link in the help menu
- Integration | Component | sidebar-frame: it should render link status, console ui panel container and yield block for app content
- Integration | Component | sidebar-frame: it should render namespace picker in sidebar footer
- Integration | Component | sidebar-nav-access: it should hide links and headings user does not have access to
- Integration | Component | sidebar-nav-access: it should render nav headings
- Integration | Component | sidebar-nav-access: it should render nav links
- Integration | Component | sidebar-nav-agents: it should render nav headings
- Integration | Component | sidebar-nav-agents: it should render nav links
- Integration | Component | sidebar-nav-cluster: it does NOT show Secrets Recovery when user is in HVD admin namespace
- Integration | Component | sidebar-nav-cluster: it should hide Monitoring heading if nav permissions is false
- Integration | Component | sidebar-nav-cluster: it should hide Support nav link in HVD managed clusters
- Integration | Component | sidebar-nav-cluster: it should hide billing metrics link when in namespace other than root or admin when Consumption Billing feature and does not have permission
- Integration | Component | sidebar-nav-cluster: it should hide client counts link in PKI-only Secrets clusters
- Integration | Component | sidebar-nav-cluster: it should hide client counts link in chroot namespace
- Integration | Component | sidebar-nav-cluster: it should hide enterprise related links in child namespace
- Integration | Component | sidebar-nav-cluster: it should hide links and headings user does not have access to
- Integration | Component | sidebar-nav-cluster: it should not show billing metrics link when in HVD admin namespace and Consumption Billing feature is enabled
- Integration | Component | sidebar-nav-cluster: it should render nav headings
- Integration | Component | sidebar-nav-cluster: it should render nav links
- Integration | Component | sidebar-nav-cluster: it should render nav links on community version
- Integration | Component | sidebar-nav-cluster: it should show Monitoring heading when consumption billing is true
- Integration | Component | sidebar-nav-cluster: it should show billing metrics link when in HVD admin namespace and Consumption Billing feature is enabled
- Integration | Component | sidebar-nav-cluster: it should show billing metrics link when in root namespace and Consumption Billing feature is enabled
- Integration | Component | sidebar-nav-reporting: it does NOT Vault Usage if the user has the necessary permission but user is on CE || OSS || community
- Integration | Component | sidebar-nav-reporting: it does NOT show Vault Usage when user is enterprise but not in root namespace
- Integration | Component | sidebar-nav-reporting: it does NOT show Vault Usage when user is user is on CE || OSS || community
- Integration | Component | sidebar-nav-reporting: it does NOT show Vault Usage when user lacks the necessary permission
- Integration | Component | sidebar-nav-reporting: it should hide links user does not have access to
- Integration | Component | sidebar-nav-reporting: it should render nav headings and links
- Integration | Component | sidebar-nav-reporting: it shows Vault Usage when user is enterprise and in root namespace
- Integration | Component | sidebar-nav-reporting: it shows Vault Usage when user is in HVD admin namespace
- Integration | Component | sidebar-nav-resilience-and-recovery: it does NOT show snapshots when user is in HVD admin namespace
- Integration | Component | sidebar-nav-resilience-and-recovery: it should hide links user does not have access to other than secrets recovery
- Integration | Component | sidebar-nav-resilience-and-recovery: it should render nav headings and links
- Integration | Component | sidebar-nav-resilience-and-recovery: it shows Seal Vault when user is enterprise and in root namespace and has nav permissions
- Integration | Component | sidebar-nav-secrets: community: it hides Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: ent (on license), activated and no permissions: it hides Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: ent (on license), activated and permissions: it shows Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: ent (on license), not activated and no permissions: it shows Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: ent (on license), not activated and permissions: it shows Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: ent but feature is not on license: it hides Secrets Sync nav link
- Integration | Component | sidebar-nav-secrets: hvd managed: it shows Secrets Sync nav link regardless of activation status or permissions
- Integration | Component | sidebar-nav-secrets: it should hide links and headings user does not have access to
- Integration | Component | sidebar-nav-secrets: it should render badge for promotional links on managed clusters
- Integration | Component | sidebar-nav-secrets: it should render nav links
- Integration | Component | sidebar-nav-tools: it should hide links user does not have access to
- Integration | Component | sidebar-nav-tools: it should render nav headings and links

</details>

---

### Interactive controls must not be nested

- **Rule:** `nested-interactive`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/nested-interactive?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/nested-interactive?application=axeAPI)
- **Affected tests:** 3

#### Source files and fix

- `ui/app/components/tree-chart.hbs`
- `ui/app/components/wizard/namespaces/step-2.hbs`

**Fix:** A Carbon Charts `<div role="img">` container (the chart holder) has `role="img"` but internally contains `<div role="button">` toolbar controls (`Show as table`, `Make fullscreen`, `More options`). Interactive controls must not be nested inside an element with `role="img"`. Move the toolbar buttons outside the `role="img"` container, or remove `role="img"` from the outer wrapper if it is redundant. In `tree-chart.hbs`, the `<div role="img">` wrapper is Vault-owned — consider using `role="figure"` with a `<figcaption>` instead, which permits interactive children.

<details>
<summary>All unique offending DOM nodes (5)</summary>

```html
<div class="toolbar-control cds--overflow-menu" role="button" aria-disabled="false" aria-label="Make fullscreen">
<div class="toolbar-control cds--overflow-menu" role="button" aria-disabled="false" aria-label="More options">
<div class="toolbar-control cds--overflow-menu" role="button" aria-disabled="false" aria-label="Show as table">
<div role="img" aria-label="Test" data-test-tree-chart="Test" class="cds--chart-holder" data-carbon-theme="g10" style="width: 600px; height: 400px;">
<div role="img" class="tree has-padding-m cds--chart-holder" data-test-tree="" data-carbon-theme="white" style="height: 400px;">
```

</details>

<details>
<summary>Affected tests (3)</summary>

- Integration | Component | page/namespaces | Namespace Wizard: it shows tree chart only when there are multiple globals, orgs, or projects
- Integration | Component | tree-chart: it handles data updates
- Integration | Component | tree-chart: it renders carbon tree chart

</details>

---

### Links must have discernible text

- **Rule:** `link-name`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/link-name?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/link-name?application=axeAPI)
- **Affected tests:** 4

#### Source files and fix

- `ui/lib/pki/addon/components/pki-generate-root.hbs`
- `ui/lib/pki/addon/components/page/pki-issuer-generate-root.hbs`

**Fix:** A `<LinkTo>` renders `{{value}}` as its text content. When the linked model is still loading during the test, `value` is `undefined`/empty, producing an `<a class="loading">` element with no visible text. Add a fallback: `{{or value "(loading…)"}}` inside the `<LinkTo>`, or add `aria-label={{or value "issuer"}}` to ensure the link always has an accessible name even when the text content is momentarily empty.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<a id="ember19468" class="ember-view loading" href="#">
```

</details>

<details>
<summary>Affected tests (4)</summary>

- Integration | Component | page/pki-issuer-generate-root: it renders correct title before and after submit
- Integration | Component | pki-generate-root: it renders Not valid after as radio controls and sends not_after when specific date selected
- Integration | Component | pki-generate-root: it sends ttl when TTL is selected in Not valid after field
- Integration | Component | pki-generate-root: it should use correct endpoint based on issuer permissions

</details>

---

### Scrollable region must have keyboard access

- **Rule:** `scrollable-region-focusable`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/scrollable-region-focusable?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/scrollable-region-focusable?application=axeAPI)
- **Affected tests:** 2

#### Source files and fix

- `ui/tests/index.html`

**Fix:** The `#ember-testing-container` div is the Testem/QUnit test harness container. It is scrollable (overflow set by the test runner styles) but has no `tabindex`, so keyboard users cannot scroll it. This is a test infrastructure element, not production UI. Add `tabindex="0"` to `#ember-testing-container` in `ui/tests/index.html` to satisfy axe during integration/acceptance tests, or suppress this rule globally in the axe configuration for the test environment.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<div id="ember-testing-container">
```

</details>

<details>
<summary>Affected tests (2)</summary>

- Integration | Component | pki | external-pki | ExternalPki::Page::DnsProviders > with DNS providers: it displays configuration details
- Integration | Component | pki | external-pki | ExternalPki::Page::DnsProviders > with DNS providers: it renders list of DNS providers

</details>

---

### [role="img"] elements must have alternative text

- **Rule:** `role-img-alt`
- **Reference:** [https://dequeuniversity.com/rules/axe/4.11/role-img-alt?application=axeAPI](https://dequeuniversity.com/rules/axe/4.11/role-img-alt?application=axeAPI)
- **Affected tests:** 1

#### Source files and fix

- `ui/app/components/tree-chart.hbs`
- `ui/app/components/wizard/namespaces/step-2.hbs`

**Fix:** `<div role="img" aria-label={{@title}}>` in `tree-chart.hbs` — axe requires a non-empty `aria-label` or `aria-labelledby` for `role="img"`. The `@title` arg is set at call sites but may be empty or undefined in some test scenarios. Add a fallback: `aria-label={{or @title "Namespace tree chart"}}`. Also verify that all call sites in `wizard/namespaces/step-2.hbs` pass a non-empty `@title`.

<details>
<summary>All unique offending DOM nodes (1)</summary>

```html
<div role="img" class="tree has-padding-m cds--chart-holder" data-test-tree="" data-carbon-theme="white" style="height: 400px;">
```

</details>

<details>
<summary>Affected tests (1)</summary>

- Integration | Component | page/namespaces | Namespace Wizard: it shows tree chart only when there are multiple globals, orgs, or projects

</details>

---

