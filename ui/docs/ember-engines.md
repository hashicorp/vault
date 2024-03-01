# [Ember Engines](https://ember-engines.com/docs)

This is a quickstart guide inspired by [ember engine quickstart](https://ember-engines.com/docs/quickstart) on how to set up an ember engine in Vault!

## Create a new in-repo engine:

`ember g in-repo-engine <engine-name>`

_This blueprint in-repo engine command will add a new folder `lib/<engine-name>` and add the engine to our main app’s `package.json`_

## Engine’s package.json:

```json
{
  "name": "<engine-name>",

  "dependencies": {
    "ember-cli-htmlbars": "*",
    "ember-cli-babel": "*"
  },

  "ember-addon": {
    "paths": ["../core"]
  }
}
```

For our application, we want to include the **ember-addon** path `../core`

By adding this **ember-addon** path, we are able to share elements between your in-repo addon and the Vault application[^1].

## Configure your Engine

In the engine’s `index.js` file:

```js
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/* eslint-disable node/no-extraneous-require */
const { buildEngine } = require('ember-engines/lib/engine-addon');

module.exports = buildEngine({
  name: '<engine-name>',
  lazyLoading: {
    enabled: false,
  },
  isDevelopingAddon() {
    return true;
  },
});
```

Within your Engine’s `config/environment.js` file:

```js
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// config/environment.js

'use strict';

module.exports = function (environment) {
  const ENV = {
    modulePrefix: '<engine-name>',
    environment: environment,
  };

  return ENV;
};
```

Within your Engine’s `addon/engine.js` file:

```js
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Engine from '@ember/engine';

import loadInitializers from 'ember-load-initializers';
import Resolver from 'ember-resolver';

import config from './config/environment';

const { modulePrefix } = config;

export default class <EngineName>Engine extends Engine {
  modulePrefix = modulePrefix;
  Resolver = Resolver;
  dependencies = {
    services: ['router', 'store', 'secret-mount-path', 'flash-messages'],
    externalRoutes: ['secrets'],
  };
}

loadInitializers(<EngineName>Engine, modulePrefix);
```

### Service Dependencies:

The services in the example above are common services that we often use in our engines. If your engine requires other services from the main application, add them to the services array.

#### Notes:

- Service dependencies are OPTIONAL. If your engine does not use any external services, you do not need to include a services dependency array.
- Remember to include any dependencies here in the engine's dependencies in app/app.js

### External Route Dependencies:

The external route dependencies allow you to link to a route outside of your engine. In this example, we list 'secrets' in the externalRoute and the route is defined in the `app.js` file.

#### Notes:

- In order to link to the other routes in the main app using the `LinkToExternal` component from your engine, you need to add the route to the `app/app.js` and your engine’s `addon/engine.js`. More information on [Linking to An External Context.](https://ember-engines.com/docs/link-to-external).

## Additional info about your engine's `application.hbs`:

- Optional step: Add some text in the engine’s `application.hbs` file (to see if your engine was set up correctly).
- **NOTE: Most of our existing engines do not keep the generated `application.hbs` template file. If nothing will be added to it and it remains as just an `{{outlet}}` it can safely be removed.**

## Register your engine with our main application:

In our `app/app.js` file in the engines object, add your engine’s name and dependencies.

```js
/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Application from '@ember/application';
import Resolver from 'ember-resolver';
import loadInitializers from 'ember-load-initializers';
import config from 'vault/config/environment';

export default class App extends Application {
	...
  engines = {
    <engine-name>: {
      dependencies: {
        services: ['router', 'store', 'secret-mount-path', 'flash-messages', <any-other-dependencies-you-have>],
        externalRoutes: {
          secrets: 'vault.cluster.secrets.backends',
        },
      },
    },
  };
}

loadInitializers(App, config.modulePrefix);
```

If you used `ember g in-repo-engine <engine-name>` to generate the engine’s blueprint, it should have added `this.mount(<engine-name>)` to the main app’s `router.js` file (this adds your engine and its associated routes). \*Move `this.mount(<engine-name>)` to match your engine’s route structure. For more information about [Routable Engines](https://ember-engines.com/docs/quickstart#routable-engines).

## Add engine path to ember-addon section of main app package.json

```json
  "ember-addon": {
    "paths": [
      "lib/core",
      "lib/your-new-engine"
    ]
  },
```

### Important Notes:

- Anytime a new engine is created, you will need to `yarn install` and **RESTART** ember server!
- To add `package.json` **dependencies** or **devDependencies**, you can copy + paste the dependency into corresponding sections. Most of the time, we will want to use "\*" in place of the version number to ensure all the dependencies have the latest version.

### Common blueprint commands:

- **Generating In-repo engines routes:** `ember generate route <route-name> --in-repo <in-repo-name>` - _generates tests and route files and templates_
- **Remove In-repo engines routes:** `ember destroy route <route-name> --in-repo <in-repo-name>` - _removes tests and route files and templates_
- **Generating In-repo engines components:** `ember generate component <component-name> --in-repo <in-repo-name>`- _generates tests and component files and templates_
- **Remove In-repo engines components:** `ember destroy component <component-name> --in-repo <in-repo-name>`- _removes tests and component files and templates_

[^1]: https://ember-engines.com/docs/quickstart#create-as-in-repo-engine
