/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { allEngines, mountableEngines } from 'vault/helpers/mountable-secret-engines';
import { allMethods, methods } from 'vault/helpers/mountable-auth-methods';

const secretTypes = mountableEngines().map((engine) => engine.type);
const allSecretTypes = allEngines().map((engine) => engine.type);
const authTypes = methods().map((auth) => auth.type);
const allAuthTypes = allMethods().map((auth) => auth.type);

module('Integration | Component | mount-backend/type-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.setType = sinon.spy();
  });

  test('it calls secrets setMountType only on next click', async function (assert) {
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @mountType="secret" @setMountType={{this.setType}} />`);

    assert
      .dom('[data-test-mount-type]')
      .exists({ count: secretTypes.length }, 'Renders all mountable engines');
    await click(`[data-test-mount-type="nomad"]`);
    assert.dom(`[data-test-mount-type="nomad"] input`).isChecked(`ssh is checked`);
    assert.ok(spy.notCalled, 'callback not called');
    await click(`[data-test-mount-type="ssh"]`);
    assert.dom(`[data-test-mount-type="ssh"] input`).isChecked(`ssh is checked`);
    assert.ok(spy.notCalled, 'callback not called');
    await click('[data-test-mount-next]');
    assert.ok(spy.calledOnceWith('ssh'));
  });

  test('it calls auth setMountType only on next click', async function (assert) {
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @setMountType={{this.setType}} />`);

    assert
      .dom('[data-test-mount-type]')
      .exists({ count: authTypes.length }, 'Renders all mountable auth methods');
    await click(`[data-test-mount-type="okta"]`);
    assert.dom(`[data-test-mount-type="okta"] input`).isChecked(`ssh is checked`);
    assert.ok(spy.notCalled, 'callback not called');
    await click(`[data-test-mount-type="github"]`);
    assert.dom(`[data-test-mount-type="github"] input`).isChecked(`ssh is checked`);
    assert.ok(spy.notCalled, 'callback not called');
    await click('[data-test-mount-next]');
    assert.ok(spy.calledOnceWith('github'));
  });

  module('Enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.version = '1.12.1+ent';
    });

    test('it renders correct items for enterprise secrets', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountType="secret" @setMountType={{this.setType}} />`);
      assert
        .dom('[data-test-mount-type]')
        .exists({ count: allSecretTypes.length }, 'Renders all secret engines');
    });

    test('it renders correct items for enterprise auth methods', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountType="secret" @setMountType={{this.setType}} />`);
      assert
        .dom('[data-test-mount-type]')
        .exists({ count: allAuthTypes.length }, 'Renders all secret engines');
    });
  });
});
