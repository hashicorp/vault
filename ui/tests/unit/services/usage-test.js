/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import sinon from 'sinon';

module('Unit | Service | usage', function (hooks) {
  setupTest(hooks);
  hooks.beforeEach(function () {
    this.originalFetch = window.fetch;
    this.fetchJsonStub = sinon.stub().resolves({ data: {} });
    this.fetchStub = sinon.stub(window, 'fetch').resolves({
      json: this.fetchJsonStub,
      ok: true,
    });
    this.auth = this.owner.lookup('service:auth');
    this.currentTokenStub = sinon.stub(this.auth, 'currentToken').value('abc123');
    this.service = this.owner.lookup('service:usage');
  });

  hooks.afterEach(function () {
    this.fetchStub.restore();
    this.currentTokenStub.restore();
    window.fetch = this.originalFetch;
  });

  test('it calls fetch with correct endpoint and headers', async function (assert) {
    await this.service.getUsageData();
    assert.ok(
      this.fetchStub.calledOnceWith(
        '/v1/sys/utilization-report',
        sinon.match({ headers: { 'X-Vault-Token': 'abc123' } })
      )
    );
  });

  test('it throws an error if fetch is not ok', async function (assert) {
    this.fetchStub.resolves({ ok: false });
    assert.rejects(this.service.getUsageData(), /Failed to fetch usage data/);
  });
});
