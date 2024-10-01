/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import EmberObject from '@ember/object';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | edit form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = EmberObject.create({
      fields: [
        { name: 'one', type: 'string' },
        { name: 'two', type: 'boolean' },
      ],
      destroyRecord() {},
      save() {},
      rollbackAttributes() {},
    });
    this.onSave = sinon.spy();
    this.renderComponent = () =>
      render(hbs`
      <EditForm @model={{this.model}} @onSave={{this.onSave}} />
    `);
  });

  test('it renders', async function (assert) {
    await this.renderComponent();
    assert.dom(GENERAL.fieldByAttr('one')).exists();
    assert.dom(GENERAL.fieldByAttr('two')).exists();
  });

  test('it calls flash message fns on save', async function (assert) {
    assert.expect(4);
    const flash = this.owner.lookup('service:flash-messages');
    this.flashSuccessSpy = sinon.spy(flash, 'success');
    await this.renderComponent();
    await click('[data-test-edit-form-submit]');
    const { saveType, model } = this.onSave.lastCall.args[0];
    const [flashMessage] = this.flashSuccessSpy.lastCall.args;
    assert.strictEqual(flashMessage, 'Saved!');
    assert.strictEqual(saveType, 'save');
    assert.strictEqual(saveType, 'save');
    assert.deepEqual(model, this.model, 'passes model to onSave');
  });
});
