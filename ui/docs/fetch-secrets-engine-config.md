# Fetch Secrets Engine Configuration Decorator

The `fetch-secrets-engine-config` decorator is available in the core addon and can be used on a route that needs to be aware of the configuration details of a secrets engine prior to model hook execution. This is useful for conditionally displaying a call to action for the user to complete the configuration.

## API

The decorator accepts a single argument with the name of the Ember Data model to be fetched.

- **modelName** [string] - name of the Ember Data model to fetch which is passed to the `queryRecord` method.

With the provided model name, the decorator fetches the record using the store `queryRecord` method in the `beforeModel` route hook. Several properties are set on the route class based on the status of the request:

- **configModel** [Model | null] - set on success with resolved Ember Data model.

- **configError** [AdapterError | null] - set if the request errors with any status other than 404.

- **promptConfig** [boolean] - set to `true` if the request returns a 404, otherwise set to `false`. This is for convenience since checking for `(!this.configModel && !this.configError)` would result in the same value.

## Usage

### Configure route

```js
@withConfig('foo/config')
export default class FooConfigureRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return this.configModel || this.store.createRecord('foo/config', { backend });
  }
}
```

In the scenario of creating/updating the configuration, the model is used to populate the form if available, otherwise the form is presented in an empty state. Fetch errors are not a concern, nor is prompting the user to configure so only the `configModel` property is used.

### Configuration route

```js
@withConfig('foo/config')
export default class FooConfigurationRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    // the error could also be thrown to display the error template
    // in this example a component is used to display the error
    return {
      configModel: this.configModel,
      configError: this.configError,
    };
  }
}
```

For configuration routes, the model and error properties may be used to determine what should be displayed to the user:

`configuration.hbs`

```hbs
{{#if @configModel}}
  {{#each @configModel.fields as |field|}}
    <InfoTableRow @label={{field.label}} @value={{field.value}} />
  {{/each}}
{{else if @configError}}
  <Page::Error @error={{@configError}} />
{{else}}
  <ConfigCta />
{{/if}}
```

### Other routes (overview etc.)

This is the most basic usage where a route only needs to be aware of whether or not to show the config prompt:

```js
@withConfig('foo/config')
export default class FooOverviewRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    const backend = this.secretMountPath.currentPath;
    return hash({
      promptConfig: this.promptConfig,
      roles: this.store.query('foo/role', { backend }).catch(() => []),
      libraries: this.store.query('foo/library', { backend }).catch(() => []),
    });
  }
}
```
