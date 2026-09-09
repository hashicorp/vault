/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit, waitUntil } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

// To use this helper for configurable engines
// define `this.mountAndConfig` and this.expectedConfigEditRoute in the beforeEach hook
// (see "Acceptance | ldap | overview" as an example)
const BASE_ROUTE = 'vault.cluster.secrets.backend';

export default (test, type) => {
  const {
    isConfigurable = false,
    configRoute = 'configuration.plugin-settings',
    engineRoute = 'list-root',
    isOnlyMountable,
  } = engineDisplayData(type);

  if (isConfigurable) {
    test('(configurable): it navigates from the list view when NOT configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines`);

      await fillIn(GENERAL.inputSearch('secret-engine-path'), backend);
      await click(GENERAL.menuTrigger);
      await click(GENERAL.menuItem('Configure'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${this.expectedConfigEditRoute}`,
        'it navigates to the configure route from the list view'
      );

      await runCmd(deleteEngineCmd(backend));
    });

    test('(configurable): it renders tabs when NOT configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines/${backend}/configuration/general-settings`);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to general settings'
      );
      assert.dom(GENERAL.tabLink('general-settings')).hasClass('active');
      assert.dom(GENERAL.tab('plugin-settings')).exists();
      assert.dom(GENERAL.tabLink('plugin-settings')).doesNotHaveClass('active');

      await click(GENERAL.tabLink('plugin-settings'));
      assert.dom(GENERAL.tabLink('plugin-settings')).hasClass('active', 'plugin-settings is now active');
      assert
        .dom(GENERAL.tabLink('general-settings'))
        .doesNotHaveClass('active', 'general-settings is no longer active');
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${this.expectedConfigEditRoute}`,
        'it redirects to the edit route for the plugin'
      );

      await runCmd(deleteEngineCmd(backend));
    });

    // The dropdown only renders for engines that can be managed using the UI, e.g. "Azure" is only mountable so skip these tests
    if (!isOnlyMountable) {
      test('(configurable): it navigates when NOT configured via dropdown', async function (assert) {
        const backend = `${type}-${uuidv4()}-nav-test`;
        await runCmd(mountEngineCmd(type, backend));

        await visit(this.overviewUrl(backend));
        await click(GENERAL.dropdownToggle('Manage'));
        await click(GENERAL.menuItem('Configure'));

        assert.strictEqual(
          currentRouteName(),
          `${BASE_ROUTE}.${this.expectedConfigEditRoute}`,
          'it redirects to the plugins edit route'
        );
        assert.dom(GENERAL.tab('general-settings')).exists();
        assert.dom(GENERAL.tabLink('plugin-settings')).hasClass('active');

        await runCmd(deleteEngineCmd(backend));
      });

      test('(configurable): it navigates when configured via dropdown', async function (assert) {
        const backend = `${type}-${uuidv4()}-nav-test`;
        await runCmd(mountEngineCmd(type, backend));

        await visit(this.overviewUrl(backend));
        await click(GENERAL.dropdownToggle('Manage'));
        await click(GENERAL.menuItem('Configure'));

        assert.strictEqual(
          currentRouteName(),
          `${BASE_ROUTE}.${this.expectedConfigEditRoute}`,
          'it redirects to the plugins edit route'
        );
        assert.dom(GENERAL.tab('general-settings')).exists();
        assert.dom(GENERAL.tabLink('plugin-settings')).hasClass('active');

        await runCmd(deleteEngineCmd(backend));
      });
    }

    test('(configurable): it navigates from the list view when configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);
      await visit(`/vault/secrets-engines`);
      await fillIn(GENERAL.inputSearch('secret-engine-path'), backend);
      await click(GENERAL.menuTrigger);
      await click(GENERAL.menuItem('Configure'));

      // For configurable engines, clicking "View configuration" will direct to its plugin settings route
      await waitUntil(() => currentRouteName() === `${BASE_ROUTE}.${configRoute}`);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configRoute}`,
        'it navigates to the configure edit route from the list view'
      );

      await runCmd(deleteEngineCmd(backend));
    });

    test('(configurable): it renders tabs when configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);

      await visit(`/vault/secrets-engines/${backend}/configuration/general-settings`);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to the "general-settings" route'
      );
      assert.dom(GENERAL.tabLink('general-settings')).hasClass('active');
      assert.dom(GENERAL.tab('plugin-settings')).exists();
      await click(GENERAL.tabLink('plugin-settings'));
      // Confirm tabs after clicking plugin-settings
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configRoute}`,
        'it navigates to the read route when configured'
      );
      assert
        .dom(GENERAL.tabLink('general-settings'))
        .doesNotHaveClass('active', 'general-settings is no longer active');
      assert.dom(GENERAL.tabLink('plugin-settings')).hasClass('active', 'plugin-settings is now active');

      // Navigate to edit, visually the tabs look the same but this is a different route
      await click(SES.configure);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${this.expectedConfigEditRoute}`,
        'it navigates to the edit route'
      );
      assert
        .dom(GENERAL.tabLink('general-settings'))
        .doesNotHaveClass('active', 'general-settings is still not active');
      assert
        .dom(GENERAL.tabLink('plugin-settings'))
        .hasClass('active', 'plugin-settings is still active after clicking edit');

      await runCmd(deleteEngineCmd(backend));
    });

    test(`(configurable): it navigates to the appropriate page when "Exit configuration" is clicked in the plugin settings route`, async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);

      await visit(`/vault/secrets-engines/${backend}/configuration/general-settings`);
      await click(GENERAL.tabLink('plugin-settings'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configRoute}`,
        'it navigates to the read route when configured'
      );
      await click(GENERAL.button('Exit configuration'));

      // added check for mountable only engines to navigate back to secrets engines list as they don't have an overview page
      const route = isOnlyMountable ? 'vault.cluster.secrets.backends' : `${BASE_ROUTE}.${engineRoute}`;
      assert.strictEqual(
        currentRouteName(),
        route,
        `it navigates to the ${route} route when "Exit configuration" is clicked`
      );

      await runCmd(deleteEngineCmd(backend));
    });
  } else {
    // NON-CONFIGURABLE ENGINES
    test('it should hide plugin settings tab', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines/${backend}/configuration/general-settings`);

      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to the "general-settings" route'
      );
      assert.dom(GENERAL.tabLink('general-settings')).hasClass('active');
      assert.dom(GENERAL.tab('plugin-settings')).doesNotExist();

      await runCmd(deleteEngineCmd(backend));
    });

    test(`it navigates to the ${engineRoute} page when "Exit configuration" is clicked in general-settings`, async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;

      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines/${backend}/configuration/general-settings`);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to the "general-settings" route'
      );

      await click(GENERAL.button('Exit configuration'));

      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${engineRoute}`,
        `it navigates to the ${engineRoute} route when "Exit configuration" is clicked`
      );

      await runCmd(deleteEngineCmd(backend));
    });
  }
};
