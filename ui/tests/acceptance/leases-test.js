/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
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
import { TTL_PICKER } from 'vault/tests/helpers/components/ttl-picker-selectors';

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

  module('renewing a lease', function (hooks) {
    hooks.beforeEach(async function () {
      await login();
      const api = this.owner.lookup('service:api');
      const capabilities = this.owner.lookup('service:capabilities');
      this.leaseId = 'secret/renewable/abcd1234';
      sinon.stub(api.sys, 'leasesReadLease').resolves({
        id: this.leaseId,
        renewable: true,
        issue_time: '2018-04-03T15:15:30.000Z',
        expire_time: '2018-04-04T15:15:30.000Z',
        last_renewal: null,
        ttl: 3600,
      });
      sinon.stub(capabilities, 'fetchPathCapabilities').resolves({
        canList: true,
        canRead: true,
        canUpdate: true,
        canDelete: true,
      });
      this.renewStub = sinon.stub(api.sys, 'leasesRenewLease').resolves({});
      this.visitLease = () => visit(`/vault/access/leases/show/${this.leaseId}`);
      this.enterIncrement = async (value) => {
        await click(TTL_PICKER.toggleByLabel('Increment'));
        await fillIn(TTL_PICKER.valueInputByLabel('Increment'), value);
        await fillIn(TTL_PICKER.unitInputByLabel('Increment'), 'd');
      };
    });

    test('it blocks an increment longer than Vault supports', async function (assert) {
      assert.expect(5);
      await this.visitLease();
      await this.enterIncrement('9999999999');
      await click('[data-test-renew-lease-submit]');

      assert.dom(TTL_PICKER.error).hasText('Increment is too large. Maximum allowed value is 106,751 days.');
      assert.true(this.renewStub.notCalled, 'the lease is not renewed');

      await fillIn(TTL_PICKER.valueInputByLabel('Increment'), '106751');
      await click('[data-test-renew-lease-submit]');
      assert.true(this.renewStub.calledOnce, 'it renews once the increment is valid');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '9223286400');
      assert.strictEqual(this.renewStub.firstCall.args[0].lease_id, this.leaseId);
    });

    test('it blocks an increment too long to be stored as a number', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await this.enterIncrement('9'.repeat(400));
      await click('[data-test-renew-lease-submit]');

      assert.dom(TTL_PICKER.error).exists('the picker shows its error');
      assert.true(this.renewStub.notCalled, 'the lease is not renewed');
    });

    test('it does not send an increment that was turned off', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await this.enterIncrement('9999999999');
      await click(TTL_PICKER.toggleByLabel('Increment'));
      await click('[data-test-renew-lease-submit]');

      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '', 'the hidden increment is not sent');
    });

    // Existing behaviour, kept as it was before the max-duration checks.
    test('it keeps sending a valid increment that was turned off', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await click(TTL_PICKER.toggleByLabel('Increment'));
      await fillIn(TTL_PICKER.valueInputByLabel('Increment'), '2');
      await fillIn(TTL_PICKER.unitInputByLabel('Increment'), 'h');
      await click(TTL_PICKER.toggleByLabel('Increment'));
      await click('[data-test-renew-lease-submit]');

      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(
        this.renewStub.firstCall.args[0].increment,
        '7200',
        'the entered increment is still sent'
      );
    });

    test('it renews with the entered increment', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await click(TTL_PICKER.toggleByLabel('Increment'));
      await fillIn(TTL_PICKER.valueInputByLabel('Increment'), '2');
      await fillIn(TTL_PICKER.unitInputByLabel('Increment'), 'h');
      await click('[data-test-renew-lease-submit]');

      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '7200', 'the entered increment is sent');
    });

    test('it renews with no increment when the picker is untouched', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await click('[data-test-renew-lease-submit]');

      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '', 'no increment is sent');
    });

    // The controller outlives the route, so an increment from an earlier visit must not carry over.
    test('it renews with the default after leaving the page with an over-limit increment', async function (assert) {
      assert.expect(3);
      await this.visitLease();
      await this.enterIncrement('9999999999');
      await visit('/vault/access/leases/list');
      await this.visitLease();
      assert.dom(TTL_PICKER.error).doesNotExist('the new picker shows no error');

      await click('[data-test-renew-lease-submit]');
      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '', 'no increment is sent');
    });

    test('it does not reuse a valid increment from an earlier visit', async function (assert) {
      assert.expect(2);
      await this.visitLease();
      await this.enterIncrement('2');
      await visit('/vault/access/leases/list');
      await this.visitLease();
      await click('[data-test-renew-lease-submit]');

      assert.true(this.renewStub.calledOnce, 'the lease is renewed');
      assert.strictEqual(this.renewStub.firstCall.args[0].increment, '', 'the picker is untouched');
    });
  });
});
