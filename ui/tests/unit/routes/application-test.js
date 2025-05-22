/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'vault/tests/helpers';

import sinon from 'sinon';

module('Unit | Route | application', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const analytics = this.owner.lookup('service:analytics');
    this.analyticsStartStub = sinon.stub(analytics, 'start');
    this.flags = this.owner.lookup('service:flags');
  });

  hooks.afterEach(function () {
    sinon.reset();
  });

  test('it sets up the analytics service when a cluster is flagged as "HVD managed"', function (assert) {
    this.flags.featureFlags = ['VAULT_CLOUD_ADMIN_NAMESPACE'];
    const route = this.owner.lookup('route:application');

    route.afterModel();

    assert.true(
      this.analyticsStartStub.calledWith('posthog'),
      'the start call was made with the correct provider'
    );
  });

  test('it does not set up the analytics service when a cluster is not flagged as "HVD managed"', function (assert) {
    this.flags.featureFlags = ['AARDVARKS_ACTIVATED'];
    const route = this.owner.lookup('route:application');

    route.afterModel();

    assert.true(this.analyticsStartStub.notCalled, 'the start call was not made');
  });
});
