# Routing

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Guidelines](#guidelines)
- [File structure](#file-structure)
- [Shared functionality](#shared-functionality)
- [Decorators](#decorators)
  - [@withConfirmLeave()](#withconfirmleave)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Guidelines

- Parent route typically serves to group related child resources
- Parent index route typically displays empty state placeholder with call to action or redirects to default child resource
- Child resource names are pluralized
- Child index route represents list view.
- Child singularized name + /details is the read view.
- _Avoid_ extending routes. This can lead to unnecessary inheritance which gets messy quickly. For [shared functionality](#shared-functionality), consider a decorator.

## File structure

The file structure can be leveraged to simplify CRUD actions and passing data. The singular resource route should live at the same level as its folder, this automatically passes its model to any route nested within the folder.
Below, `details.js` and `edit.js` will automatically receive the model returned by the model hook in `resource-foo.js`. Alternately, if defining a custom model hook in those routes, we can use methods like `this.modelFor` instead of re-querying records.

```
├── routes/vault/cluster/access
│   ├── parent/
│   │   ├── index.js
│   │   ├── resource-foos /
│   │   │   ├── resource-foo.js
│   │   │   ├── create.js
│   │   │   ├── index.js
│   │   │   ├── resource-foo/
│   │   │   │   ├── details.js
│   │   │   │   ├── edit.js
```

> For example, [OIDC](../app/routes/vault/cluster/access/oidc/) route structure [_original PR_](https://github.com/hashicorp/vault/pull/16028):

```
├── routes/vault/cluster/access
│   ├── oidc/
│   │   ├── index.js
│   │   ├── clients/
│   │   │   ├── client.js
│   │   │   ├── create.js
│   │   │   ├── index.js
│   │   │   ├── client/
│   │   │   │   ├── details.js
│   │   │   │   ├── edit.js
│   │   │   │   ├── providers.js <- utilizes the modelFor method to get id about parent's clientId
```

## Shared functionality

To guide users, we sometimes have a call to action that depends on a resource's state. For example, if a secret engine hasn't been configured routing to the first step to do so, and otherwise navigating to its overview page.

Instead of extending route classes to share this `isConfigured` state, consider a decorator! [withConfig()](../../ui/lib/kubernetes/addon/decorators/fetch-config.js) is a great example.

## Decorators

### [@withConfirmLeave()](../lib/core/addon/decorators/confirm-leave.js)

- Renders `window.confirm()` alert that a user has unsaved changes if navigaing away from route with the decorator
- Unloads or rolls back Ember data model record

<!-- TODO add withConfig() if we refactor for more general use -->

<!-- ### [withConfig()](../../ui/lib/kubernetes/addon/decorators/fetch-config.js)

We sometimes have a call to action guiding users that depends on a resource's state. For example, if a secret engine hasn't been configured the UI renders an empty state linking to the first configuration step. Otherwise, it routes to the overview page of that engine.

Sample use:

```js
import { withConfig } from '../decorators/fetch-config';
@withConfig()
export default class SomeRouter extends Route {
  model() {
    // in case of any error other than 404 we want to display that to the user
    if (this.configError) {
      throw this.configError;
    }
    return {
      config: this.configModel, // configuration data to determine UI state
    };
  }
}
``` -->
