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

  module('Plugin Catalog Integration', function (hooks) {
    hooks.beforeEach(function () {
      this.apiService = this.owner.lookup('service:api');
      this.mockPluginCatalogResponse = {
        data: {
          detailed: [
            {
              name: 'aws',
              type: 'secret',
              builtin: true,
              version: 'v1.12.0+builtin.vault',
              deprecation_status: 'supported',
            },
            {
              name: 'kv',
              type: 'secret',
              builtin: true,
              version: 'v0.13.0+builtin',
              deprecation_status: 'supported',
            },
          ],
          secret: ['aws', 'kv'],
        },
      };
    });

    test('it displays loading state while fetching plugin catalog', async function (assert) {
      // Mock a slow API response
      const slowPromise = new Promise((resolve) => {
        setTimeout(() => resolve(this.mockPluginCatalogResponse), 100);
      });
      sinon.stub(this.apiService, 'getPluginCatalog').returns(slowPromise);

      render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Check for loading state
      assert.dom('[data-test-application-state-header]').hasText('Loading plugin information...');

      // Wait for the API call to complete
      await slowPromise;
    });

    test('it displays version information when plugin catalog is loaded', async function (assert) {
      sinon.stub(this.apiService, 'getPluginCatalog').resolves(this.mockPluginCatalogResponse);

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Check that version information is displayed for engines with plugin data
      const awsCard = assert.dom(MOUNT_BACKEND_FORM.mountType('aws'));
      awsCard.exists('AWS engine card exists');
      awsCard.includesText('v1.12.0+builtin.vault', 'AWS version is displayed');

      const kvCard = assert.dom(MOUNT_BACKEND_FORM.mountType('kv'));
      kvCard.exists('KV engine card exists');
      kvCard.includesText('v0.13.0+builtin', 'KV version is displayed');
    });

    test('it falls back gracefully when plugin catalog API fails', async function (assert) {
      sinon.stub(this.apiService, 'getPluginCatalog').rejects(new Error('API Error'));

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Should still render engines without version info
      for (const type of secretTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable secret engine`);
      }

      // Loading state should not be visible after failure
      assert
        .dom('[data-test-application-state-header]')
        .doesNotExist('Loading state is hidden after API failure');
    });

    test('it does not fetch plugin catalog for auth methods', async function (assert) {
      const getPluginCatalogSpy = sinon.spy(this.apiService, 'getPluginCatalog');

      await render(hbs`<MountBackend::TypeForm @mountCategory="auth" @setMountType={{this.setType}} />`);

      // Should not call plugin catalog API for auth methods
      assert.ok(getPluginCatalogSpy.notCalled, 'Plugin catalog API not called for auth methods');
    });
  });
});
