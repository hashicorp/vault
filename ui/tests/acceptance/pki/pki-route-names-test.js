/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentRouteName, currentURL, visit } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';

import { login } from 'vault/tests/helpers/auth/auth-helpers';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import { runCmd } from 'vault/tests/helpers/commands';

const ROUTE_PREFIX = 'vault.cluster.secrets.backend.pki';

/**
 * Frozen baseline for the PKI engine route table (lib/pki/addon/routes.js). Internal routes are
 * visited by URL and the full route name is asserted, so inlining the engine into the main app must
 * keep every URL and route name unchanged. The configuration and cross-sign routes are covered by
 * pki-configuration-test.js and pki-cross-sign-test.js, which need their own mount setup.
 */
module('Acceptance | pki | route names', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    await login();
    this.mountPath = `pki-routes-${uuidv4()}`;
    await enablePage.enable('pki', this.mountPath);
    await runCmd([
      `write ${this.mountPath}/root/generate/internal common_name="Route Test Root" issuer_name="route-issuer" key_name="route-key"`,
      `write ${this.mountPath}/roles/route-role issuer_ref="default" allow_any_name=true max_ttl="720h"`,
      `write ${this.mountPath}/issue/route-role common_name="route.example.com"`,
    ]);
    const list = async (path) => {
      const res = await fetch(`/v1/${this.mountPath}/${path}?list=true`, {
        headers: { 'X-Vault-Token': 'root' },
      });
      return (await res.json()).data.keys;
    };
    [this.issuerId] = await list('issuers');
    [this.keyId] = await list('keys');
    [this.serial] = (await list('certs')).map((s) => s.replaceAll('-', ':'));
    await login();
    this.base = `/vault/secrets-engines/${this.mountPath}/pki`;
  });

  hooks.afterEach(async function () {
    await login();
    await runCmd([`delete sys/mounts/${this.mountPath}`]);
  });

  const staticRoutes = {
    '/overview': 'overview',
    '/roles': 'roles.index',
    '/roles/create': 'roles.create',
    '/issuers': 'issuers.index',
    '/issuers/import': 'issuers.import',
    '/issuers/generate-root': 'issuers.generate-root',
    '/issuers/generate-intermediate': 'issuers.generate-intermediate',
    '/keys': 'keys.index',
    '/keys/create': 'keys.create',
    '/keys/import': 'keys.import',
    '/certificates': 'certificates.index',
    '/tidy': 'tidy.index',
    '/tidy/auto/configure': 'tidy.auto.configure',
    '/tidy/manual': 'tidy.manual',
  };
  for (const [path, name] of Object.entries(staticRoutes)) {
    test(`static route ${path} resolves to ${name}`, async function (assert) {
      await visit(`${this.base}${path}`);
      assert.strictEqual(currentURL(), `${this.base}${path}`, 'it keeps its URL');
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.${name}`, 'it resolves the full route name');
    });
  }

  const dynamicRoutes = {
    '/roles/route-role/details': 'roles.role.details',
    '/roles/route-role/edit': 'roles.role.edit',
    '/roles/route-role/generate': 'roles.role.generate',
    '/roles/route-role/sign': 'roles.role.sign',
    '/issuers/:issuerId/details': 'issuers.issuer.details',
    '/issuers/:issuerId/edit': 'issuers.issuer.edit',
    '/issuers/:issuerId/sign': 'issuers.issuer.sign',
    '/issuers/:issuerId/rotate-root': 'issuers.issuer.rotate-root',
    '/keys/:keyId/details': 'keys.key.details',
    '/keys/:keyId/edit': 'keys.key.edit',
    '/certificates/:serial/details': 'certificates.certificate.details',
  };
  for (const [template, name] of Object.entries(dynamicRoutes)) {
    test(`dynamic route ${template} resolves to ${name}`, async function (assert) {
      const path = template
        .replace(':issuerId', this.issuerId)
        .replace(':keyId', this.keyId)
        .replace(':serial', this.serial);
      await visit(`${this.base}${path}`);
      assert.strictEqual(currentURL(), `${this.base}${path}`, 'it keeps its URL');
      assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.${name}`, 'it resolves the full route name');
    });
  }

  test('the engine index route redirects to overview', async function (assert) {
    await visit(this.base);
    assert.strictEqual(currentURL(), `${this.base}/overview`, 'it redirects to the overview URL');
    assert.strictEqual(currentRouteName(), `${ROUTE_PREFIX}.overview`, 'it lands on the overview route');
  });
});
