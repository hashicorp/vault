/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';

module('Acceptance | enterprise | pki | external | acme-accounts route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.acmeAccountsListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigAcmeAccount');
    this.acmeAccountReadStub = sinon.stub(api.secrets, 'pkiExternalCaReadConfigAcmeAccount');
    this.mountPath = `pki-external-ca-${uuidv4()}`;

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.acmeAccountsURL = `/vault/secrets-engines/${this.mountPath}/pki/external/acme-accounts`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it fetches ACME account details for each account', async function (assert) {
    this.acmeAccountsListStub.resolves({ keys: ['account-1', 'account-2', 'account-3'] });

    // Mock individual account reads
    this.acmeAccountReadStub.withArgs('account-1', this.mountPath).resolves({
      name: 'account-1',
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      email_contacts: ['admin1@example.com'],
      active_key_version: 0,
      account_keys: {
        0: { key_version: 0, key_type: 'ec-256', key_creation_date: '2024-01-15T10:30:00Z' },
      },
    });

    this.acmeAccountReadStub.withArgs('account-2', this.mountPath).resolves({
      name: 'account-2',
      directory_url: 'https://acme-staging-v02.api.letsencrypt.org/directory',
      email_contacts: ['admin2@example.com'],
      active_key_version: 0,
      account_keys: {
        0: { key_version: 0, key_type: 'rsa-2048', key_creation_date: '2024-01-16T10:30:00Z' },
      },
    });

    this.acmeAccountReadStub.withArgs('account-3', this.mountPath).resolves({
      name: 'account-3',
      directory_url: 'https://acme.example.com/directory',
      email_contacts: ['admin3@example.com'],
      active_key_version: 1,
      account_keys: {
        0: { key_version: 0, key_type: 'ec-384', key_creation_date: '2024-01-17T10:30:00Z' },
        1: { key_version: 1, key_type: 'ec-256', key_creation_date: '2024-01-18T10:30:00Z' },
      },
    });

    await visit(this.acmeAccountsURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.acme-accounts',
      'it navigates to acme-accounts route'
    );
    assert.true(this.acmeAccountsListStub.calledOnce, 'accounts list called once in parent route');
    assert.strictEqual(this.acmeAccountReadStub.callCount, 3, 'read called once for each account');
    assert.dom(GENERAL.cardContainer()).exists({ count: 3 });
    const accounts = ['account-1', 'account-2', 'account-3'];
    accounts.forEach((a, idx) => {
      assert.dom(`${GENERAL.cardContainer(a)} ${GENERAL.textDisplay()}`).hasText(a);
      assert
        .dom(`${GENERAL.cardContainer(a)} ${GENERAL.infoRowValue('Email contacts')}`)
        .hasText(`admin${idx + 1}@example.com`);
    });
  });

  test('it handles a 404', async function (assert) {
    this.acmeAccountsListStub.rejects(getErrorResponse()); // Throws 404
    await visit(this.acmeAccountsURL);
    assert.dom(GENERAL.emptyStateTitle).hasText('No ACME accounts exist yet');
    // Implementation select should be visible
    assert.dom(GENERAL.radioCardByAttr()).exists({ count: 2 }, 'it renders automation snippets');
    assert.dom(GENERAL.textDisplay()).exists({ count: 1 }).hasText('Configure an ACME account');
  });

  test('it displays 403 permission denied error for list response', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.acmeAccountsListStub.rejects(getErrorResponse(error, 403));
    await visit(this.acmeAccountsURL);
    assert.true(this.acmeAccountsListStub.calledOnce, 'accounts list called once');
    assert.strictEqual(currentURL(), this.acmeAccountsURL, 'it renders acme-accounts URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });

  test('it displays 500 internal server error', async function (assert) {
    const error = { errors: ['Internal server error'] };
    this.acmeAccountsListStub.rejects(getErrorResponse(error, 500));
    await visit(this.acmeAccountsURL);
    assert.strictEqual(currentURL(), this.acmeAccountsURL, 'it renders acme-accounts URL');
    assert.true(this.acmeAccountsListStub.calledOnce, 'accounts list called once');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(500)).exists().hasText('ERROR 500 Error');
  });

  test('it handles partial failures when reading individual accounts', async function (assert) {
    this.acmeAccountsListStub.resolves({ keys: ['account-1', 'account-2', 'account-3'] });
    // First account fails with 404
    this.acmeAccountReadStub.withArgs('account-1', this.mountPath).rejects(getErrorResponse());

    // Second account succeeds
    this.acmeAccountReadStub.withArgs('account-2', this.mountPath).resolves({
      name: 'account-2',
      directory_url: 'https://acme-v02.api.letsencrypt.org/directory',
      email_contacts: ['admin@example.com'],
      active_key_version: 1,
      account_keys: {
        1: { key_version: 1, key_type: 'ec-256', key_creation_date: '2024-01-15T10:30:00Z' },
      },
    });

    // Third account fails with 403
    this.acmeAccountReadStub
      .withArgs('account-3', this.mountPath)
      .rejects(getErrorResponse({ errors: ['1 error occurred:\n\t* permission denied\n\n'] }, 403));

    await visit(this.acmeAccountsURL);
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.acme-accounts',
      'stays on acme-accounts route despite partial failure'
    );
    // Should still render account name for each
    assert.dom(GENERAL.cardContainer()).exists({ count: 3 });
    const accounts = ['account-1', 'account-2', 'account-3'];
    accounts.forEach((a) => {
      assert.dom(`${GENERAL.cardContainer(a)} ${GENERAL.textDisplay()}`).hasText(a);
    });
    assert.strictEqual(this.acmeAccountReadStub.callCount, 3, 'read called for all accounts');

    // Successful account should be displayed
    assert
      .dom(`${GENERAL.cardContainer('account-2')} ${GENERAL.infoRowValue('Directory URL')}`)
      .hasText('https://acme-v02.api.letsencrypt.org/directory');

    // Failed accounts should show errors
    assert
      .dom(`${GENERAL.cardContainer('account-1')} ${GENERAL.infoRowValue('Error')}`)
      .hasText('An error occurred, please try again');
    assert
      .dom(`${GENERAL.cardContainer('account-3')} ${GENERAL.infoRowValue('Error')}`)
      .hasText('You do not have permission to read configurations for this account');
  });
});
