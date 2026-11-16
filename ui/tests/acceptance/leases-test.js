/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, visit } from '@ember/test-helpers';
import { module, skip, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import sinon from 'sinon';
import { v4 as uuidv4 } from 'uuid';
import secretList from 'vault/tests/pages/secrets/backend/list';
import secretEdit from 'vault/tests/pages/secrets/backend/kv/edit-secret';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// Lease lifecycle tests remain skipped because generating leases needs -dev-leased-kv.

module('Acceptance | leases', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  const createSecret = async (context, isRenewable) => {
    context.name = `secret-${uuidv4()}`;
    await secretList.visitRoot({ backend: context.enginePath });
    await click(SES.createSecretLink);
    if (isRenewable) {
      await secretEdit.createSecret(context.name, 'ttl', '30h');
    } else {
      await secretEdit.createSecret(context.name, 'foo', 'bar');
    }
  };

  const navToDetail = async (context) => {
    await visit('/vault/access/leases/');
    await click(`[data-test-lease-link="${context.enginePath}/"]`);
    await click(`[data-test-lease-link="${context.enginePath}/${context.name}/"]`);
    await click(`[data-test-lease-link]:eq(0)`);
  };

  skip('it renders the show page', function (assert) {
    createSecret(this);
    navToDetail(this);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert
      .dom('[data-test-lease-renew-picker]')
      .doesNotExist('non-renewable lease does not render a renew picker');
  });

  // skip for now until we find an easy way to generate a renewable lease
  skip('it renders the show page with a picker', function (assert) {
    createSecret(this, true);
    navToDetail(this);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.leases.show',
      'a lease for the secret is in the list'
    );
    assert
      .dom('[data-test-lease-renew-picker]')
      .exists({ count: 1 }, 'renewable lease renders a renew picker');
  });

  skip('it removes leases upon revocation', async function (assert) {
    createSecret(this);
    navToDetail(this);
    await click('[data-test-lease-revoke] button');
    await click(GENERAL.confirmButton);
    assert.strictEqual(
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

  skip('it removes branches when a prefix is revoked', async function (assert) {
    createSecret(this);
    await visit(`/vault/access/leases/list/${this.enginePath}`);
    await click('[data-test-lease-revoke-prefix] button');
    await click(GENERAL.confirmButton);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.access.leases.list-root',
      'it navigates back to the leases root on revocation'
    );
    assert
      .dom(`[data-test-lease-link="${this.enginePath}/"]`)
      .doesNotExist('link to the prefix was removed with revocation');
  });

  test('lease not found', async function (assert) {
    await login();
    const api = this.owner.lookup('service:api');
    const capabilities = this.owner.lookup('service:capabilities');
    sinon.stub(api.sys, 'leasesReadLease').rejects(getErrorResponse({ errors: ['invalid lease'] }, 400));
    sinon.stub(capabilities, 'fetchPathCapabilities').resolves({
      canList: true,
      canRead: true,
      canUpdate: true,
    });

    await visit('/vault/access/leases/show/not-found');
    assert
      .dom('[data-test-lease-error]')
      .hasText('not-found is not a valid lease ID.', 'it explains that the lease ID is invalid');
  });

  test('lease lookup 404 renders the no lease state', async function (assert) {
    await login();
    const api = this.owner.lookup('service:api');
    const capabilities = this.owner.lookup('service:capabilities');
    sinon.stub(api.sys, 'leasesReadLease').rejects(getErrorResponse({}, 404));
    sinon.stub(capabilities, 'fetchPathCapabilities').resolves({
      canList: true,
      canRead: true,
      canUpdate: true,
    });

    await visit('/vault/access/leases/show/not-found');
    assert
      .dom('[data-test-empty-state-title]')
      .hasText('No leases with that ID', 'it explains that no lease exists for the ID');
  });
});
