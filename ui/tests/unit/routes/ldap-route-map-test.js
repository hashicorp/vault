/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

const PREFIX = 'vault.cluster.secrets.backend.ldap';
const BASE = '/vault/secrets-engines/my-ldap/ldap';

/**
 * Pins the URL to route name mapping of the LDAP routes through the public router service so the
 * map stays identical when the engine is inlined into the main app.
 */
module('Unit | Route | ldap route map', function (hooks) {
  setupTest(hooks);

  test('it recognizes urls as the expected route names and params', function (assert) {
    const router = this.owner.lookup('service:router');
    const cases = [
      ['overview', 'overview', { backend: 'my-ldap' }],
      ['configure', 'configure', { backend: 'my-ldap' }],
      ['configuration', 'configuration', { backend: 'my-ldap' }],
      ['roles', 'roles.index', { backend: 'my-ldap' }],
      ['roles/create', 'roles.create', { backend: 'my-ldap' }],
      ['roles/static/my-role/details', 'roles.role.details', { type: 'static', name: 'my-role' }],
      ['roles/dynamic/my-role/edit', 'roles.role.edit', { type: 'dynamic', name: 'my-role' }],
      ['roles/static/my-role/credentials', 'roles.role.credentials', { type: 'static', name: 'my-role' }],
      ['roles/static/subdirectory/admin/', 'roles.subdirectory', { type: 'static', path_to_role: 'admin/' }],
      ['libraries', 'libraries.index', { backend: 'my-ldap' }],
      ['libraries/create', 'libraries.create', { backend: 'my-ldap' }],
      ['libraries/subdirectory/admin/', 'libraries.subdirectory', { path_to_library: 'admin/' }],
      ['libraries/my-lib/details/accounts', 'libraries.library.details.accounts', { name: 'my-lib' }],
      [
        'libraries/my-lib/details/configuration',
        'libraries.library.details.configuration',
        { name: 'my-lib' },
      ],
      ['libraries/my-lib/edit', 'libraries.library.edit', { name: 'my-lib' }],
      ['libraries/my-lib/check-out', 'libraries.library.check-out', { name: 'my-lib' }],
    ];

    for (const [path, routeName, params] of cases) {
      const info = router.recognize(`/ui${BASE}/${path}`);
      assert.strictEqual(info?.name, `${PREFIX}.${routeName}`, `${path} resolves to ${routeName}`);
      let current = info;
      const allParams = {};
      while (current) {
        Object.assign(allParams, current.params);
        current = current.parent;
      }
      assert.propContains(allParams, params, `${path} has the expected params`);
    }
  });
});
