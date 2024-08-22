/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { setupMirage } from 'ember-cli-mirage/test-support';
import { module, test } from 'qunit';
import Sinon from 'sinon';

import { setupTest } from 'vault/tests/helpers';

module('Unit | Adapter | application', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.server.get('/some-url', function () {
      return {
        warnings: ['this is a warning'],
      };
    });
    this.adapter = this.owner.lookup('adapter:application');
  });

  test('it triggers info flash message when warnings returned from API', async function (assert) {
    const flashSuccessSpy = Sinon.spy(this.owner.lookup('service:flash-messages'), 'info');
    await this.adapter.ajax('/v1/some-url', 'GET', { skipWarnings: true });
    assert.true(flashSuccessSpy.notCalled, 'flash is not called when skipWarnings option passed');
    await this.adapter.ajax('/v1/some-url', 'GET', {});
    assert.true(flashSuccessSpy.calledOnce);
    assert.true(flashSuccessSpy.calledWith('this is a warning'));
  });
});
