# Writing and consuming components

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Page components for every route](#page-components-for-every-route)
- [Conditional rendering](#conditional-rendering)
- [Reusable components](#reusable-components)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Components can range from small, highly reusable "atoms" to large units with lots of business logic specific to one workflow or action. In any scenario, these are things to keep in mind while developing components for the Vault UI.

Please note that these guidelines are aspirational and you will see instances of anti-patterns in the codebase. Many of these should be updated as we move forward. As with any rule set, sometimes it is appropriate to break the rule.

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

- Consider yielding something instead of passing a new arg
- Less is more! Adding lots of rendering logic means the component is likely doing too much

| 💡 Tips for reusability                                           | Example                                                                    |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------- |
| ✅ Add splattributes to the top level                             | <pre>`<div ...attributes> Something! </div>`</pre>                         |
| ✅ Pass a class or helper that controls a style                   | <pre>`<Block @title="Example" class="padding-0" />`</pre>                  |
| ❌ Don't pass a new arg that controls a style                     | <pre>`<Block @title="Example" @hasPadding={{false}} />` </pre>             |
| ✅ Minimize args passed, pass one arg that is rendered if present | <pre>`<Block @title="Example" @icon="key" />`</pre>                        |
| ❌ Don't pass separate args required for icon to render           | <pre>`<Block @title="Example" @hasIcon={{true}} @iconName="key" />` </pre> |
