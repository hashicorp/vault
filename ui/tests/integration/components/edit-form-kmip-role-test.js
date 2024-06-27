/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, run, _cancelTimers as cancelTimers } from '@ember/runloop';
import { resolve } from 'rsvp';
import EmberObject, { computed } from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import { setupEngine } from 'ember-engines/test-support';
import { COMPUTEDS } from 'vault/models/kmip/role';

const flash = Service.extend({
  success: sinon.stub(),
});
const namespace = Service.extend({});

const fieldToCheckbox = (field) => ({ name: field, type: 'boolean' });

const createModel = (options) => {
  const model = EmberObject.extend(COMPUTEDS, {
    /* eslint-disable ember/avoid-leaking-state-in-ember-objects */
    newFields: [
      'role',
      'operationActivate',
      'operationAddAttribute',
      'operationAll',
      'operationCreate',
      'operationDestroy',
      'operationDiscoverVersion',
      'operationGet',
      'operationGetAttributes',
      'operationLocate',
      'operationNone',
      'operationRekey',
      'operationRevoke',
      'tlsClientKeyBits',
      'tlsClientKeyType',
      'tlsClientTtl',
    ],
    fields: computed('operationFields', function () {
      return this.operationFields.map(fieldToCheckbox);
    }),
    destroyRecord() {
      return resolve();
    },
    save() {
      return resolve();
    },
    rollbackAttributes() {},
  });
  return model.create({
    ...options,
  });
};

module('Integration | Component | edit form kmip role', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kmip');

  hooks.beforeEach(function () {
    this.context = { owner: this.engine }; // this.engine set by setupEngine
    run(() => {
      this.engine.unregister('service:flash-messages');
      this.engine.register('service:flash-messages', flash);
      this.engine.register('service:namespace', namespace);
    });
  });

  test('it renders: new model', async function (assert) {
    assert.expect(3);
    const model = createModel({ isNew: true });
    this.set('model', model);
    this.onSave = ({ model }) => {
      assert.false(model.operationNone, 'callback fires with operationNone as false');
      assert.true(model.operationAll, 'callback fires with operationAll as true');
    };
    await render(hbs`<EditFormKmipRole @model={{this.model}} @onSave={{this.onSave}} />`, this.context);

    assert.dom('[data-test-input="operationAll"]').isChecked('sets operationAll');
    await click('[data-test-edit-form-submit]');
  });

  test('it renders: operationAll', async function (assert) {
    assert.expect(3);
    const model = createModel({ operationAll: true });
    this.set('model', model);
    this.onSave = ({ model }) => {
      assert.false(model.operationNone, 'callback fires with operationNone as false');
      assert.true(model.operationAll, 'callback fires with operationAll as true');
    };
    await render(hbs`<EditFormKmipRole @model={{this.model}} @onSave={{this.onSave}} />`, this.context);
    assert.dom('[data-test-input="operationAll"]').isChecked('sets operationAll');
    await click('[data-test-edit-form-submit]');
  });

  test('it renders: operationNone', async function (assert) {
    assert.expect(2);
    const model = createModel({ operationNone: true, operationAll: undefined });
    this.set('model', model);

    this.onSave = ({ model }) => {
      assert.true(model.operationNone, 'callback fires with operationNone as true');
    };
    await render(hbs`<EditFormKmipRole @model={{this.model}} @onSave={{this.onSave}} />`, this.context);
    assert.dom('[data-test-input="operationNone"]').isNotChecked('sets operationNone');
    await click('[data-test-edit-form-submit]');
  });

  test('it renders: choose operations', async function (assert) {
    assert.expect(3);
    const model = createModel({ operationGet: true });
    this.set('model', model);
    this.onSave = ({ model }) => {
      assert.false(model.operationNone, 'callback fires with operationNone as false');
    };
    await render(hbs`<EditFormKmipRole @model={{this.model}} @onSave={{this.onSave}} />`, this.context);

    assert.dom('[data-test-input="operationNone"]').isChecked('sets operationNone');
    assert.dom('[data-test-input="operationAll"]').isNotChecked('sets operationAll');
    await click('[data-test-edit-form-submit]');
  });

  test('it saves operationNone=true when unchecking operationAll box', async function (assert) {
    assert.expect(15);
    const model = createModel({ isNew: true });
    this.set('model', model);
    this.onSave = ({ model }) => {
      assert.true(model.operationNone, 'callback fires with operationNone as true');
      assert.false(model.operationAll, 'callback fires with operationAll as false');
    };

    await render(hbs`<EditFormKmipRole @model={{this.model}} @onSave={{this.onSave}} />`, this.context);
    await click('[data-test-input="operationAll"]');
    for (const field of model.fields) {
      const { name } = field;
      if (name === 'operationNone') continue;
      assert.dom(`[data-test-input="${name}"]`).isNotChecked(`${name} is unchecked`);
    }

    assert.dom('[data-test-input="operationAll"]').isNotChecked('sets operationAll');
    assert
      .dom('[data-test-input="operationNone"]')
      .isChecked('operationNone toggle is true which means allow operations');
    await click('[data-test-edit-form-submit]');
  });

  const savingTests = [
    [
      'setting operationAll',
      { operationNone: true, operationGet: true },
      'operationNone',
      {
        operationAll: true,
        operationNone: false,
        operationGet: true,
      },
      {
        operationGet: null,
        operationNone: false,
      },
      5,
    ],
    [
      'setting operationNone',
      { operationAll: true, operationCreate: true },
      'operationNone',
      {
        operationAll: false,
        operationNone: true,
        operationCreate: true,
      },
      {
        operationNone: true,
        operationCreate: null,
        operationAll: false,
      },
      6,
    ],

    [
      'setting choose, and selecting an additional item',
      { operationAll: true, operationGet: true, operationCreate: true },
      'operationAll,operationDestroy',
      {
        operationAll: false,
        operationCreate: true,
        operationGet: true,
      },
      {
        operationGet: true,
        operationCreate: true,
        operationDestroy: true,
        operationAll: false,
      },
      7,
    ],
  ];
  for (const testCase of savingTests) {
    const [name, initialState, displayClicks, stateBeforeSave, stateAfterSave, assertionCount] = testCase;
    test(name, async function (assert) {
      assert.expect(assertionCount);
      const model = createModel(initialState);
      this.set('model', model);
      const clickTargets = displayClicks.split(',');
      await render(hbs`<EditFormKmipRole @model={{this.model}} />`, this.context);

      for (const clickTarget of clickTargets) {
        await click(`label[for=${clickTarget}]`);
      }
      for (const beforeStateKey of Object.keys(stateBeforeSave)) {
        assert.strictEqual(
          model.get(beforeStateKey),
          stateBeforeSave[beforeStateKey],
          `sets ${beforeStateKey}`
        );
      }

      click('[data-test-edit-form-submit]');

      later(() => cancelTimers(), 50);
      await settled();

      for (const afterStateKey of Object.keys(stateAfterSave)) {
        assert.strictEqual(
          model.get(afterStateKey),
          stateAfterSave[afterStateKey],
          `sets ${afterStateKey} on save`
        );
      }
    });
  }
});
