/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { fillIn, waitFor } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { deleteNS, runCmd } from 'vault/tests/helpers/commands';
import { setupMirage } from 'ember-cli-mirage/test-support';
import parseURL from 'core/utils/parse-url';
import { login, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { AUTH_FORM } from 'vault/tests/helpers/auth/auth-form-selectors';

module('Acceptance | Enterprise | oidc auth namespace test', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  test('oidc: request is made to auth_url when a namespace is inputted', async function (assert) {
    assert.expect(2);
    const namespace = 'test-ns';
    const mountPath = 'ns-oidc';

    // stubs the auth_url for the OIDC method configured in the namespace, NOT for the root namespace
    // should only be hit once when the oidc tab is selected after inputting a namespace
    this.server.post(`/auth/${mountPath}/oidc/auth_url`, (schema, req) => {
      const { redirect_uri } = JSON.parse(req.requestBody);
      const { pathname, search } = parseURL(redirect_uri);
      assert.strictEqual(
        pathname + search,
        `/ui/vault/auth/${mountPath}/oidc/callback?namespace=${namespace}`,
        'request made to correct auth_url when namespace is filled in'
      );
    });

    await login();
    await runCmd([
      `write sys/namespaces/${namespace} -force`,
      `write ${namespace}/sys/auth/${mountPath} type=oidc`,
      `write auth/${mountPath}/config default_role="myrole" oidc_discovery_url="https://example.com"`,
      // show method as tab
      `write ${namespace}/sys/auth/${mountPath}/tune listing_visibility="unauth"`,
    ]);
    await logout();

    await fillIn(GENERAL.inputByAttr('namespace'), namespace);
    await waitFor(AUTH_FORM.tabBtn('oidc')); // no need to click because selected by default
    assert.dom(AUTH_FORM.tabBtn('oidc')).exists('renders oidc method tab for child namespace');

    // cleanup
    await fillIn(GENERAL.inputByAttr('namespace'), ''); // clear namespace input
    await login();
    await runCmd(deleteNS(namespace));
  });
});
