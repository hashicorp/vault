/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { allowAllCapabilitiesStub, noopStub } from 'vault/tests/helpers/stubs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { mountBackend } from 'vault/tests/helpers/components/mount-backend-form-helpers';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';

import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';
import AuthMethodForm from 'vault/forms/auth/method';

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
    this.onMountSuccess = sinon.spy();
  });

  module('auth method', function (hooks) {
    hooks.beforeEach(function () {
      const defaults = {
        config: { listing_visibility: false },
      };
      this.model = new AuthMethodForm(defaults, { isNew: true });
    });

    test('it renders default state', async function (assert) {
      assert.expect(15);
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      assert
        .dom(GENERAL.title)
        .hasText('Enable an Authentication Method', 'renders auth header in default state');

      for (const method of filterEnginesByMountCategory({
        mountCategory: 'auth',
        isEnterprise: false,
      }).filter((engine) => engine.type !== 'token')) {
        assert
          .dom(GENERAL.cardContainer(method.type))
          .hasText(method.displayName, `renders type:${method.displayName} picker`);
      }
    });

    test('it changes path when type is changed', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );

      await click(GENERAL.cardContainer('aws'));
      assert.dom(GENERAL.inputByAttr('path')).hasValue('aws', 'sets the value of the type');
      await click(GENERAL.backButton);
      await click(GENERAL.cardContainer('approle'));
      assert.dom(GENERAL.inputByAttr('path')).hasValue('approle', 'updates the value of the type');
    });

    test('it keeps path value if the user has changed it', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await click(GENERAL.cardContainer('approle'));
      assert.strictEqual(this.model.type, 'approle', 'Updates type on model');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('approle', 'defaults to approle (first in the list)');
      await fillIn(GENERAL.inputByAttr('path'), 'newpath');
      assert.strictEqual(this.model.path, 'newpath', 'Updates path on model');
      await click(GENERAL.backButton);
      assert.strictEqual(this.model.type, '', 'Clears type on back');
      assert.strictEqual(this.model.path, 'newpath', 'Path is still newPath');
      await click(GENERAL.cardContainer('aws'));
      assert.strictEqual(this.model.type, 'aws', 'Updates type on model');
      assert.dom(GENERAL.inputByAttr('path')).hasValue('newpath', 'keeps custom path value');
    });

    test('it does not show a selected token type when first mounting an auth method', async function (assert) {
      await render(
        hbs`<MountBackendForm @mountModel={{this.model}} @onMountSuccess={{this.onMountSuccess}} />`
      );
      await click(GENERAL.cardContainer('github'));
      await click(GENERAL.button('Method Options'));
      assert
        .dom(GENERAL.inputByAttr('config.token_type'))
        .hasValue('', 'token type does not have a default value.');
      const selectOptions = document.querySelector(GENERAL.inputByAttr('config.token_type')).options;
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
      await mountBackend('approle', 'foo');

      assert.true(spy.calledOnce, 'calls the passed success method');
      assert.true(
        this.flashSuccessSpy.calledWith('Successfully mounted the approle auth method at foo.'),
        'Renders correct flash message'
      );
    });
  });
});
