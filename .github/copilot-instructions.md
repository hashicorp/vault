# GitHub Copilot Instructions for Vault

This document outlines coding standards, best practices, and project-specific conventions for our monorepo. GitHub Copilot should refer to these guidelines when providing suggestions, completing code, and assisting with code reviews.

**Prioritize reusability, testability, and maintainability in all code changes across the entire monorepo.**

---

## UI (Ember.js - `vault/ui` Directory) Guidelines

These instructions apply specifically to code within the `vault/ui` (Ember.js) directory.

### For Files in the `changelog` Directory

**Instruction:** Ensure `changelog` entries for enterprise-only changes explicitly indicate "enterprise" in their description.

* **DON'T:** `ui: descriptive text`
* **DO:** `ui (enterprise): descriptive text`

### For Files in the `components` Directory

**Instruction:** Do not use `@tracked` on properties that are immutable or do not change throughout the component's lifecycle. `@tracked` is for reactive state.

* **DON'T:** `@secret={{concat "role/" @model.id}}` (Here, `@secret` is an argument, not a reactive property of the component itself.)
* **DO:** `@tracked myReactiveProperty = "initialValue";` (Use `@tracked` for component-internal state that can change.)

**Instruction:** Component class names should always match their corresponding file names for clarity and Ember's convention.

**Instruction:** Co-locate `.hbs` templates in the same directory as their corresponding `.js` or `.ts` files (e.g., `vault/components/my-component/my-component.js` and `vault/components/my-component/my-component.hbs`).

**Instruction:** Avoid unnecessary quotation marks around dynamic data attributes in templates.

* **DON'T:** `data-test-namespace-link="{{option.label}}"`
* **DO:** `data-test-namespace-link={{option.label}}`

### For Files in the `models` Directory

**Instruction:** Use single-line syntax for `@attr` declarations if they have one or fewer keys in their options object.

* **DON'T:**
    ```javascript
    @attr('string', {
      label: 'Client ID',
    }) clientId;
    ```
* **DO:** `@attr('string', { label: 'Client ID' }) clientId;`

**Instruction:** Do not add extra blank lines between consecutive single-line `@attr` declarations to maintain compactness.

* **DON'T:**
    ```javascript
    @attr('string', { label: 'Client ID' }) clientId;

    @attr('string', { label: 'Client Secret' }) clientSecret;
    ```
* **DO:**
    ```javascript
    @attr('string', { label: 'Client ID' }) clientId;
    @attr('string', { label: 'Client Secret' }) clientSecret;
    ```

### For Files Ending in `.js` or `.ts` (within `vault/ui`)

**Instruction:** Include JSDoc-style documentation (e.g., `/** @module ComponentName */` for modules, `/** ... */` for classes, methods, and properties) for all new files, especially components, helpers, services, and complex functions.

**Instruction:** Remove all unused code, imports, or constants to keep the codebase clean and efficient.

**Instruction:** Avoid using `console.error`. Instead, use Ember's `assert` for runtime checks or `debug` for logging during development.

**Instruction:** Ensure all comments are up-to-date and accurately reflect the current code logic. Remove or update outdated comments.

**Instruction:** Avoid using the `any` type in TypeScript files. Strive for strict typing to improve code predictability and maintainability.

**Instruction:** Only use `@task` from `ember-concurrency` when you specifically need to leverage its advanced features like `isRunning`, `perform`, `cancel`, or task modifiers. Otherwise, prefer standard `async`/`await` with `try`/`catch` blocks for asynchronous operations.

**Instruction:** If `setTimeout` is used, consider whether `requestAnimationFrame`, `@tracked` properties with getters, or other reactive Ember patterns might be more appropriate. `setTimeout` can lead to testing complexities and event loop management issues.

**Instruction:** Do not use `new Date()` directly. To ensure consistent UTC date calculations across all environments, use `Date.UTC(...)` and `getUTCFullYear()`, `getUTCMonth()`, etc. for date manipulation.

### For Files in the `tests` Directory (within `vault/ui`)

**Instruction:** Replace `setTimeout` with Ember's `run.later` method within tests for better control over the Ember runloop.

**Instruction:** When using `ember-cli-mirage`, place all `this.server` setup steps at the top of the `test` or `module` block for clarity and consistency.

**Instruction:** Avoid using `assert.ok()`. Prefer `assert.true()` or `assert.false()` for clearer boolean assertions.

**Instruction:** Provide clear, descriptive messages for all assertions to aid in debugging test failures.

* **DON'T:** `assert.dom(GENERAL.messageError).hasText('Error');`
* **DO:** `assert.dom(GENERAL.messageError).hasText('Error', "Verify that the error message is displayed correctly.");`

**Instruction:** Use `data-test-*` selectors for all DOM interactions within tests. This decouples tests from presentational styles and markup changes.

**Instruction:** Avoid introducing new Page Object patterns. Instead, interact with elements directly using general selectors (e.g., `data-test-*` selectors, standard CSS selectors) within your test files. This simplifies test maintenance and reduces abstraction layers.

### For Files Ending in `.hbs` (Handlebars Templates within `vault/ui`)

**Instruction:** Avoid using the `.length` property in logical operators within `{{if}}` or `{{unless}}` helpers for arrays. Empty arrays are already considered falsy in Ember templates.

* **DON'T:** `{{#if (gt this.model.allowed_roles.length 0)}}`
* **DO:** `{{#if this.model.allowed_roles}}`

**Instruction:** Avoid using the `{{concat}}` helper. Prefer string interpolation directly within attributes for cleaner syntax.

* **DON'T:** `@secret={{concat "role/" @model.id}}`
* **DO:** `@secret="role/{{@model.id}}"`

**Instruction:** Ensure that any links leading to `vaut/docs` are rendered using `Hds::Link::Inline` and not a standard `<button>`.

