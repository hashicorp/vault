/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { fillIn, render, waitFor } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { keyMgmtMockModel } from 'vault/tests/helpers/secret-engine/mocks';

module('Integration | Component | SecretEngine::Card::Version', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.model = {
      ...keyMgmtMockModel,
      pinnedVersion: null,
    };
  });

  test('it shows version card information', async function (assert) {
    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    assert.dom(`${GENERAL.cardContainer('version')} h2`).hasText('Version');
    assert.dom(GENERAL.infoRowValue('type')).hasAnyText(keyMgmtMockModel.secretsEngine.type);
    assert.dom(GENERAL.inputByAttr('plugin-version')).doesNotExist();
  });

  test('it shows info message when pinned version exists', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.15',
        running_plugin_version: 'v0.17',
      },
      versions: ['v0.16', 'v0.17', 'v0.18'],
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show info message about pinned version override
    assert.dom(GENERAL.inlineAlert).exists('Info message shown when pinned version exists');
    assert
      .dom(`${GENERAL.inlineAlert} .hds-alert__description`)
      .hasText('Pinned plugin version (v0.16) is overridden by the running version (v0.17).');
  });

  test('it shows "(Pinned)" label for pinned version in dropdown', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        running_plugin_version: 'v0.17',
      },
      versions: ['v0.16', 'v0.17', 'v0.18'],
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Check that dropdown exists when versions are available
    assert.dom(GENERAL.inputByAttr('plugin-version')).exists();

    // Running version (v0.17) should be filtered out, so we expect 3 total options:
    // placeholder + v0.16 (Pinned) + v0.18
    assert
      .dom(`select${GENERAL.inputByAttr('plugin-version')} option`)
      .exists({ count: 3 }, 'Has 3 options: placeholder + 2 versions (running version filtered out)');

    // Find the option with the pinned version and verify it shows "(Pinned)"
    const options = this.element.querySelectorAll(`select${GENERAL.inputByAttr('plugin-version')} option`);
    let foundPinnedOption = false;
    let foundNonPinnedOption = false;
    let foundRunningVersionOption = false;

    options.forEach((option) => {
      if (option.value === 'v0.16 (Pinned)') {
        foundPinnedOption = true;
      } else if (option.value === 'v0.18') {
        foundNonPinnedOption = true;
      } else if (option.value === 'v0.17') {
        foundRunningVersionOption = true;
      }
    });

    assert.true(foundPinnedOption, 'Found pinned version option');
    assert.true(foundNonPinnedOption, 'Found non-pinned version option');
    assert.false(foundRunningVersionOption, 'Running version is filtered out of dropdown');

    // Verify the text content of the specific options
    const pinnedOption = this.element.querySelector(
      `select${GENERAL.inputByAttr('plugin-version')} option[value="v0.16 (Pinned)"]`
    );
    const nonPinnedOption = this.element.querySelector(
      `select${GENERAL.inputByAttr('plugin-version')} option[value="v0.18"]`
    );

    assert.strictEqual(
      pinnedOption.textContent.trim(),
      'v0.16 (Pinned)',
      'Pinned version shows "(Pinned)" label'
    );
    assert.strictEqual(
      nonPinnedOption.textContent.trim(),
      'v0.18',
      'Non-pinned version does not show "(Pinned)" label'
    );
  });

  test('it shows override pinned version alert when selecting different version from pinned', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        running_plugin_version: 'v0.17',
      },
      versions: ['v0.16', 'v0.17', 'v0.18'],
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Initially no override alert should be shown
    assert
      .dom(GENERAL.inlineAlertByAttr('override-pinned-version'))
      .doesNotExist('Override alert not shown initially');

    // Select a version different from pinned version
    await fillIn(GENERAL.inputByAttr('plugin-version'), 'v0.18');

    // Alert should now be visible
    assert
      .dom(GENERAL.inlineAlertByAttr('override-pinned-version'))
      .exists('Override alert shown when selecting version different from pinned');

    // Verify the alert text is correct
    assert
      .dom(`${GENERAL.inlineAlertByAttr('override-pinned-version')} .hds-alert__title`)
      .hasText('Override pinned version');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('override-pinned-version')} .hds-alert__description`)
      .hasText(
        'You have selected v0.18, but version v0.16 is pinned for this plugin. Updating to this version will override the pinned version for this mount.'
      );
  });

  test('it does not show info message when no pinned version exists', async function (assert) {
    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should not show info message when no pinned version exists
    assert.dom(GENERAL.inlineAlert).doesNotExist('No info message when no pinned version');
  });

  test('it does not show version mismatch when running equals pinned', async function (assert) {
    const pinnedVersion = 'v0.16';

    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.14',
        running_plugin_version: pinnedVersion,
      },
      pinnedVersion: pinnedVersion,
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should not show version mismatch alert when running === pinned
    assert
      .dom(GENERAL.inlineAlertByAttr('alert-message'))
      .doesNotExist('No version mismatch shown when running version equals pinned version');
  });

  test('it shows version mismatch when current != running and running != pinned', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.15',
        running_plugin_version: 'v0.16',
      },
      pinnedVersion: 'v0.14',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show version mismatch alert when current != running AND running != pinned
    assert
      .dom(GENERAL.inlineAlertByAttr('alert-message'))
      .exists('Version mismatch shown when current != running and running != pinned');
  });

  test('it shows no version mismatch when running equals plugin but differs from pinned', async function (assert) {
    const runningAndPluginVersion = 'v0.15';

    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: runningAndPluginVersion,
        running_plugin_version: runningAndPluginVersion,
      },
      pinnedVersion: 'v0.14',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should not show version mismatch alert when running === plugin even if pinned differs
    assert
      .dom(GENERAL.inlineAlertByAttr('alert-message'))
      .doesNotExist('No version mismatch shown when running equals plugin version');
  });

  test('it shows no override alert when selecting pinned version', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        running_plugin_version: 'v0.17',
      },
      versions: ['v0.16', 'v0.17', 'v0.18'],
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Select the pinned version (should have "(Pinned)" label)
    await fillIn(GENERAL.inputByAttr('plugin-version'), 'v0.16 (Pinned)');

    // No override alert should be shown when selecting the pinned version
    assert
      .dom(GENERAL.inlineAlertByAttr('override-pinned-version'))
      .doesNotExist('No override alert shown when selecting the pinned version');
  });

  test('it shows manual override info message when running != pinned but running == plugin', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.17',
        running_plugin_version: 'v0.17',
      },
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show manual override info message
    assert.dom(GENERAL.inlineAlertByAttr('info-message')).exists('Manual override info message shown');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('info-message')} .hds-alert__description`)
      .hasText('This engine has a manual override of the pinned version (v0.16) by version v0.17.');
  });

  test('it shows configured override info message when running == pinned', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.15',
        running_plugin_version: 'v0.16',
      },
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show configured plugin override info message
    assert.dom(GENERAL.inlineAlertByAttr('info-message')).exists('Configured override info message shown');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('info-message')} .hds-alert__description`)
      .hasText('Configured plugin version (v0.15) is overridden by the pinned version (v0.16).');
  });

  test('it does not show override message when plugin and pinned versions are the same', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.16',
        running_plugin_version: 'v0.16',
      },
      pinnedVersion: 'v0.16',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should not show override message when versions are the same
    assert
      .dom(GENERAL.inlineAlert)
      .doesNotExist('No override message when plugin and pinned versions are the same');
  });

  test('it shows version mismatch info and alert when plugin != running != pinned', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.15',
        running_plugin_version: 'v0.16',
      },
      pinnedVersion: 'v0.14',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show version mismatch info message
    assert.dom(GENERAL.inlineAlertByAttr('info-message')).exists('Version mismatch info message shown');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('info-message')} .hds-alert__description`)
      .hasText('Pinned plugin version (v0.14) is overridden by the running version (v0.16).');

    // Should also show version mismatch alert with reload button
    assert.dom(GENERAL.inlineAlertByAttr('alert-message')).exists('Version mismatch alert shown');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('alert-message')} .hds-alert__title`)
      .hasText('Version mismatch detected');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('alert-message')} .hds-alert__description`)
      .hasText(
        'This plugin is configured to use version v0.15 but is currently running version v0.16. Reload the plugin to sync the running version with the configured version.'
      );
    assert.dom(GENERAL.button('reload-plugin')).exists('Reload button shown in mismatch alert');
  });

  test('it shows only running version row and hides pinned/configured versions', async function (assert) {
    this.server.get('/sys/plugins/pins/secret/:name', () => ({
      data: { version: 'v0.16' },
    }));

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should only show running version since pinned/configured version rows were removed from template
    assert.dom(GENERAL.infoRowValue('running-plugin-version')).exists('Running version row shown');
  });

  test('it shows "(Pinned)" suffix when running version equals pinned version', async function (assert) {
    const pinnedAndRunningVersion = 'v0.16';

    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        running_plugin_version: pinnedAndRunningVersion,
      },
      pinnedVersion: pinnedAndRunningVersion,
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Should show "(Pinned)" suffix in running version
    assert.dom(GENERAL.infoRowValue('running-plugin-version')).containsText('(Pinned)');
  });

  test('it filters out empty string versions but includes running version in dropdown', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.17',
        running_plugin_version: 'v0.17',
      },
      pinnedVersion: 'v0.17',
      versions: ['', 'v0.16', '', 'v0.17', '', 'v0.18', ''],
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Wait for the dropdown options to be rendered
    await waitFor(`select${GENERAL.inputByAttr('plugin-version')} option`);

    // Check that empty versions are filtered out
    const options = this.element.querySelectorAll(`select${GENERAL.inputByAttr('plugin-version')} option`);
    const optionValues = Array.from(options).map((opt) => opt.value);

    // Filter out the placeholder option when checking for empty strings from data
    const emptyDataOptions = Array.from(options).filter(
      (opt) => opt.value === '' && opt.textContent.trim() !== 'Select version'
    );
    assert.strictEqual(
      emptyDataOptions.length,
      0,
      'Empty string versions from data are filtered out of dropdown options'
    );

    // Should have: Select version placeholder + 2 non-empty, non-running versions (v0.16, v0.18)
    // v0.17 is running version and should be filtered out
    assert.strictEqual(optionValues.length, 3, 'Should have 3 total options: placeholder + 2 versions');

    // Count non-placeholder options
    const nonPlaceholderOptions = Array.from(options).filter(
      (opt) => opt.textContent.trim() !== 'Select version'
    );
    assert.strictEqual(
      nonPlaceholderOptions.length,
      2,
      'Should have exactly 2 non-placeholder version options (running version filtered out)'
    );

    // Verify the actual version options present (running version should be filtered out)
    assert.true(optionValues.includes('v0.16'), 'v0.16 is included');
    assert.false(optionValues.includes('v0.17'), 'v0.17 is filtered out (running version)');
    assert.true(optionValues.includes('v0.18'), 'v0.18 is included');
  });

  test('it uses correct alert configurations from computed properties', async function (assert) {
    this.model = {
      ...keyMgmtMockModel,
      secretsEngine: {
        ...keyMgmtMockModel.secretsEngine,
        plugin_version: 'v0.15',
        running_plugin_version: 'v0.16',
      },
      pinnedVersion: 'v0.14',
    };

    await render(hbs`
      <SecretEngine::Card::Version @model={{this.model}} />
    `);

    // Verify version mismatch alert uses versionMismatchAlert computed property
    assert
      .dom(`${GENERAL.inlineAlertByAttr('alert-message')} .hds-alert__title`)
      .hasText('Version mismatch detected');
    assert
      .dom(`${GENERAL.inlineAlertByAttr('alert-message')} .hds-alert__description`)
      .hasText(
        'This plugin is configured to use version v0.15 but is currently running version v0.16. Reload the plugin to sync the running version with the configured version.'
      );

    // Verify reload button is shown from showReloadButton property
    assert.dom(GENERAL.button('reload-plugin')).exists('Reload button shown via computed property');
  });
});
