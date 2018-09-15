## Module Report
### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/adapters/application.js` at line 8

```js
const { APP } = config;
const { POLLING_URLS, NAMESPACE_ROOT_URLS } = APP;
const { inject, assign, set, RSVP } = Ember;

export default DS.RESTAdapter.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/adapters/cluster.js` at line 6

```js

const { AdapterError } = DS;
const { assert, inject } = Ember;

const ENDPOINTS = ['health', 'seal-status', 'tokens', 'token', 'seal', 'unseal', 'init', 'capabilities-self'];
```

### Unknown Global

**Global**: `Ember.String`

**Location**: `app/adapters/cluster.js` at line 52

```js

  pathForType(type) {
    return type === 'cluster' ? 'clusters' : Ember.String.pluralize(type);
  },

```

### Unknown Global

**Global**: `Ember.String`

**Location**: `app/adapters/transit-key.js` at line 43

```js
        break;
      default:
        path = Ember.String.pluralize(type);
        break;
    }
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/application.js` at line 4

```js
import config from '../config/environment';

const { Controller, computed, inject } = Ember;
export default Controller.extend({
  env: config.environment,
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/helpers/has-feature.js` at line 2

```js
import Ember from 'ember';
const { Helper, inject, observer } = Ember;

const FEATURES = [
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/helpers/is-active-route.js` at line 3

```js
import Ember from 'ember';

const { Helper, inject, observer } = Ember;

const exact = (a, b) => a === b;
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/helpers/is-version.js` at line 2

```js
import Ember from 'ember';
const { Helper, inject, observer } = Ember;

export default Helper.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/helpers/nav-to-route.js` at line 3

```js
import Ember from 'ember';

const { Helper, inject } = Ember;

export default Helper.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/mount-accessor-select.js` at line 4

```js
import { task } from 'ember-concurrency';

const { inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/mount-backend-form.js` at line 6

```js
import { engines } from 'vault/helpers/mountable-secret-engines';

const { inject, computed, Component } = Ember;
const METHODS = methods();
const ENGINES = engines();
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/namespace-link.js` at line 3

```js
import Ember from 'ember';

const { Component, computed, inject } = Ember;

export default Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/namespace-reminder.js` at line 3

```js
import Ember from 'ember';

const { Component, inject, computed } = Ember;

export default Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/namespace-picker.js` at line 7

```js

const { ancestorKeysForKey } = keyUtils;
const { Component, computed, inject } = Ember;
const DOT_REPLACEMENT = '☃';
const ANIMATION_DURATION = 250;
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/not-found.js` at line 3

```js
import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/replication-mode-summary.js` at line 2

```js
import Ember from 'ember';
const { computed, get, getProperties, Component, inject } = Ember;

const replicationAttr = function(attr) {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/replication-summary.js` at line 5

```js
import ReplicationActions from 'vault/mixins/replication-actions';

const { computed, get, Component, inject } = Ember;

const DEFAULTS = {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/role-edit.js` at line 5

```js
import keys from 'vault/lib/keycodes';

const { get, set, computed, inject } = Ember;
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/secret-edit.js` at line 9

```js
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
const { get, computed, inject } = Ember;

export default Ember.Component.extend(FocusOnInsertMixin, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/splash-page.js` at line 3

```js
import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/shamir-flow.js` at line 4

```js
import base64js from 'base64-js';

const { Component, inject, computed, get } = Ember;
const { camelize } = Ember.String;

```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/status-menu.js` at line 3

```js
import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/transit-edit.js` at line 5

```js
import keys from 'vault/lib/keycodes';

const { get, set, computed, inject } = Ember;
const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/ui-wizard.js` at line 4

```js
import { matchesState } from 'xstate';

const { inject, computed } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/wizard-content.js` at line 3

```js
import Ember from 'ember';

