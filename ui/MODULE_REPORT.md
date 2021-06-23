## Module Report
### Unknown Global

**Global**: `Ember.onerror`

**Location**: `tests/helpers/wait-for-error.js` at line 5

```js

export default function waitForError(opts) {
  const orig = Ember.onerror;

  let error = null;
```

### Unknown Global

**Global**: `Ember.onerror`

**Location**: `tests/helpers/wait-for-error.js` at line 5

```js

export default function waitForError(opts) {
  const orig = Ember.onerror;

  let error = null;
```

### Unknown Global

**Global**: `Ember.onerror`

**Location**: `tests/helpers/wait-for-error.js` at line 8

```js

  let error = null;
  Ember.onerror = err => {
    error = err;
  };
```

### Unknown Global

**Global**: `Ember.onerror`

**Location**: `tests/helpers/wait-for-error.js` at line 13

```js

  return waitUntil(() => error, opts).finally(() => {
    Ember.onerror = orig;
  });
}
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 15

```js

  hooks.beforeEach(function() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 16

```js
  hooks.beforeEach(function() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 17

```js
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authPage.login();
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 18

```js
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authPage.login();
  });
```

### Unknown Global

**Global**: `Ember.Test`

**Location**: `tests/acceptance/not-found-test.js` at line 23

```js

  hooks.afterEach(function() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    return logout.visit();
```

### Unknown Global

**Global**: `Ember.Logger`

**Location**: `tests/acceptance/not-found-test.js` at line 24

```js
  hooks.afterEach(function() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    return logout.visit();
  });
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/components/auth-jwt.js` at line 9

```js

/* eslint-disable ember/no-ember-testing-in-module-scope */
const WAIT_TIME = Ember.testing ? 0 : 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete.  Please click Sign In to try again.';
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/components/auth-jwt.js` at line 119

```js
  exchangeOIDC: task(function*(event, oidcWindow) {
    // in non-incognito mode we need to use a timeout because it takes time before oidcState is written to local storage.
    let oidcState = Ember.testing
      ? event.storageArea.getItem('oidcState')
      : yield timeout(1000).then(() => event.storageArea.getItem('oidcState'));
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/routes/vault.js` at line 7

```js
import Ember from 'ember';
/* eslint-disable ember/no-ember-testing-in-module-scope */
const SPLASH_DELAY = Ember.testing ? 0 : 300;

export default Route.extend({
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/services/auth.js` at line 267

```js
  checkShouldRenew: task(function*() {
    while (true) {
      if (Ember.testing) {
        return;
      }
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `lib/core/addon/components/ttl-form.js` at line 82

```js
    this.set('time', parsedTime);
    this.handleChange();
    if (Ember.testing) {
      return;
    }
```
