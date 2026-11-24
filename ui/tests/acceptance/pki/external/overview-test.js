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
import { assertTabState, EXTERNAL_TABS } from 'vault/tests/helpers/pki/assertion-helpers';

// Tests logic in the pki.external.ts route because the overview route inherits the model from its parent route
module('Acceptance | enterprise | pki | external | overview route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
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

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it navigates to external overview', async function (assert) {
    await visit(`vault/secrets-engines/${this.mountPath}/pki/external`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it redirects from the index route to overview'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'navigating to pki.external.index redirects to overview'
    );
    assert.dom(GENERAL.hdsPageHeaderTitle).exists().hasText(this.mountPath);
    assert.dom(GENERAL.breadcrumb).exists({ count: 3 });
    assert.dom(GENERAL.breadcrumbs).hasText(`Vault Secrets engines ${this.mountPath}`);
    assertTabState(assert, 'Overview', EXTERNAL_TABS);
  });

  test('only "Overview" tab renders when no resources exist but user has permission to list everything', async function (assert) {
    // getErrorResponse() throws 404 by default
    this.acmeListStub.rejects(getErrorResponse());
    this.dnsListStub.rejects(getErrorResponse());
    this.rolesListStub.rejects(getErrorResponse());
    await visit(`/vault/secrets-engines/${this.mountPath}/pki/external/overview`);
    assert.strictEqual(
      currentURL(),
      `/vault/secrets-engines/${this.mountPath}/pki/external/overview`,
      'it navigates to url'
    );
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.overview',
      'it navigates to route'
    );
    assert.dom(GENERAL.linkTo('Overview')).exists().hasClass('active');
    const hidden = EXTERNAL_TABS.filter((t) => t !== 'Overview');
    hidden.forEach((t) => assert.dom(GENERAL.linkTo(t)).doesNotExist());
  });

  test('it catches 404 errors', async function (assert) {
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

  test('it catches 403 permissions errors and hides cards', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.acmeListStub.rejects(getErrorResponse(error, 403));
    this.dnsListStub.rejects(getErrorResponse(error, 403));
    this.rolesListStub.rejects(getErrorResponse(error, 403));

    await visit(this.engineURL);
    assert.dom(GENERAL.overviewCard.container('ACME accounts')).doesNotExist();
    assert.dom(GENERAL.overviewCard.container('DNS providers')).doesNotExist();
    assert.dom(GENERAL.overviewCard.container('Roles')).doesNotExist();
    // Only order and cert lookups render
    assert.dom(GENERAL.overviewCard.container('View certificate')).exists();
    assert.dom(GENERAL.overviewCard.container('Orders')).exists();
  });

  test('it catches and displays non-404/non-403 error messages', async function (assert) {
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
