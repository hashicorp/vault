/**
 * Copyright (c) HashiCorp, Inc.
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';

module('Integration | Component | usage | Page::Usage', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(async function () {
    this.api = this.owner.lookup('service:api');
    this.generateUtilizationReportStub = sinon.stub(this.api.sys, 'generateUtilizationReport').resolves({});
  });

  hooks.afterEach(function () {
    this.generateUtilizationReportStub.restore();
  });

  test('it provides the correct fetch function to the dashboard component', async function (assert) {
    await render(hbs`<Usage::Page />`);
    assert.true(this.generateUtilizationReportStub.calledOnce, 'fetch function is called on render');
  });
});
