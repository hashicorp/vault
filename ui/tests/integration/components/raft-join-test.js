/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

module('Integration | Component | raft-join', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`<RaftJoin />`);
    assert.dom('[data-test-join-choice]').exists();
  });

  test('it shows the join form when clicking next', async function (assert) {
    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    assert.dom('[data-test-join-header]').exists();
  });
  test('it returns to the first screen when clicking back', async function (assert) {
    await render(hbs`<RaftJoin />`);
    await click('[data-test-next]');
    assert.dom('[data-test-join-header]').exists();
    await click('[data-test-cancel-button]');
    assert.dom('[data-test-join-choice]').exists();
  });

  test('it calls onDismiss when a user chooses to init', async function (assert) {
    const spy = sinon.spy();
    this.set('onDismiss', spy);
    await render(hbs`<RaftJoin @onDismiss={{this.onDismiss}} />`);

    await click('[data-test-join-init]');
    await click('[data-test-next]');
    assert.ok(spy.calledOnce, 'it calls the passed onDismiss');
  });
});
