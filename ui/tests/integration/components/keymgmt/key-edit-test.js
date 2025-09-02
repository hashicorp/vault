/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import sinon from 'sinon';
import EmberObject from '@ember/object';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import timestamp from 'core/utils/timestamp';

module('Integration | Component | keymgmt/key-edit', function (hooks) {
  setupRenderingTest(hooks);

  hooks.before(function () {
    this.timestampStub = sinon.replace(timestamp, 'now', sinon.fake.returns(new Date('2018-04-03T14:15:30')));
  });
  hooks.beforeEach(function () {
    const now = this.timestampStub();
    const model = EmberObject.create({
      name: 'Unicorns',
      id: 'Unicorns',
      minEnabledVersion: 1,
      versions: [
        {
          id: 1,
          creation_time: now.toString(),
        },
        {
          id: 2,
          creation_time: now.toString(),
        },
      ],
      canDelete: true,
    });
    this.model = model;
    this.tab = '';
  });

  // TODO: Add capabilities tests
  test('it renders show view as default', async function (assert) {
    assert.expect(8);
    await render(hbs`<Keymgmt::KeyEdit @model={{this.model}} @tab={{this.tab}} />`);
    assert.dom('[data-test-secret-header]').hasText('Unicorns', 'Shows key name');
    assert.dom('[data-test-keymgmt-key-toolbar]').exists('Subnav toolbar exists');
    assert.dom('[data-test-tab="Details"]').exists('Details tab exists');
    assert.dom('[data-test-tab="Versions"]').exists('Versions tab exists');
    assert.dom('[data-test-keymgmt-key-destroy]').isDisabled('Destroy button is disabled');
    assert.dom('[data-test-keymgmt-dist-empty-state]').exists('Distribution empty state exists');

    this.set('tab', 'versions');
    assert.dom('[data-test-keymgmt-key-version]').exists({ count: 2 }, 'Renders two version list items');
    assert
      .dom('[data-test-keymgmt-key-current-min]')
      .exists({ count: 1 }, 'Checks only one as current minimum');
  });

  test('it renders the correct elements on edit view', async function (assert) {
    assert.expect(4);
    const model = EmberObject.create({
      name: 'Unicorns',
      id: 'Unicorns',
    });
    this.set('mode', 'edit');
    this.set('model', model);

    await render(hbs`<Keymgmt::KeyEdit @model={{this.model}} @mode={{this.mode}} />`);
    assert.dom('[data-test-secret-header]').hasText('Edit Key', 'Shows edit header');
    assert.dom('[data-test-keymgmt-key-toolbar]').doesNotExist('Subnav toolbar does not exist');
    assert.dom('[data-test-tab="Details"]').doesNotExist('Details tab does not exist');
    assert.dom('[data-test-tab="Versions"]').doesNotExist('Versions tab does not exist');
  });

  test('it renders the correct elements on create view', async function (assert) {
    assert.expect(4);
    const model = EmberObject.create({});
    this.set('mode', 'create');
    this.set('model', model);

    await render(hbs`<Keymgmt::KeyEdit @model={{this.model}} @mode={{this.mode}} />`);
    assert.dom('[data-test-secret-header]').hasText('Create Key', 'Shows edit header');
    assert.dom('[data-test-keymgmt-key-toolbar]').doesNotExist('Subnav toolbar does not exist');
    assert.dom('[data-test-tab="Details"]').doesNotExist('Details tab does not exist');
    assert.dom('[data-test-tab="Versions"]').doesNotExist('Versions tab does not exist');
  });

  test('it defaults to keyType rsa-2048', async function (assert) {
    assert.expect(1);
    const store = this.owner.lookup('service:store');
    this.model = store.createRecord('keymgmt/key');
    this.set('mode', 'create');
    await render(hbs`<Keymgmt::KeyEdit @model={{this.model}} @mode={{this.mode}} />`);
    assert.dom('[data-test-input="type"]').hasValue('rsa-2048', 'Has type rsa-2048 by default');
  });
});
