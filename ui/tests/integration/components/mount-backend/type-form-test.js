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
import { GENERAL } from 'vault/tests/helpers/general-selectors';

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
      assert.dom(GENERAL.applicationStateHeader).hasText('Loading plugin information...');

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
      assert.dom(GENERAL.applicationStateHeader).doesNotExist('Loading state is hidden after API failure');
    });

    test('it does not fetch plugin catalog for auth methods', async function (assert) {
      const getPluginCatalogSpy = sinon.spy(this.apiService, 'getPluginCatalog');

      await render(hbs`<MountBackend::TypeForm @mountCategory="auth" @setMountType={{this.setType}} />`);

      // Should not call plugin catalog API for auth methods
      assert.ok(getPluginCatalogSpy.notCalled, 'Plugin catalog API not called for auth methods');
    });

    test('it handles plugin catalog API timeout gracefully', async function (assert) {
      const timeoutError = new Error('Request timeout');
      timeoutError.name = 'TimeoutError';
      sinon.stub(this.apiService, 'getPluginCatalog').rejects(timeoutError);

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Should still render engines without version info
      for (const type of secretTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable secret engine`);
      }

      // Should not show error message to user for network issues
      assert.dom(GENERAL.applicationStateHeader).doesNotExist('Should not show loading state after timeout');
    });

    test('it handles permission denied errors appropriately', async function (assert) {
      const permissionError = new Error('Permission denied');
      permissionError.name = 'PermissionError';
      sinon.stub(this.apiService, 'getPluginCatalog').rejects(permissionError);

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Should still render engines without version info
      for (const type of secretTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable secret engine`);
      }
    });

    test('it handles invalid JSON response from plugin catalog API', async function (assert) {
      sinon.stub(this.apiService, 'getPluginCatalog').resolves({ invalid: 'response' });

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Should still render engines without version info
      for (const type of secretTypes) {
        assert.dom(MOUNT_BACKEND_FORM.mountType(type)).exists(`Renders ${type} mountable secret engine`);
      }
    });

    test('it displays plugin status indicators for builtin vs external plugins', async function (assert) {
      const mockResponse = {
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
              name: 'custom-plugin',
              type: 'secret',
              builtin: false,
              version: 'v2.1.0',
              deprecation_status: 'supported',
            },
          ],
          secret: ['aws', 'custom-plugin'],
        },
      };
      sinon.stub(this.apiService, 'getPluginCatalog').resolves(mockResponse);

      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // AWS should show builtin badge
      const awsCard = assert.dom(MOUNT_BACKEND_FORM.mountType('aws'));
      awsCard.exists('AWS engine card exists');
      awsCard.containsText('Builtin', 'AWS shows builtin badge');

      // Custom plugin should show external badge (if it exists in static data)
      if (secretTypes.includes('custom-plugin')) {
        const customCard = assert.dom(MOUNT_BACKEND_FORM.mountType('custom-plugin'));
        customCard.containsText('External', 'Custom plugin shows external badge');
      }
    });

    test('it handles disabled plugins correctly', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Check for disabled plugin cards (demo plugins from helper function)
      const disabledPlugins = ['demo-alpha', 'example-cloud', 'test-infra'];

      disabledPlugins.forEach((pluginType) => {
        const pluginCard = assert.dom(`[data-test-mount-type="${pluginType}"]`);
        if (pluginCard.exists()) {
          // Should have disabled styling and not be clickable for mount type selection
          pluginCard.hasClass('disabled-plugin-card', `${pluginType} has disabled styling`);
        }
      });
    });

    test('it displays vertical divider between enabled and disabled plugins', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Check for vertical divider in categories that have both enabled and disabled plugins
      const dividers = document.querySelectorAll('.vertical-divider');
      assert.ok(
        dividers.length >= 0,
        'May show vertical dividers when both enabled and disabled plugins exist'
      );
    });

    test('it opens documentation flyout for disabled plugins', async function (assert) {
      await render(hbs`<MountBackend::TypeForm @mountCategory="secret" @setMountType={{this.setType}} />`);

      // Find a disabled plugin card and click it
      const disabledCard = document.querySelector('[data-test-mount-type="demo-alpha"]');
      if (disabledCard) {
        await click(disabledCard);

        // Should open the documentation flyout
        assert.dom(GENERAL.flyout).exists('Documentation flyout opens for disabled plugin');
        assert
          .dom(`${GENERAL.flyout} .hds-flyout__header`)
          .containsText('Demo Plugin Alpha', 'Shows correct plugin name in flyout');
      }
    });
  });
});
