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
import { setRunOptions } from 'ember-a11y-testing/test-support';

const secretTypes = mountableEngines().map((engine) => engine.type);
const allSecretTypes = allEngines().map((engine) => engine.type);
const authTypes = methods().map((auth) => auth.type);
const allAuthTypes = allMethods().map((auth) => auth.type);

module('Integration | Component | mount-backend/type-form', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.setType = sinon.spy();
  });

  test('it calls secrets setMountType when type is selected', async function (assert) {
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @mountType="secret" @setMountType={{this.setType}} />`);

    assert
      .dom('[data-test-mount-type]')
      .exists({ count: secretTypes.length }, 'Renders all mountable engines');
    await click(`[data-test-mount-type="ssh"]`);
    assert.ok(spy.calledOnceWith('ssh'));
  });

  test('it calls auth setMountType when type is selected', async function (assert) {
    const spy = sinon.spy();
    this.set('setType', spy);
    await render(hbs`<MountBackend::TypeForm @setMountType={{this.setType}} />`);

    assert
      .dom('[data-test-mount-type]')
      .exists({ count: authTypes.length }, 'Renders all mountable auth methods');
    await click(`[data-test-mount-type="okta"]`);
    assert.ok(spy.calledOnceWith('okta'));
  });

  module('Enterprise', function (hooks) {
    hooks.beforeEach(function () {
      this.version = this.owner.lookup('service:version');
      this.version.type = 'enterprise';
    });

    test('it renders correct items for enterprise secrets', async function (assert) {
      setRunOptions({
        rules: {
          // TODO: Fix disabled enterprise options with enterprise badge
          'color-contrast': { enabled: false },
        },
      });
      await render(hbs`<MountBackend::TypeForm @mountType="secret" @setMountType={{this.setType}} />`);
      assert
        .dom('[data-test-mount-type]')
        .exists({ count: allSecretTypes.length }, 'Renders all secret engines');
    });

    test('it renders correct items for enterprise auth methods', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountType="auth" @setMountType={{this.setType}} />`);
      assert.dom('[data-test-mount-type]').exists({ count: allAuthTypes.length }, 'Renders all auth methods');
    });
  });
});
