/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { later, _cancelTimers as cancelTimers } from '@ember/runloop';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub, noopStub } from 'vault/tests/helpers/stubs';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import mountBackendForm from '../../pages/components/mount-backend-form';

import sinon from 'sinon';

const component = create(mountBackendForm);

module('Integration | Component | mount backend form', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.flashMessages = this.owner.lookup('service:flash-messages');
    this.flashMessages.registerTypes(['success', 'danger']);
    this.flashSuccessSpy = sinon.spy(this.flashMessages, 'success');
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.server.post('/sys/auth/foo', noopStub());
    this.server.post('/sys/mounts/foo', noopStub());
    this.onMountSuccess = sinon.spy();
  });

  hooks.afterEach(function () {
    this.server.shutdown();
  });

  module('auth method', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.store.createRecord('auth-method');
      this.model.set('config', this.store.createRecord('mount-config'));
    });

    test('it renders default state', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert.strictEqual(
        component.header,
        'Enable an Authentication Method',
        'renders auth header in default state'
      );
      assert.ok(component.types.length > 0, 'renders type picker');
    });

    test('it changes path when type is changed', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.selectType('aws');
      assert.strictEqual(component.pathValue, 'aws', 'sets the value of the type');
      await component.back();
      await component.selectType('approle');
      assert.strictEqual(component.pathValue, 'approle', 'updates the value of the type');
    });

    test('it keeps path value if the user has changed it', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.selectType('approle');
      assert.strictEqual(this.model.type, 'approle', 'Updates type on model');
      assert.strictEqual(component.pathValue, 'approle', 'defaults to approle (first in the list)');
      await component.path('newpath');
      assert.strictEqual(this.model.path, 'newpath', 'Updates path on model');
      await component.back();
      assert.strictEqual(this.model.type, '', 'Clears type on back');
      assert.strictEqual(this.model.path, 'newpath', 'Path is still newPath');
      await component.selectType('aws');
      assert.strictEqual(this.model.type, 'aws', 'Updates type on model');
      assert.strictEqual(component.pathValue, 'newpath', 'keeps custom path value');
    });

    test('it does not show a selected token type when first mounting an auth method', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.selectType('github');
      await component.toggleOptions();
      assert
        .dom('[data-test-input="config.tokenType"]')
        .hasValue('', 'token type does not have a default value.');
      const selectOptions = document.querySelector('[data-test-input="config.tokenType"]').options;
      assert.strictEqual(selectOptions[1].text, 'default-service', 'first option is default-service');
      assert.strictEqual(selectOptions[2].text, 'default-batch', 'second option is default-batch');
      assert.strictEqual(selectOptions[3].text, 'batch', 'third option is batch');
      assert.strictEqual(selectOptions[4].text, 'service', 'fourth option is service');
    });

    test('it calls mount success', async function (assert) {
      assert.expect(3);

      this.server.post('/sys/auth/foo', () => {
        assert.ok(true, 'it calls enable on an auth method');
        return [204, { 'Content-Type': 'application/json' }];
      });
      const spy = sinon.spy();
      this.set('onMountSuccess', spy);

      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.mount('approle', 'foo');
      later(() => cancelTimers(), 50);
      await settled();

      assert.true(spy.calledOnce, 'calls the passed success method');
      assert.true(
        this.flashSuccessSpy.calledWith('Successfully mounted the approle auth method at foo.'),
        'Renders correct flash message'
      );
    });
  });

  module('secrets engine', function (hooks) {
    hooks.beforeEach(function () {
      this.model = this.store.createRecord('secret-engine');
      this.model.set('config', this.store.createRecord('mount-config'));
    });

    test('it renders secret specific headers', async function (assert) {
      await render(
        hbs`<MountBackendForm  @mountType="secret" @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert.strictEqual(component.header, 'Enable a Secrets Engine', 'renders secrets header');
      assert.ok(component.types.length > 0, 'renders type picker');
    });

    test('it changes path when type is changed', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountType="secret" @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.selectType('azure');
      assert.strictEqual(component.pathValue, 'azure', 'sets the value of the type');
      await component.back();
      await component.selectType('nomad');
      assert.strictEqual(component.pathValue, 'nomad', 'updates the value of the type');
    });

    test('it keeps path value if the user has changed it', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountType="secret" @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await component.selectType('kv');
      assert.strictEqual(this.model.type, 'kv', 'Updates type on model');
      assert.strictEqual(component.pathValue, 'kv', 'path matches mount type');
      await component.path('newpath');
      assert.strictEqual(this.model.path, 'newpath', 'Updates path on model');
      await component.back();
      assert.strictEqual(this.model.type, '', 'Clears type on back');
      assert.strictEqual(this.model.path, 'newpath', 'path is still newpath');
      await component.selectType('ssh');
      assert.strictEqual(this.model.type, 'ssh', 'Updates type on model');
      assert.strictEqual(component.pathValue, 'newpath', 'path stays the same');
    });

    test('it calls mount success', async function (assert) {
      assert.expect(3);

      this.server.post('/sys/mounts/foo', () => {
        assert.ok(true, 'it calls enable on an secrets engine');
        return [204, { 'Content-Type': 'application/json' }];
      });
      const spy = sinon.spy();
      this.set('onMountSuccess', spy);

      await render(
        hbs`<MountBackendForm @mountType="secret" @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      await component.mount('ssh', 'foo');
      later(() => cancelTimers(), 50);
      await settled();

      assert.true(spy.calledOnce, 'calls the passed success method');
      assert.true(
        this.flashSuccessSpy.calledWith('Successfully mounted the ssh secrets engine at foo.'),
        'Renders correct flash message'
      );
    });
  });
});
