/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import { setRunOptions } from 'ember-a11y-testing/test-support';
import { MOUNT_BACKEND_FORM } from 'vault/tests/helpers/components/mount-backend-form-selectors';

const secretTypes = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: false })
  .filter((engine) => engine.type !== 'cubbyhole')
  .map((engine) => engine.type);
const allSecretTypes = filterEnginesByMountCategory({ mountCategory: 'secret', isEnterprise: true })
  .filter((engine) => engine.type !== 'cubbyhole')
  .map((engine) => engine.type);
const authTypes = filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: false })
  .filter((engine) => engine.type !== 'token')
  .map((auth) => auth.type);
const allAuthTypes = filterEnginesByMountCategory({ mountCategory: 'auth', isEnterprise: true })
  .filter((engine) => engine.type !== 'token')
  .map((auth) => auth.type);

module('Integration | Component | mount-backend/type-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.setType = sinon.spy();
  });

  test('it calls secrets setMountType when type is selected', async function (assert) {
    assert.expect(secretTypes.length + 1, 'renders all mountable engines plus calls a spy');
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

    for (const type of secretTypes) {
      assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable secret engine`);
    }
    await click(MOUNT_BACKEND_FORM.mountType('ssh'));
    assert.ok(spy.calledOnceWith('ssh'));
  });

  test('it calls auth setMountType when type is selected', async function (assert) {
    assert.expect(authTypes.length + 1, 'renders all mountable auth methods plus calls a spy');
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @setMountType={{this.setType}} />`);

    for (const type of authTypes) {
      assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable auth engine`);
    }
    await click(MOUNT_BACKEND_FORM.mountType('okta'));
    assert.ok(spy.calledOnceWith('okta'));
  });

  module('Enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';
    });

    test('it renders correct items for enterprise secrets', async function (assert) {
      assert.expect(allSecretTypes.length, 'renders all enterprise secret engines');
      setRunOptions({
        rules: {
          // TODO: Fix disabled enterprise options with enterprise badge
          'color-contrast': { enabled: false },
        },
      });
      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);
      for (const type of allSecretTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} secret engine`);
      }
    });

    test('it renders correct items for enterprise auth methods', async function (assert) {
      assert.expect(allAuthTypes.length, 'renders all enterprise auth engines');
      await render(hbs`<MountBackend::TypeForm @mountCategory="auth" @setMountType={{this.setType}} />`);
      for (const type of allAuthTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} auth engine`);
      }
    });
  });
});
