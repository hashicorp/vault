/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, run, _cancelTimers as cancelTimers } from '@ember/runloop';
import { resolve } from 'rsvp';
import EmberObject from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { create } from 'ember-cli-page-object';
import editForm from 'vault/tests/pages/components/edit-form';

const component = create(editForm);

const flash = Service.extend({
  success: sinon.stub(),
});

const createModel = (canDelete = true) => {
  return EmberObject.create({
    fields: [
      { name: 'one', type: 'string' },
      { name: 'two', type: 'boolean' },
    ],
    canDelete,
    destroyRecord() {
      return resolve();
    },
    save() {
      return resolve();
    },
    rollbackAttributes() {},
  });
};

module('Integration | Component | edit form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    run(() => {
      this.owner.unregister('service:flash-messages');
      this.owner.register('service:flash-messages', flash);
    });
  });

  test('it renders', async function (assert) {
    this.set('model', createModel());
    await render(hbs`{{edit-form model=this.model}}`);

    assert.ok(component.fields.length, 2);
  });

  test('it calls flash message fns on save', async function (assert) {
    assert.expect(4);
    const onSave = () => {
      return resolve();
    };
    this.set('model', createModel());
    this.set('onSave', onSave);
    const saveSpy = sinon.spy(this, 'onSave');

    await render(hbs`{{edit-form model=this.model onSave=this.onSave}}`);

    component.submit();
    later(() => cancelTimers(), 50);
    await settled();

    assert.true(saveSpy.calledOnce, 'calls passed onSave');
    assert.strictEqual(saveSpy.getCall(0).args[0].saveType, 'save');
    assert.deepEqual(saveSpy.getCall(0).args[0].model, this.model, 'passes model to onSave');
    const flash = this.owner.lookup('service:flash-messages');
    assert.strictEqual(flash.success.callCount, 1, 'calls flash message success');
  });
});
