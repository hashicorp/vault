import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import Ember from 'ember';

let adapterException;
let loggerError;

moduleForAcceptance('Acceptance | not-found', {
  beforeEach() {
    loggerError = Ember.Logger.error;
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    Ember.Logger.error = () => {};
    return authLogin();
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    Ember.Logger.error = loggerError;
    return authLogout();
  },
});

test('top-level not-found', function(assert) {
  visit('/404');
  andThen(() => {
    assert.ok(find('[data-test-not-found]').length, 'renders the not found component');
    assert.ok(find('[data-test-header-without-nav]').length, 'renders the not found component with a header');
  });
});

test('vault route not-found', function(assert) {
  visit('/vault/404');
  andThen(() => {
    assert.dom('[data-test-not-found]').exists('renders the not found component');
    assert.ok(find('[data-test-header-with-nav]').length, 'renders header with nav');
  });
});

test('cluster route not-found', function(assert) {
  visit('/vault/secrets/secret/404/show');
  andThen(() => {
    assert.dom('[data-test-not-found]').exists('renders the not found component');
    assert.ok(find('[data-test-header-with-nav]').length, 'renders header with nav');
  });
});

test('secret not-found', function(assert) {
  visit('/vault/secrets/secret/show/404');
  andThen(() => {
    assert.dom('[data-test-secret-not-found]').exists('renders the message about the secret not being found');
  });
});
