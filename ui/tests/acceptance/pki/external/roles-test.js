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

module('Acceptance | enterprise | pki | external | roles route', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    // Test setup
    const api = this.owner.lookup('service:api');
    this.rolesListStub = sinon.stub(api.secrets, 'pkiExternalCaListRole');
    this.mountPath = `pki-external-ca-${uuidv4()}`;

    // Setup External PKI engine
    await login();
    await runCmd(mountEngineCmd('pki-external-ca', this.mountPath));
    // assertion helpers
    this.rolesURL = `/vault/secrets-engines/${this.mountPath}/pki/external/roles`;
  });

  hooks.afterEach(async function () {
    // cleanup after
    await runCmd([`delete sys/mounts/${this.mountPath}`], false);
  });

  test('it uses parent route data and does not re-fetch', async function (assert) {
    this.rolesListStub.resolves({ keys: ['role-1', 'role-2', 'role-3'] });
    await visit(this.rolesURL);
    assert.strictEqual(currentURL(), this.rolesURL, 'it has expected URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.roles.index',
      'navigated to roles route'
    );
    assert.true(this.rolesListStub.calledOnce, 'roles list called once in parent route');

    // Verify roles are displayed
    assert.dom(GENERAL.listItem()).exists({ count: 3 }, 'displays all roles');
    assert.dom(GENERAL.linkTo('role-1')).exists();
    assert.dom(GENERAL.linkTo('role-2')).exists();
    assert.dom(GENERAL.linkTo('role-3')).exists();
  });

  test('it handles a 404', async function (assert) {
    this.rolesListStub.rejects(getErrorResponse()); // Throws 404
    await visit(this.rolesURL);
    assert.strictEqual(currentURL(), this.rolesURL, 'navigated to roles route');
    assert.dom(GENERAL.emptyStateTitle).hasText('No roles exist yet');
    // Implementation select should be visible
    assert.dom(GENERAL.radioCardByAttr()).exists({ count: 2 }, 'it renders automation snippets');
    assert.dom(GENERAL.textDisplay()).exists({ count: 1 }).hasText('Create a role');
  });

  test('it displays 403 permission denied error', async function (assert) {
    const error = { errors: ['1 error occurred:\n\t* permission denied\n\n'] };
    this.rolesListStub.rejects(getErrorResponse(error, 403));
    await visit(this.rolesURL);
    assert.true(this.rolesListStub.calledOnce, 'roles list called once');
    assert.strictEqual(currentURL(), this.rolesURL, 'it renders roles URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(403)).exists().hasText('ERROR 403 Not authorized');
  });

  test('it displays 500 internal server error', async function (assert) {
    const error = { errors: ['Internal server error'] };
    this.rolesListStub.rejects(getErrorResponse(error, 500));
    await visit(this.rolesURL);
    assert.true(this.rolesListStub.calledOnce, 'roles list called once');
    assert.strictEqual(currentURL(), this.rolesURL, 'it renders roles URL');
    assert.strictEqual(
      currentRouteName(),
      'vault.cluster.secrets.backend.pki.external.error',
      'it redirects to error route'
    );
    assert.dom(GENERAL.pageError.title(500)).exists().hasText('ERROR 500 Error');
  });
});
