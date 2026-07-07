/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';
import { visit } from '@ember/test-helpers';
import sinon from 'sinon';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { getErrorResponse } from 'vault/tests/helpers/api/error-response';
import { stubCapabilitiesForPaths } from 'vault/tests/helpers/capabilities/stubs';

// Tests logic in the pki.external.ts route because the overview route inherits the model from its parent route
module('Acceptance | enterprise | pki | external | overview route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.capabilities = this.owner.lookup('service:capabilities');
    this.acmeListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigAcmeAccount');
    this.dnsListStub = sinon.stub(api.secrets, 'pkiExternalCaListConfigDns');
    this.rolesListStub = sinon.stub(api.secrets, 'pkiExternalCaListRole');
    this.mountPath = `pki-external-ca-${uuidv4()}`;
    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.engineURL = `vault/secrets-engines/${this.mountPath}/pki/external/overview`;
  });

  test('it handles 404 errors gracefully', async function (assert) {
    stubCapabilitiesForPaths(
      this.capabilities,
      {
        pkiExternalConfigAcmeAccount: { canList: true },
        pkiExternalConfigDns: { canList: true },
        pkiExternalRole: { canList: true },
      },
      { backend: this.mountPath }
    );
    this.acmeListStub.resolves({ keys: ['account-1', 'account-2'] });
    this.dnsListStub.rejects(getErrorResponse());
    this.rolesListStub.rejects(getErrorResponse());
    await visit(this.engineURL);
    assert.true(this.acmeListStub.calledOnce, 'request made to list ACME accounts');
    assert.true(this.dnsListStub.calledOnce, 'request made to list DNS providers');
    assert.true(this.rolesListStub.calledOnce, 'request made to list roles');
    assert.dom(GENERAL.overviewCard.content('ACME accounts')).hasText('2');
    assert.dom(GENERAL.overviewCard.content('DNS providers')).hasText('0');
    assert.dom(GENERAL.overviewCard.content('Roles')).hasText('0');
  });

  test('it fetches data when user canList', async function (assert) {
    const capabilitiesStub = stubCapabilitiesForPaths(
      this.capabilities,
      {
        pkiExternalConfigAcmeAccount: { canList: true },
        pkiExternalConfigDns: { canList: true },
        pkiExternalRole: { canList: true },
      },
      { backend: this.mountPath }
    );
    this.acmeListStub.resolves({ keys: ['account-1', 'account-2'] });
    this.dnsListStub.resolves({ keys: ['provider-1'] });
    this.rolesListStub.resolves({ keys: ['role-1', 'role-2', 'role-3'] });

    await visit(this.engineURL);
    const [paths] = capabilitiesStub.lastCall.args;
    assert.propEqual(
      paths,
      [
        `${this.mountPath}/config/acme-account`,
        `${this.mountPath}/config/dns`,
        `${this.mountPath}/role`,
        `${this.mountPath}/lookup/orders`,
      ],
      'it requests capabilities service with expected paths'
    );
    assert.true(capabilitiesStub.calledOnce, 'capabilities are only requested once');
    assert.true(this.acmeListStub.calledOnce, 'request made to list ACME accounts');
    assert.true(this.dnsListStub.calledOnce, 'request made to list DNS providers');
    assert.true(this.rolesListStub.calledOnce, 'request made to list roles');
  });

  test('it fetches data when user canRead', async function (assert) {
    const capabilitiesStub = stubCapabilitiesForPaths(
      this.capabilities,
      {
        pkiExternalConfigAcmeAccount: { canRead: true },
        pkiExternalConfigDns: { canRead: true },
        pkiExternalRole: { canRead: true },
      },
      { backend: this.mountPath }
    );
    this.acmeListStub.resolves({ keys: ['account-1', 'account-2'] });
    this.dnsListStub.resolves({ keys: ['provider-1'] });
    this.rolesListStub.resolves({ keys: ['role-1', 'role-2', 'role-3'] });

    await visit(this.engineURL);
    const [paths] = capabilitiesStub.lastCall.args;
    assert.propEqual(
      paths,
      [
        `${this.mountPath}/config/acme-account`,
        `${this.mountPath}/config/dns`,
        `${this.mountPath}/role`,
        `${this.mountPath}/lookup/orders`,
      ],
      'it requests capabilities service with expected paths'
    );
    assert.true(capabilitiesStub.calledOnce, 'capabilities are only requested once');
    assert.true(this.acmeListStub.calledOnce, 'request made to list ACME accounts');
    assert.true(this.dnsListStub.calledOnce, 'request made to list DNS providers');
    assert.true(this.rolesListStub.calledOnce, 'request made to list roles');
  });

  test('it does NOT fetch data when user lacks permissions', async function (assert) {
    const capabilitiesStub = stubCapabilitiesForPaths(
      this.capabilities,
      {
        pkiExternalConfigAcmeAccount: { canList: false, canRead: false },
        pkiExternalConfigDns: { canList: false, canRead: false },
        pkiExternalRole: { canList: false, canRead: false },
      },
      { backend: this.mountPath }
    );
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.acmeListStub.rejects(getErrorResponse(error, 403));
    this.dnsListStub.rejects(getErrorResponse(error, 403));
    this.rolesListStub.rejects(getErrorResponse(error, 403));

    await visit(this.engineURL);
    const [paths] = capabilitiesStub.lastCall.args;
    assert.propEqual(
      paths,
      [
        `${this.mountPath}/config/acme-account`,
        `${this.mountPath}/config/dns`,
        `${this.mountPath}/role`,
        `${this.mountPath}/lookup/orders`,
      ],
      'it requests capabilities service with expected paths'
    );
    assert.true(capabilitiesStub.calledOnce, 'capabilities are only requested once');
    assert.true(this.acmeListStub.notCalled, 'request NOT made to list ACME accounts');
    assert.true(this.dnsListStub.notCalled, 'request NOT made to list DNS providers');
    assert.true(this.rolesListStub.notCalled, 'request NOT made to list roles');
  });

  test('it catches and displays non-404 error messages', async function (assert) {
    stubCapabilitiesForPaths(
      this.capabilities,
      {
        pkiExternalConfigAcmeAccount: { canList: true, canRead: false },
        pkiExternalConfigDns: { canList: true, canRead: false },
        pkiExternalRole: { canList: true, canRead: false },
      },
      { backend: this.mountPath }
    );

    const acmeError = { errors: ['Internal server error for ACME accounts'] };
    const dnsError = { errors: ['DNS provider service unavailable'] };
    this.acmeListStub.rejects(getErrorResponse(acmeError, 500));
    this.dnsListStub.rejects(getErrorResponse(dnsError, 503));
    this.rolesListStub.resolves({ keys: ['role-1'] });

    await visit(this.engineURL);
    assert
      .dom(`${GENERAL.overviewCard.container('ACME accounts')} ${GENERAL.messageError}`)
      .exists()
      .hasText('Error Internal server error for ACME accounts');
    assert.dom(GENERAL.overviewCard.content('ACME accounts')).doesNotExist();
    assert
      .dom(`${GENERAL.overviewCard.container('DNS providers')} ${GENERAL.messageError}`)
      .exists()
      .hasText('Error DNS provider service unavailable');
    assert.dom(GENERAL.overviewCard.content('DNS providers')).doesNotExist();
    // Roles still render just fine
    assert.dom(GENERAL.overviewCard.content('Roles')).hasText('1');
  });
});
