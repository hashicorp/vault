/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentRouteName, fillIn, visit } from '@ember/test-helpers';
import { v4 as uuidv4 } from 'uuid';
import engineDisplayData from 'vault/helpers/engines-display-data';
import { deleteEngineCmd, mountEngineCmd, runCmd } from 'vault/tests/helpers/commands';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { SECRET_ENGINE_SELECTORS as SES } from 'vault/tests/helpers/secret-engine/secret-engine-selectors';

// To use this helper for configurable engines
// define `this.mountAndConfig` in the beforeEach hook
// (see "Acceptance | ldap | overview" as an example)
const BASE_ROUTE = 'vault.cluster.secrets.backend';

export default (test, type) => {
  const {
    isConfigurable = false,
    configReadRoute = 'configuration.plugin-settings',
    configEditRoute = 'configuration.edit',
    engineRoute = 'list-root',
  } = engineDisplayData(type);

  if (isConfigurable) {
    test('(configurable): it navigates from the list view when NOT configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines`);

      await fillIn(GENERAL.inputSearch('secret-engine-path'), backend);
      await click(GENERAL.menuTrigger);
      await click(GENERAL.menuItem('View configuration'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configEditRoute}`,
        'it navigates to the configure route from the list view'
      );

      await runCmd(deleteEngineCmd(backend));
    });

    test('(configurable): it renders tabs when NOT configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines/${backend}/configuration`);
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to the "general-settings" route'
      );
      assert.dom(GENERAL.tabLink('general-settings')).hasClass('active');
      assert.dom(GENERAL.tab('plugin-settings')).exists();
      await click(GENERAL.tabLink('plugin-settings'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configEditRoute}`,
        'clicking plugin settings navigates to edit route when not configured'
      );
      assert
        .dom(GENERAL.tabLink('general-settings'))
        .doesNotHaveClass('active', 'general-settings is no longer active');
      assert.dom(GENERAL.tabLink('plugin-settings')).hasClass('active', 'plugin-settings is now active');

      await runCmd(deleteEngineCmd(backend));
    });

    test('(configurable): it navigates from the list view when configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);

      await visit(`/vault/secrets-engines`);

      await fillIn(GENERAL.inputSearch('secret-engine-path'), backend);
      await click(GENERAL.menuTrigger);
      await click(GENERAL.menuItem('View configuration'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configReadRoute}`,
        'it navigates to the configure route from the list view'
      );

      await runCmd(deleteEngineCmd(backend));
    });

    test('(configurable): it renders tabs when configured', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);

      await visit(`/vault/secrets-engines/${backend}/configuration`);
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
        `${BASE_ROUTE}.${configReadRoute}`,
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
        `${BASE_ROUTE}.${configEditRoute}`,
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

    test(`(configurable): it navigates to the ${engineRoute} page when "Exit configuration" is clicked in the plugin settings route`, async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await this.mountAndConfig(backend);

      await visit(`/vault/secrets-engines/${backend}/configuration`);
      await click(GENERAL.tabLink('plugin-settings'));
      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${configReadRoute}`,
        'it navigates to the read route when configured'
      );
      await click(GENERAL.button('Exit configuration'));

      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.${engineRoute}`,
        `it navigates to the ${engineRoute} route when "Exit configuration" is clicked`
      );

      await runCmd(deleteEngineCmd(backend));
    });
  } else {
    // NON-CONFIGURABLE ENGINES
    test('it should hide plugin settings tab', async function (assert) {
      const backend = `${type}-${uuidv4()}-nav-test`;
      await runCmd(mountEngineCmd(type, backend));

      await visit(`/vault/secrets-engines/${backend}/configuration`);

      assert.strictEqual(
        currentRouteName(),
        `${BASE_ROUTE}.configuration.general-settings`,
        'it navigates to the "general-settings" route'
      );
      assert.dom(GENERAL.tabLink('general-settings')).hasClass('active');
      assert.dom(GENERAL.tab('plugin-settings')).doesNotExist();

      await runCmd(deleteEngineCmd(backend));
    });
  }

  test(`it navigates to the ${engineRoute} page when "Exit configuration" is clicked in general-settings`, async function (assert) {
    const backend = `${type}-${uuidv4()}-nav-test`;

    await runCmd(mountEngineCmd(type, backend));

    await visit(`/vault/secrets-engines/${backend}/configuration`);
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
};
