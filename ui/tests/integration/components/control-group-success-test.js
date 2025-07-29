/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, fillIn, find, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { CONTROL_GROUP } from 'vault/tests/helpers/components/control-group-selectors';

module('Integration | Component | control group success', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');
    this.controlGroup = this.owner.lookup('service:control-group');
    this.markTokenForUnwrapStub = sinon.stub(this.controlGroup, 'markTokenForUnwrap');
    this.model = {
      approved: false,
      requestPath: 'foo/bar',
      id: 'accessor',
      requestEntity: { id: 'requestor', name: 'entity8509' },
      reload: sinon.stub(),
    };
  });

  test('render with saved token', async function (assert) {
    assert.expect(3);
    const response = {
      uiParams: { url: '/foo' },
      token: 'token',
    };
    this.set('response', response);
    await render(hbs`<ControlGroupSuccess @model={{this.model}} @controlGroupResponse={{this.response}} />`);
    assert
      .dom(CONTROL_GROUP.navMessage)
      .hasText(
        'You have been granted access to foo/bar. Be careful, you can only access this data once. If you need access again in the future you will need to get authorized again. Visit'
      );

    await click(GENERAL.button('Visit'));
    const [transition] = this.transitionStub.lastCall.args;
    const [accessor] = this.markTokenForUnwrapStub.lastCall.args;
    assert.strictEqual(accessor, 'accessor', 'marks token for unwrap');
    assert.strictEqual(transition, '/foo', 'calls router transition');
  });

  test('render without token', async function (assert) {
    assert.expect(2);
    await render(hbs`<ControlGroupSuccess @model={{this.model}} />`);
    assert.dom(CONTROL_GROUP.unwrapForm).exists();
    assert.dom(GENERAL.inputByAttr('token')).hasValue('');
  });

  test('it unwraps data on submit', async function (assert) {
    assert.expect(2);

    sinon.stub(this.owner.lookup('service:api').sys, 'unwrap').resolves({ data: { foo: 'bar' } });

    await render(hbs`<ControlGroupSuccess @model={{this.model}} />`);
    assert.dom(GENERAL.inputByAttr('token')).hasValue('');

    await fillIn(GENERAL.inputByAttr('token'), 'token');
    await click(GENERAL.submitButton);

    const actual = find(CONTROL_GROUP.jsonViewer).innerText;
    const expected = JSON.stringify({ foo: 'bar' }, null, 2);
    assert.strictEqual(actual, expected, `it renders unwrapped data: ${actual}`);
  });
});
