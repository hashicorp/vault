import Ember from 'ember';
import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import page from 'vault/tests/pages/settings/configure-secret-backends/pki/section';
let adapterException;
moduleForAcceptance('Acceptance | settings/configure/secrets/pki/urls', {
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;
    return authLogin();
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
  },
});

test('it saves urls config', function(assert) {
  const path = `pki-${new Date().getTime()}`;
  mountSupportedSecretBackend(assert, 'pki', path);
  page.visit({ backend: path, section: 'urls' });
  andThen(() => {
    assert.equal(currentRouteName(), 'vault.cluster.settings.configure-secret-backend.section');
  });

  page.form.fields(0).input('foo').change();
  page.form.submit();

  andThen(() => {
    assert.ok(page.form.hasError, 'shows error on invalid input');
  });

  page.form.fields(0).input('foo.example.com').change();
  page.form.submit();

  andThen(() => {
    assert.equal(page.lastMessage, 'The urls config for this backend has been updated.');
  });
});
