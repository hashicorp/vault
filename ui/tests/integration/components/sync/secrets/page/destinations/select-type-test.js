/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { click, render } from '@ember/test-helpers';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';
import { destinationTypes } from 'vault/helpers/sync-destinations';

module('Integration | Component | sync | page | destinations | select-type', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');

  test('it transitions to selected type', async function (assert) {
    const types = destinationTypes();
    assert.expect(types.length);
    const transitionStub = sinon.stub(this.owner.lookup('service:router'), 'transitionTo');

    await render(
      hbs`
     <Secrets::Page::Destinations::SelectType />
    `,
      { owner: this.engine }
    );

    for (const type of types) {
      await click(PAGE.create.selectType(type));
      const transition = transitionStub.calledWith(
        'vault.cluster.sync.secrets.destinations.create.destination',
        type
      );
      assert.true(transition, `transitionTo called with param: ${type}`);
    }
  });
});
