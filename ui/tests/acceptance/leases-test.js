import { test } from 'qunit';
import moduleForAcceptance from 'vault/tests/helpers/module-for-acceptance';
import Ember from 'ember';

let adapterException;
// testing error states is terrible in ember acceptance tests so these weird Ember bits are to work around that
// adapted from https://github.com/emberjs/ember.js/issues/12791#issuecomment-244934786
moduleForAcceptance('Acceptance | leases', {
  beforeEach() {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => null;
    return authLogin();
  },
  afterEach() {
    Ember.Test.adapter.exception = adapterException;
    return authLogout();
  },
});

const createSecret = (context, isRenewable) => {
  const now = new Date().getTime();
  const secretContents = { secret: 'foo' };
  if (isRenewable) {
    secretContents.ttl = '30h';
  }
  context.secret = {
    name: isRenewable ? `renew-secret-${now}` : `secret-${now}`,
    text: JSON.stringify(secretContents),
  };
  //create a secret so we have a lease (server is running in -dev-leased-kv mode)
  visit('/vault/secrets/secret/list');
  click('[data-test-secret-create]');
  fillIn('[data-test-secret-path]', context.secret.name);
  andThen(() => {
    const codeMirror = find('.CodeMirror');
    // UI keeps state so once we flip to json, we don't need to again
    if (!codeMirror.length) {
      click('[data-test-secret-json-toggle]');
    }
  });
  andThen(() => {
    find('.CodeMirror').get(0).CodeMirror.setValue(context.secret.text);
  });
  click('[data-test-secret-save]');
};

const navToDetail = secret => {
  visit('/vault/access/leases/');
  click('[data-test-lease-link="secret/"]');
  click(`[data-test-lease-link="secret/${secret.name}/"]`);
  click(`[data-test-lease-link]:eq(0)`);
};

test('it renders the show page', function(assert) {
  createSecret(this);
  navToDetail(this.secret);
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert.equal(
      find('[data-test-lease-renew-picker]').length,
      0,
      'non-renewable lease does not render a renew picker'
    );
  });
});

test('it renders the show page with a picker', function(assert) {
  createSecret(this, true);
  navToDetail(this.secret);
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert.equal(find('[data-test-lease-renew-picker]').length, 1, 'renewable lease renders a renew picker');
  });
});

test('it removes leases upon revocation', function(assert) {
  createSecret(this);
  navToDetail(this.secret);
  click('[data-test-lease-revoke] button');
  click('[data-test-confirm-button]');
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'it navigates back to the leases root on revocation'
    );
  });
  click('[data-test-lease-link="secret/"]');
  andThen(() => {
    assert.equal(
      find(`[data-test-lease-link="secret/${this.secret.name}/"]`).length,
      0,
      'link to the lease was removed with revocation'
    );
  });
});

test('it removes branches when a prefix is revoked', function(assert) {
  createSecret(this);
  visit('/vault/access/leases/list/secret/');
  click('[data-test-lease-revoke-prefix] button');
  click('[data-test-confirm-button]');
  andThen(() => {
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'it navigates back to the leases root on revocation'
    );
    assert.equal(
      find('[data-test-lease-link="secret/"]').length,
      0,
      'link to the prefix was removed with revocation'
    );
  });
});

test('lease not found', function(assert) {
  visit('/vault/access/leases/show/not-found');
  andThen(() => {
    assert.equal(
      find('[data-test-lease-error]').text().trim(),
      'not-found is not a valid lease ID',
      'it shows an error when the lease is not found'
    );
  });
});
