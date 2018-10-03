import { click, currentRouteName, visit } from '@ember/test-helpers';
// TESTS HERE ARE SKPPED
// running vault with -dev-leased-kv flag lets you run some of these tests
// but generating leases programmatically is currently difficult
//
// TODO revisit this when it's easier to create leases

import { module, skip } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import secretList from 'vault/tests/pages/secrets/backend/list';
import secretEdit from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';

module('Acceptance | leases', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await authPage.login();
    this.enginePath = `kv-for-lease-${new Date().getTime()}`;
    // need a version 1 mount for leased secrets here
    return mountSecrets
      .visit()
      .path(this.enginePath)
      .type('kv')
      .version(1)
      .submit();
  });

  hooks.afterEach(function() {
    return logout.visit();
  });

  const createSecret = async (context, isRenewable) => {
    context.name = `secret-${new Date().getTime()}`;
    await secretList.visitRoot({ backend: context.enginePath });
    await secretList.create();
    if (isRenewable) {
      await secretEdit.createSecret(context.name, 'ttl', '30h');
    } else {
      await secretEdit.createSecret(context.name, 'foo', 'bar');
    }
  };

  const navToDetail = async context => {
    await visit('/vault/access/leases/');
    await click(`[data-test-lease-link="${context.enginePath}/"]`);
    await click(`[data-test-lease-link="${context.enginePath}/${context.name}/"]`);
    await click(`[data-test-lease-link]:eq(0)`);
  };

  skip('it renders the show page', function(assert) {
    createSecret(this);
    navToDetail(this);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert
      .dom('[data-test-lease-renew-picker]')
      .doesNotExist('non-renewable lease does not render a renew picker');
  });

  // skip for now until we find an easy way to generate a renewable lease
  skip('it renders the show page with a picker', function(assert) {
    createSecret(this, true);
    navToDetail(this);
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert
      .dom('[data-test-lease-renew-picker]')
      .exists({ count: 1 }, 'renewable lease renders a renew picker');
  });

  skip('it removes leases upon revocation', async function(assert) {
    createSecret(this);
    navToDetail(this);
    await click('[data-test-lease-revoke] button');
    await click('[data-test-confirm-button]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'it navigates back to the leases root on revocation'
    );
    await click(`[data-test-lease-link="${this.enginePath}/"]`);
    await click('[data-test-lease-link="data/"]');
    assert
      .dom(`[data-test-lease-link="${this.enginePath}/data/${this.name}/"]`)
      .doesNotExist('link to the lease was removed with revocation');
  });

  skip('it removes branches when a prefix is revoked', async function(assert) {
    createSecret(this);
    await visit(`/vault/access/leases/list/${this.enginePath}`);
    await click('[data-test-lease-revoke-prefix] button');
    await click('[data-test-confirm-button]');
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'it navigates back to the leases root on revocation'
    );
    assert
      .dom(`[data-test-lease-link="${this.enginePath}/"]`)
      .doesNotExist('link to the prefix was removed with revocation');
  });

  skip('lease not found', async function(assert) {
    await visit('/vault/access/leases/show/not-found');
    assert
      .dom('[data-test-lease-error]')
      .hasText('not-found is not a valid lease ID', 'it shows an error when the lease is not found');
  });
});
