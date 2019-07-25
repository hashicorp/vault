import { later, run } from '@ember/runloop';
import { resolve } from 'rsvp';
import EmberObject, { computed } from '@ember/object';
import Service from '@ember/service';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { click, find, render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import engineResolverFor from 'ember-engines/test-support/engine-resolver-for';
import { COMPUTEDS } from 'vault/models/kmip/role';
const resolver = engineResolverFor('kmip');

const flash = Service.extend({
  success: sinon.stub(),
});
const namespace = Service.extend({});

const fieldToCheckbox = field => ({ name: field, type: 'boolean' });

const createModel = options => {
  let model = EmberObject.extend(COMPUTEDS, {
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
    fields: computed('operationFields', function() {
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

module('Integration | Component | edit form kmip role', function(hooks) {
  setupRenderingTest(hooks, { resolver });

  hooks.beforeEach(function() {
    run(() => {
      this.owner.unregister('service:flash-messages');
      this.owner.register('service:flash-messages', flash);
      this.owner.register('service:namespace', namespace);
    });
  });

  test('it renders: new model', async function(assert) {
    let model = createModel({ isNew: true });
    this.set('model', model);
    await render(hbs`<EditFormKmipRole @model={{model}} />`);

    assert.dom('[data-test-input="operationAll"').isChecked('sets operationAll');
  });

  test('it renders: operationAll', async function(assert) {
    let model = createModel({ operationAll: true });
    this.set('model', model);
    await render(hbs`<EditFormKmipRole @model={{model}} />`);
    assert.dom('[data-test-input="operationAll"').isChecked('sets operationAll');
  });

  test('it renders: operationNone', async function(assert) {
    let model = createModel({ operationNone: true });
    this.set('model', model);
    await render(hbs`<EditFormKmipRole @model={{model}} />`);

    assert.dom('[data-test-input="operationNone"]').isNotChecked('sets operationNone');
  });

  test('it renders: choose operations', async function(assert) {
    let model = createModel({ operationGet: true });
    this.set('model', model);
    await render(hbs`<EditFormKmipRole @model={{model}} />`);

    assert.dom('[data-test-input="operationNone"]').isChecked('sets operationNone');
    assert.dom('[data-test-input="operationAll"').isNotChecked('sets operationAll');
  });

  let savingTests = [
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
    ],
  ];
  for (let testCase of savingTests) {
    let [name, initialState, displayClicks, stateBeforeSave, stateAfterSave] = testCase;
    test(name, async function(assert) {
      let model = createModel(initialState);
      this.set('model', model);
      let clickTargets = displayClicks.split(',');
      await render(hbs`<EditFormKmipRole @model={{model}} />`);

      for (let clickTarget of clickTargets) {
        await click(`label[for=${clickTarget}]`);
      }
      for (let beforeStateKey of Object.keys(stateBeforeSave)) {
        assert.equal(model.get(beforeStateKey), stateBeforeSave[beforeStateKey], `sets ${beforeStateKey}`);
      }

      click('[data-test-edit-form-submit]');

      later(() => run.cancelTimers(), 50);
      return settled().then(() => {
        for (let afterStateKey of Object.keys(stateAfterSave)) {
          assert.equal(
            model.get(afterStateKey),
            stateAfterSave[afterStateKey],
            `sets ${afterStateKey} on save`
          );
        }
      });
    });
  }
});
