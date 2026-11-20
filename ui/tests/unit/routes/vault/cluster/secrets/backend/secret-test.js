/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster/secrets/backend secret', function (hooks) {
  setupTest(hooks);

  test('it fetches generic-secret capabilities in strict mode', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.show');
    const fetch = sinon.stub(route.capabilitiesService, 'fetch').resolves({
      'generic-secret/readable': { canRead: true, canUpdate: false },
    });

    sinon.stub(route, 'paramsFor').returns({ backend: 'generic-secret' });
    sinon.stub(route, 'modelFor').returns({ engineType: 'generic' });

    const actual = await route.capabilities('readable', 'secret');

    assert.deepEqual(
      actual,
      { canRead: true, canUpdate: false },
      'the route returns the requested capabilities'
    );
    assert.deepEqual(
      fetch.firstCall.args,
      [['generic-secret/readable'], { strict: true }],
      'the route requests capabilities for the full path and rejects permissive fallbacks'
    );
  });

  test('it permits write-without-read only when update capability is present', async function (assert) {
    const route = this.owner.lookup('route:vault.cluster.secrets.backend.show');
    const error = new Error('capability lookup failed');
    const model = await route.handleSecretModelError(
      { canUpdate: true, canDelete: false },
      'generic-secret/readable',
      'secret',
      error,
      'generic-secret'
    );

    assert.deepEqual(
      model,
      {
        id: 'generic-secret/readable',
        type: 'secret',
        backend: 'generic-secret',
        canDelete: false,
        failedServerRead: true,
        secretData: {},
      },
      'the fallback contains no secret data'
    );
    await assert.rejects(
      route.handleSecretModelError(
        { canUpdate: false },
        'generic-secret/readable',
        'secret',
        error,
        'generic-secret'
      ),
      (actual) => actual === error && actual.httpStatus === 403,
      'a user without update capability receives the read error'
    );
  });
});
