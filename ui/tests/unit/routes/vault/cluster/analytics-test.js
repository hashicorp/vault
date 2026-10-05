/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Route | vault/cluster | analytics identify', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    this.route = this.owner.lookup('route:vault/cluster');
    this.analytics = this.owner.lookup('service:analytics');
    this.version = this.owner.lookup('service:version');
    this.flags = this.owner.lookup('service:flags');

    // Pretend analytics is already consented/started so identify runs.
    this.analytics.activated = true;
    this.identifyStub = sinon.stub(this.analytics, 'identifyUser');

    // addAnalyticsService reads the license status over the API; stub it out.
    sinon.stub(this.route.api.sys, 'systemReadLicenseStatus').resolves({ data: { autoloaded: {} } });

    // Self-managed (consent-gated) path. startVaultSmAnalytics is already a no-op
    // here because activated is true, but stub it to avoid the network read.
    sinon.stub(this.flags, 'isHvdManaged').value(false);
    sinon.stub(this.route, 'startVaultSmAnalytics').resolves();

    this.model = {
      id: 'cluster-1',
      version: { version: '1.99.0+ent' },
      storageType: 'raft',
      replicationMode: 'dr',
      // license deliberately absent: it is loaded by a separate fetch and is
      // routinely null when addAnalyticsService runs.
      license: null,
    };
  });

  hooks.afterEach(function () {
    sinon.restore();
  });

  test('reports isEnterprise from the version service, not model.license', async function (assert) {
    sinon.stub(this.version, 'isEnterprise').value(true);

    await this.route.addAnalyticsService(this.model);

    assert.true(this.identifyStub.calledOnce, 'identifyUser is called');
    assert.true(
      this.identifyStub.firstCall.args[1].isEnterprise,
      'isEnterprise is true for an enterprise cluster even though model.license is null'
    );
  });

  test('reports isEnterprise false for a community cluster', async function (assert) {
    sinon.stub(this.version, 'isEnterprise').value(false);

    await this.route.addAnalyticsService(this.model);

    assert.true(this.identifyStub.calledOnce, 'identifyUser is called');
    assert.false(
      this.identifyStub.firstCall.args[1].isEnterprise,
      'isEnterprise is false when the version service reports community'
    );
  });
});
