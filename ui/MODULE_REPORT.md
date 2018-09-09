## Module Report
### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/routes/vault.js` at line 6

```js
import Route from '@ember/routing/route';
import Ember from 'ember';
const SPLASH_DELAY = Ember.testing ? 0 : 300;

export default Route.extend({
```

### Unknown Global

**Global**: `Ember.testing`

**Location**: `app/routes/vault/cluster.js` at line 81

```js
    // when testing, the polling loop causes promises to never settle so acceptance tests hang
    // to get around that, we just disable the poll in tests
    return Ember.testing
      ? null
      : later(() => {
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
