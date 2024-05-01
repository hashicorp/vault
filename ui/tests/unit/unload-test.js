/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { run } from '@ember/runloop';
import Model, { attr } from '@ember-data/model';
import { setupTest } from 'ember-qunit';
import { module, test } from 'qunit';
import Adapter from '@ember/test/adapter';

/**
 * This test is testing ember internals for what we need available on lazyPaginatedQuery
 */
module('Unit | Model | unloadAll works as expected', function (hooks) {
  setupTest(hooks);

  hooks.beforeEach(function () {
    const Company = Model.extend({
      name: attr('string'),
    });

    const CompanyAdapter = Adapter.extend({
      updateRecord: () => {
        return undefined;
      },
      createRecord: () => {
        return {
          id: '4',
          data: {
            name: 'Foobar',
          },
        };
      },
    });

    this.owner.register('model:company', Company);
    this.owner.register('adapter:company', CompanyAdapter);
    this.store = this.owner.lookup('service:store');
  });

  test('edit then unload correctly removes all records', async function (assert) {
    this.store.push({
      data: [
        {
          id: '1',
          type: 'company',
          attributes: {
            name: 'ACME',
          },
        },
        {
          id: '2',
          type: 'company',
          attributes: {
            name: 'EMCA',
          },
        },
      ],
    });
    assert.strictEqual(this.store.peekAll('company').length, 2, '2 companies loaded');
    const editRecord = this.store.peekRecord('company', '1');
    editRecord.name = 'Rebrand';
    await editRecord.save();
    assert.false(editRecord.hasDirtyAttributes, 'edit record does not have dirty attrs after save');
    this.store.peekAll('company').length;
    run(() => {
      this.store.unloadAll('company');
    });

    assert.strictEqual(this.store.peekAll('company').length, 0, 'peekAll 0 - companies unloaded');
    assert.strictEqual(
      this.store.peekAll('company').slice().length,
      0,
      'peekAll array 0 - companies unloaded'
    );
  });
  test('create then unload correctly removes all records', async function (assert) {
    this.store.push({
      data: [
        {
          id: '1',
          type: 'company',
          attributes: {
            name: 'ACME',
          },
        },
        {
          id: '2',
          type: 'company',
          attributes: {
            name: 'EMCA',
          },
        },
      ],
    });
    assert.strictEqual(this.store.peekAll('company').length, 2, '2 companies loaded');
    const newRecord = this.store.createRecord('company', { name: 'Foobar' });

    await newRecord.save();
    assert.false(newRecord.hasDirtyAttributes, 'new record does not have dirty attrs after save');
    this.store.peekAll('company').length;
    run(() => {
      this.store.unloadAll('company');
    });

    assert.strictEqual(this.store.peekAll('company').length, 0, 'peekAll 0 - companies unloaded');
    assert.strictEqual(
      this.store.peekAll('company').slice().length,
      0,
      'peekAll array 0 - companies unloaded'
    );
  });
});