const { Component, inject } = Ember;
export default Component.extend({
  wizard: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/access/namespaces/create.js` at line 3

```js
import Ember from 'ember';

const { inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/access/namespaces/index.js` at line 3

```js
import Ember from 'ember';

const { computed, inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/access/namespaces/create.js` at line 4

```js
import UnloadModel from 'vault/mixins/unload-model-route';

const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/access/namespaces/index.js` at line 4

```js
import UnloadModel from 'vault/mixins/unload-model-route';

const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/helpers/set-flash-message.js` at line 3

```js
import Ember from 'ember';

const { Helper, inject } = Ember;

export default Helper.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/mixins/cluster-route.js` at line 3

```js
import Ember from 'ember';

const { get, inject, Mixin, RSVP } = Ember;
const INIT = 'vault.cluster.init';
const UNSEAL = 'vault.cluster.unseal';
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/mixins/policy-edit-controller.js` at line 3

```js
import Ember from 'ember';

let { inject } = Ember;

export default Ember.Mixin.create({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/mixins/replication-actions.js` at line 2

```js
import Ember from 'ember';
const { inject, computed } = Ember;

export default Ember.Mixin.create({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/application.js` at line 4

```js
import ControlGroupError from 'vault/lib/control-group-error';

const { inject } = Ember;
export default Ember.Route.extend({
  controlGroup: inject.service(),
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/routes/vault.js` at line 2

```js
import Ember from 'ember';
const SPLASH_DELAY = Ember.testing ? 0 : 300;

export default Ember.Route.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/models/cluster.js` at line 5

```js
import { fragment } from 'ember-data-model-fragments/attributes';
const { hasMany, attr } = DS;
const { computed, get, inject } = Ember;
const { alias, gte, not } = computed;

```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/auth-form.js` at line 5

```js
import { task } from 'ember-concurrency';
const BACKENDS = supportedAuthBackends();
const { computed, inject, get } = Ember;

const DEFAULTS = {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/auth-info.js` at line 3

```js
import Ember from 'ember';

const { Component, inject, computed, run } = Ember;
export default Component.extend({
  auth: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/config-pki-ca.js` at line 3

```js
import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/control-group-success.js` at line 4

```js
import { task } from 'ember-concurrency';

const { inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/config-pki.js` at line 3

```js
import Ember from 'ember';

const { get, inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/edit-form.js` at line 5

```js
import DS from 'ember-data';

const { inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/control-group.js` at line 4

```js
import { task } from 'ember-concurrency';

const { get, computed, inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/generate-credentials.js` at line 3

```js
import Ember from 'ember';

const { get, set, computed, Component, inject } = Ember;

const MODEL_TYPES = {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/list-item.js` at line 4

```js
import { task } from 'ember-concurrency';

const { inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/logo-splash.js` at line 3

```js
import Ember from 'ember';

const { inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/unseal.js` at line 4

```js
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/policies/index.js` at line 2

```js
import Ember from 'ember';
let { inject } = Ember;

export default Ember.Controller.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/settings/mount-secret-backend.js` at line 6

```js
const SUPPORTED_BACKENDS = supportedSecretBackends();

const { inject, Controller } = Ember;

export default Controller.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/policies/index.js` at line 4

```js
import ClusterRoute from 'vault/mixins/cluster-route';

const { inject } = Ember;

export default Ember.Route.extend(ClusterRoute, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/access/control-group-accessor.js` at line 3

```js
import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';
const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/policies/create.js` at line 5

```js
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

const { inject } = Ember;
export default Ember.Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  version: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/access/control-groups.js` at line 3

```js
import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';
const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/settings/control-groups.js` at line 3

```js
import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';
const { inject } = Ember;

export default Ember.Route.extend(UnloadModel, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/secrets/backend.js` at line 2

```js
import Ember from 'ember';
const { inject } = Ember;
export default Ember.Route.extend({
  flashMessages: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/access/identity/index.js` at line 4

```js
import ListController from 'vault/mixins/list-controller';

const { inject } = Ember;

export default Ember.Controller.extend(ListController, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/access/leases/list.js` at line 4

```js
import utils from 'vault/lib/key-utils';

const { inject, computed, Controller } = Ember;
export default Controller.extend({
  flashMessages: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/services/auth.js` at line 6

```js
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';

const { get, isArray, computed, getOwner, Service, inject } = Ember;

const TOKEN_SEPARATOR = '☃';
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/services/control-group.js` at line 7

```js
const CONTROL_GROUP_PREFIX = 'vault:cg-';
const TOKEN_SEPARATOR = '☃';
const { Service, inject, RSVP } = Ember;

// list of endpoints that return wrapped responses
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/services/namespace.js` at line 5

```js
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const { Service, computed, inject } = Ember;
const ROOT_NAMESPACE = '';
export default Service.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/services/version.js` at line 4

```js
import { task } from 'ember-concurrency';

const { Service, inject, computed } = Ember;

const hasFeatureMethod = (context, featureKey) => {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/services/wizard.js` at line 4

```js
import { Machine } from 'xstate';

const { Service, inject } = Ember;

import getStorage from 'vault/lib/token-storage';
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster.js` at line 6

```js

const POLL_INTERVAL_MS = 10000;
const { inject, Route, getOwner } = Ember;

export default Route.extend(ModelBoundaryRoute, ClusterRoute, {
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/routes/vault/cluster.js` at line 69

```js
    // when testing, the polling loop causes promises to never settle so acceptance tests hang
    // to get around that, we just disable the poll in tests
    return Ember.testing
      ? null
      : Ember.run.later(() => {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster.js` at line 3

```js
import Ember from 'ember';

const { Controller, computed, observer, inject } = Ember;
export default Controller.extend({
  auth: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/auth-config-form/config.js` at line 5

```js
import DS from 'ember-data';

const { inject } = Ember;

const AuthConfigBase = Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/identity/item-details.js` at line 3

```js
import Ember from 'ember';

const { inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/identity/_popup-base.js` at line 2

```js
import Ember from 'ember';
const { assert, inject, Component } = Ember;

export default Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/identity/edit-form.js` at line 5

```js
import { humanize } from 'vault/helpers/humanize';

const { computed, inject } = Ember;
export default Ember.Component.extend({
  flashMessages: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/console/ui-panel.js` at line 13

```js
} from 'vault/lib/console-helpers';

const { inject, computed, getOwner, run } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/identity/lookup-input.js` at line 5

```js
import { underscore } from 'vault/helpers/underscore';

const { inject } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/wizard/mounts-wizard.js` at line 8

```js
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
const supportedAuth = supportedAuthBackends();
const { inject, computed } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/components/wizard/features-selection.js` at line 3

```js
import Ember from 'ember';

const { inject, computed } = Ember;

export default Ember.Component.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/auth.js` at line 4

```js
import { task, timeout } from 'ember-concurrency';

const { inject, computed, Controller } = Ember;
export default Controller.extend({
  vaultController: inject.controller('vault'),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/controllers/vault/cluster/settings.js` at line 3

```js
import Ember from 'ember';

const { inject, Controller } = Ember;
export default Controller.extend({
  namespaceService: inject.service('namespace'),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/auth.js` at line 5

```js
import config from 'vault/config/environment';

const { inject } = Ember;

export default ClusterRouteBase.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/logout.js` at line 4

```js
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

const { inject } = Ember;
export default Ember.Route.extend(ModelBoundaryRoute, {
  auth: inject.service(),
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/policies.js` at line 5

```js

const ALLOWED_TYPES = ['acl', 'egp', 'rgp'];
const { inject } = Ember;

export default Ember.Route.extend(ClusterRoute, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/init.js` at line 4

```js
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/policy.js` at line 5

```js

const ALLOWED_TYPES = ['acl', 'egp', 'rgp'];
const { inject } = Ember;

export default Ember.Route.extend(ClusterRoute, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/replication.js` at line 3

```js
import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';
const { inject } = Ember;

export default Ember.Route.extend(ClusterRoute, {
```

### Unknown Global

**Global**: `Ember.inject`

**Location**: `app/routes/vault/cluster/settings/auth/configure/section.js` at line 5

```js
import UnloadModelRoute from 'vault/mixins/unload-model-route';

const { RSVP, inject } = Ember;
export default Ember.Route.extend(UnloadModelRoute, {
  modelPath: 'model.model',
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/leases-test.js` at line 19

```js
moduleForAcceptance('Acceptance | leases', {
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;

```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/leases-test.js` at line 20

```js
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;

    authLogin();
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/leases-test.js` at line 28

```js
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    return authLogout();
  },
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 10

```js
moduleForAcceptance('Acceptance | policies (old)', {
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 11

```js
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 12

```js
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 13

```js
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
  },
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 17

```js
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
  },
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/policies-acl-old-test.js` at line 18

```js
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
  },
});
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 10

```js
moduleForAcceptance('Acceptance | not-found', {
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 11

```js
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 12

```js
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 13

```js
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
  },
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 17

```js
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    return authLogout();
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 18

```js
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    return authLogout();
  },
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/settings/configure-secret-backends/pki/section-urls-test.js` at line 8

```js
moduleForAcceptance('Acceptance | settings/configure/secrets/pki/urls', {
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;
    return authLogin();
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/settings/configure-secret-backends/pki/section-urls-test.js` at line 9

```js
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;
    return authLogin();
  },
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/settings/configure-secret-backends/pki/section-urls-test.js` at line 13

```js
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
  },
});
```
