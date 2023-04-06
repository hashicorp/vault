# Routing

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents**

- [Routing](#routing)
  - [Guidelines](#guidelines)
  - [File structure](#file-structure)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Guidelines

- Parent route typically serves to group related child resources
- Parent index route typically displays empty state placeholder with call to action or redirects to default child resource
- Child resource names are pluralized
- Child index route represents list view.
- Child singularized name /details is the read

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

> Example with [OIDC](../app/routes/vault/cluster/access/oidc/) in codebase [sample PR](https://github.com/hashicorp/vault/pull/16028):

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
│   │   │   │   ├── providers.js <- utilizes the `modelFor` method to get parent `client`
```
