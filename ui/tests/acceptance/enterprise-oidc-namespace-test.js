/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { currentURL, click, fillIn } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { runCmd, createNS } from 'vault/tests/helpers/commands';
import { setupMirage } from 'ember-cli-mirage/test-support';
import parseURL from 'core/utils/parse-url';
import { login, loginNs, logout } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  authTab: (path) => `[data-test-auth-method="${path}"] a`,
};

module('Acceptance | Enterprise | oidc auth namespace test', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.namespace = 'test-ns';
    this.rootOidc = 'root-oidc';
    this.nsOidc = 'ns-oidc';

    this.server.post(`/auth/:path/config`, () => {});

    this.enableOidc = (path, role = '') => {
      return runCmd([
        `write sys/auth/${path} type=oidc`,
        `write auth/${path}/config default_role="${role}" oidc_discovery_url="https://example.com"`,
        // show method as tab
        `write sys/auth/${path}/tune listing_visibility="unauth"`,
      ]);
    };

    this.disableOidc = (path) => runCmd([`delete /sys/auth/${path}`]);
  });

  test('oidc: request is made to auth_url when a namespace is inputted', async function (assert) {
    assert.expect(4);

    // stubs the auth_url for the OIDC method configured in the namespace, NOT for the root namespace
    // should only be hit once when the oidc tab is clicked after inputting a namespace
    this.server.post(`/auth/${this.nsOidc}/oidc/auth_url`, (schema, req) => {
      const { redirect_uri } = JSON.parse(req.requestBody);
      const { pathname, search } = parseURL(redirect_uri);
      assert.strictEqual(
        pathname + search,
        `/ui/vault/auth/${this.nsOidc}/oidc/callback?namespace=${this.namespace}`,
        'request made to correct auth_url when namespace is filled in'
      );
    });

    await login();
    // enable oidc in root namespace, without default role
    await this.enableOidc(this.rootOidc);
    // create child namespace to enable oidc
    await runCmd(createNS(this.namespace), false);
    // enable oidc in child namespace with default role
    await loginNs(this.namespace);
    await this.enableOidc(this.nsOidc, `${this.nsOidc}-role`);
    // check root namespace for method tab
    await logout();
    await fillIn(GENERAL.inputByAttr('namespace'), '');
    assert.dom(SELECTORS.authTab(this.rootOidc)).exists('renders oidc method tab for root');
    // check child namespace for method tab
    await fillIn(GENERAL.inputByAttr('namespace'), this.namespace);
    assert.dom(SELECTORS.authTab(this.nsOidc)).exists('renders oidc method tab for child namespace');
    // clicking on the tab should update with= queryParam
    await click(`[data-test-auth-method="${this.nsOidc}"] a`); // fires request to /auth_url and assertion
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${this.namespace}&with=${this.nsOidc}%2F`,
      'url updates with namespace value'
    );
    // disable methods to cleanup test state for re-running
    await fillIn(GENERAL.inputByAttr('namespace'), '');
    await login();
    await this.disableOidc(this.rootOidc);
    await this.disableOidc(this.nsOidc);
    await runCmd([`delete /sys/auth/${this.namespace}`]);
  });
});