**Instruction:** Do not use unnecessary quotation marks around double curly brace expressions when passing dynamic values to component arguments.

* **DON'T:** `@user="{{@model.name}}"`
* **DO:** `@user={{@model.name}}`

**Instruction:** For `selected` attributes, the passed-in property should almost always be dynamic. If a static value (e.g., `true` or a fixed string) is used, suggest a review to confirm it is intentional and correct.

* **DON'T:** `selected="true"`
* **DO:** `selected={{eq this.selectedAuthMethod type}}`

**Instruction:** If a conditional wraps an entire element, refactor it so the conditional wraps only the dynamic *content* within the element, improving readability and reducing HTML boilerplate.

* **DON'T:**
    ```handlebars
    {{#if this.version.isEnterprise}}
      <PH.Title>Enterprise things</PH.Title>
    {{else}}
      <PH.Title>Community things</PH.Title>
    {{/if}}
    ```
* **DO:** `<PH.Title>{{if this.version.isEnterprise "Enterprise things" "Community things"}}</PH.Title>`

**Instruction:** Avoid using the `style` attribute for inline styling. Define and use CSS classes within `.scss` files instead.

* **DON'T:**
    ```handlebars
    <Hds::Button style="color: red;" />
    ```
* **DO:**
    ```handlebars
    <Hds::Button class="my-custom-button" />
    ```
    (with `.my-custom-button { color: red; }` in a SCSS file)

**Instruction:** Place `data-test-*` selectors as the *last* attribute on an HTML element for consistency and ease of parsing.

* **DON'T:**
    ```handlebars
    <Hds::Button
      data-test-save
      @text={{if (eq @mountType "secret") "Enable engine" "Enable method"}}
      type="submit"
    />
    ```
* **DO:**
    ```handlebars
    <Hds::Button
      @text={{if (eq @mountType "secret") "Enable engine" "Enable method"}}
      type="submit"
      data-test-save
    />
    ```

### For Files Ending in `.scss` (within `vault/ui`)

**Instruction:** Avoid using `z-index`. Instead, manage element stacking order by adjusting their natural order in the template (DOM structure).

* **DON'T:**
    ```css
    .modal {
      position: absolute;
      z-index: 10;
    }
    ```
* **DO:** (This instruction implies the HTML change, so the SCSS `DO` is the absence of `z-index`).
    ```css
    /* Adjust HTML structure to control stacking instead of z-index */
    .modal {
      position: absolute;
      /* no z-index needed if stacking context is managed by DOM order */
    }
    ```

**Instruction:** Avoid using `!important`. Instead, achieve desired style overrides by increasing CSS specificity (e.g., by targeting elements with more specific selectors or by using component-scoped styles).

* **DON'T:**
    ```css
    .button {
      color: red !important;
    }
    ```
* **DO:**
    ```css
    .namespace-picker .button { /* More specific selector */
      color: red;
    }
    ```

### For Changes to `package.json` (within `vault/ui` - if applicable)

**Instruction:** Ensure `package.json` changes are independent and minimal, relating only to the immediate feature or bug fix. The `yarn.lock` file is the only expected co-change.

**Instruction:** When adding or modifying dependencies, prefer pinning exact versions or using the tilde (`~`) for patch-level updates. Avoid using the caret (`^`) operator for dependency versions.

* **DON'T:** `"ansi-html": "^0.0.8"`
* **DO:** `"ansi-html": "~0.0.8"` or `"ansi-html": "0.0.8"`

**Instruction:** Dependencies within the `resolutions` block must always be pinned to an exact version (no `~` or `^`).

---

## General UI Guidelines (within `vault/ui`)

These instructions provide general coding principles and best practices specifically for the `vault/ui` (Ember.js) directory.

**Instruction:** Replace deprecated `Route.extend`, `Model.extend`, or `Component.extend` syntax with their modern Ember Octane class-based equivalents.

**Instruction:** Place comments directly above the code lines or blocks they describe, not below or interspersed.

**Instruction:** Use sentence case for all titles (e.g., for component arguments like `@title`).

* **DON'T:** `@title="Upload Users'sProfile"`
* **DO:** `@title="Upload user's profile"`

**Instruction:** Ensure all subtitles and descriptive text follow proper grammar rules, including ending sentences with periods.

**Instruction:** All new components (but not tests) should be co-located within their logical structure inside the `vault/components` directory.

**Instruction:** Adhere to Ember's built-in best practices for code structure and syntax, consulting the official Ember documentation for the version you are using.

**Instruction:** Prioritize creating reusable and maintainable code. Avoid overly complex or one-off component/route implementations without strong justification.

**Instruction:** When referring to KV secret engines, use the precise terms "KVv2" or "KVv1". When space allows, prefer spelling out "KV version 2" or "KV version 1".

* **DON'T:** `KV v2`
* **DO:** `KV version 2` or `KVv2` (depending on space/context)

---

## Debugging Tips (Informational, not for Copilot to "enforce")

These tips are for developers working within the `vault/ui` directory.

* **Tip:** When using `{{debugger}}` inside a template, you can inspect values in the browser console.
    * **Example:**
        ```handlebars
        {{#each this.items as |item|}}
          {{debugger}}
        {{/each}}
        ```
    * **In the console:**
        * Use `get('item.name')` to inspect an item's property.
        * Use `context` to explore the current rendering context.

---

## API (e.g., [Your API Language/Framework] - `vault/api` Directory) Guidelines

These instructions apply specifically to code within the `vault/api` directory.

**Instruction:** Add your API-specific guidelines here.

* *Example: All database migrations must be reversible.*
* *Example: Ensure all API endpoints include proper authentication and authorization checks.*
* *Example: Prefer using [Your API ORM] for database interactions.*