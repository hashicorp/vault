# Writing and consuming components

Components can range from small, highly reusable "atoms" to large units with lots of business logic specific to one workflow or action. In any scenario, these are things to keep in mind while developing components for the Vault UI.

Please note that these guidelines are aspirational and you will see instances of antipatterns in the codebase. Many of these should be updated as we move forward. As with any ruleset, sometimes it is appropriate to break the rule.

## Page components for every route

Route templates should render a `Page` component, which includes breadcrumbs, page title, and then renders whatever else should be on the page (often another scoped component).

- This component should be named something like `<Page::CreateFoo />` and you can create it like `ember g component page/create-foo -gc`.
- The Route should pass the model hook data into the component in the template. However, if the model hook returns multiple objects they should each be passed into the Page component as separate args. For example: within a route's template whose model hook returns two different data models, the route's template would look like:

```hbs
<Page::CreateFoo @config={{this.model.config}} @foo={{this.model.foo}} />
```

## Conditional rendering

Generally, we want the burden of deciding whether a component should render to live in the parent rather than the child.

- **Readability** - it's easier to tell at a glance that a component will sometimes not render if the `{{#if}}` block is on the parent.
- **Performance** - when a component is in charge of its own rendering logic, the component's lifecycle hooks will fire whether or not the component will render on the page. This can lead to degraded performance, for example if hundreds of the same component are listed on the page.

## Reusable components

When developing components, make sure to:

- Add splattributes to the top level, eg:

  ```hbs
  <div data-test-stuff ...attributes>Stuff!</div>
  ```

- Consider passing splattributes or yielding something instead of passing a new arg

  ❌ **Instead of:** passing a new arg that controls a style

  ```
  <Block @title="Example" @hasPadding={{false}} />
  ```

  ✅ **Prefer:** passing a class or helper that controls a style

  ```
  <Block @title="Example" class="padding-0" />
  ```

- Minimize the number of args that must be passed

  ❌ **Instead of:** Passing in separate args that are both required for icon to render

  ```
  <Block @title="Example" @hasIcon={{true}} @iconName="key" />
  ```

  ✅ **Prefer:** One arg that is rendered if present

  ```
  <Block @title="Example" @icon="key" />
  ```
