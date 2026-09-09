/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render, waitFor } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { Response } from 'miragejs';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { UPGRADE_INFO } from 'vault/constants/upgrade-info';

const SELECTORS = {
  cardDescription: `[data-test-card-description]`,
  cardTitle: `[data-test-card-title]`,
};

module('Integration | Component | Upgrade Path Analyzer', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer' },
    ];
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.19.5';
    this.onSetUpgradeInfo = () => {};
    this.server.get('/sys/release-info', () => ({ data: { versions: UPGRADE_INFO } }));
    // Stub returns one version per minor — the last patch of each line.
    // filteredTargetVersions will drop <= 1.19.5, leaving 1.20.1, 1.21.7, 2.0.1.
    this.server.get('/sys/vault-versions', () => ({
      data: { versions: ['1.18.0', '1.19.5', '1.20.1', '1.21.7', '2.0.1'] },
    }));
  });

  test('it renders', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    assert.dom(GENERAL.breadcrumbs).exists('Breadcrumbs are rendered');
    assert.dom(GENERAL.hdsPageHeaderTitle).hasText('Upgrade path analyzer', 'Page title is correct');
    assert.dom(GENERAL.cardContainer('version-selection')).exists('Card container is rendered');
    assert
      .dom(GENERAL.cardContainer('Known issues'))
      .doesNotExist('Known issues card is not rendered during the initial state');
  });

  test('Analyze button is enabled on page load with the latest version pre-selected', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    assert
      .dom(GENERAL.button('Analyze'))
      .isNotDisabled('Analyze button is enabled because the latest version (2.0.1) is pre-selected');
    assert
      .dom(GENERAL.selectByAttr('target version'))
      .hasValue('2.0.1', 'Latest available version is pre-selected in the dropdown');
  });

  test('it shows an "already on latest" banner when current version is the latest available', async function (assert) {
    this.version.version = '2.0.1';
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    assert
      .dom('[data-test-latest-empty-state]')
      .exists('Already-on-latest banner is shown when no newer versions exist');
    assert
      .dom(GENERAL.button('Analyze'))
      .doesNotExist('Analyze button is hidden when already on the latest version');
    assert
      .dom(GENERAL.selectByAttr('target version'))
      .doesNotExist('Target version dropdown is hidden when already on the latest version');
    assert
      .dom(GENERAL.selectByAttr('2.0.1'))
      .doesNotExist('Current version dropdown is hidden when already on the latest version');
  });

  test('it detects the current version', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    assert.dom(GENERAL.selectByAttr('1.19.5')).exists('Current version is detected');
  });

  test('it displays the overview cards with correct title, description, count, and link', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );

    await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
    await click(GENERAL.button('Analyze'));
    await waitFor(GENERAL.cardContainer('Known issues'));

    assert
      .dom(`${GENERAL.cardContainer('Known issues')} ${SELECTORS.cardTitle}`)
      .hasText('Known issues', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('Known issues')} ${SELECTORS.cardDescription}`)
      .hasText(
        'These are all the known issues documented with the version selected.',
        'Card description is correct'
      );
    assert
      .dom(`${GENERAL.cardContainer('Known issues')} ${GENERAL.badge()}`)
      .hasText('6', 'Known issues count is correct');
    assert.dom(GENERAL.linkTo('Known issues')).exists('Known issues view link exists');

    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${SELECTORS.cardTitle}`)
      .hasText('Breaking changes', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${SELECTORS.cardDescription}`)
      .hasText('These are functional changes from one version to the other.', 'Card description is correct');
    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${GENERAL.badge()}`)
      .hasText('4', 'Breaking changes count is correct');
    assert.dom(GENERAL.linkTo('Breaking changes')).exists('Breaking changes view link exists');

    assert
      .dom(`${GENERAL.cardContainer('New behavior')} ${SELECTORS.cardTitle}`)
      .hasText('New behavior', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('New behavior')} ${SELECTORS.cardDescription}`)
      .hasText(
        'New behavior introduced and released in the version selected.',
        'Card description is correct'
      );
    assert
      .dom(`${GENERAL.cardContainer('New behavior')} ${GENERAL.badge()}`)
      .hasText('2', 'New behavior count is correct');
    assert.dom(GENERAL.linkTo('New behavior')).exists('New behavior view link exists');

    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${SELECTORS.cardTitle}`)
      .hasText('Rollback steps', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${SELECTORS.cardDescription}`)
      .hasText('Follow these steps to safely rollback.', 'Card description is correct');
    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${GENERAL.badge()}`)
      .hasText('8', 'Rollback steps count is correct');
    assert.dom(GENERAL.linkTo('Rollback steps')).exists('Rollback steps view link exists');
  });

  test('it displays the Upgrade steps section with download button', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
    await click(GENERAL.button('Analyze'));
    await waitFor(GENERAL.cardContainer('Known issues'));

    assert.dom('[data-test-upgrade-steps-title]').hasText('Upgrade steps', 'Upgrade steps title is correct');
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'Single instance: upgrade the current Vault instance after creating a backup',
        'Upgrade alert title is shown'
      );
    assert.dom(GENERAL.button('Download steps')).exists('Download steps button is rendered');
  });

  test('it surfaces the release info error message to the user', async function (assert) {
    this.server.get(
      '/sys/release-info',
      () => new Response(500, {}, { errors: ['Release info request failed'] })
    );

    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );

    await waitFor(GENERAL.inlineError);

    assert
      .dom(GENERAL.inlineError)
      .hasText('500: Release info request failed', 'The API error message is shown to the user');
  });

  test('it surfaces the vault versions error message to the user', async function (assert) {
    this.server.get(
      '/sys/vault-versions',
      () => new Response(500, {}, { errors: ['Vault versions request failed'] })
    );

    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );

    await waitFor(GENERAL.inlineError);

    assert
      .dom(GENERAL.inlineError)
      .hasText('500: Vault versions request failed', 'The API error message is shown to the user');
  });

  module('release info filtering', function () {
    test('it filters breaking changes and new behavior to only those introduced between currentVersion (exclusive) and targetVersion (inclusive)', async function (assert) {
      await render(
        hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
      );
      await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
      await click(GENERAL.button('Analyze'));
      await waitFor(GENERAL.cardContainer('Breaking changes'));

      assert
        .dom(`${GENERAL.cardContainer('Breaking changes')} ${GENERAL.badge()}`)
        .hasText('4', '4 breaking changes introduced in (1.19.5, 1.20.1]; 1.20.4 and 1.21.0 items excluded');

      assert
        .dom(`${GENERAL.cardContainer('New behavior')} ${GENERAL.badge()}`)
        .hasText('2', '2 new behavior entries introduced in (1.19.5, 1.20.1]; later items excluded');
    });

    test('it filters known issues to those found in range and not fixed by targetVersion', async function (assert) {
      await render(
        hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
      );
      await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
      await click(GENERAL.button('Analyze'));
      await waitFor(GENERAL.cardContainer('Known issues'));

      assert
        .dom(`${GENERAL.cardContainer('Known issues')} ${GENERAL.badge()}`)
        .hasText('6', '6 known issues found in range and still unresolved at 1.20.1');
    });

    test('it excludes known issues already fixed by the targetVersion', async function (assert) {
      await render(
        hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
      );
      await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
      await click(GENERAL.button('Analyze'));
      await waitFor(GENERAL.cardContainer('Known issues'));

      assert
        .dom(`${GENERAL.cardContainer('Known issues')} ${GENERAL.badge()}`)
        .hasText('6', '6 known issues still unresolved at 1.20.1 (4 fixed by 1.20.1 are excluded)');
    });

    test('it excludes items found before the currentVersion', async function (assert) {
      this.version.version = '1.19.5';

      await render(
        hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
      );
      await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
      await click(GENERAL.button('Analyze'));
      await waitFor(GENERAL.cardContainer('Known issues'));

      assert
        .dom(`${GENERAL.cardContainer('Known issues')} ${GENERAL.badge()}`)
        .hasText('6', 'Items with found version <= currentVersion (1.18.4, 1.19.0) are excluded');
    });

    test('it passes the correctly filtered info to the onSetUpgradeInfo callback', async function (assert) {
      assert.expect(2);
      this.onSetUpgradeInfo = (info) => {
        assert.strictEqual(
          info.breaking_changes.length,
          4,
          '4 breaking changes introduced in (1.19.5, 1.20.1]'
        );
        assert.strictEqual(
          info.known_issues.length,
          6,
          '6 known issues found in range and unresolved at 1.20.1'
        );
      };

      await render(
        hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
      );
      await fillIn(GENERAL.selectByAttr('target version'), '1.20.1');
      await click(GENERAL.button('Analyze'));
      await waitFor(GENERAL.cardContainer('Known issues'));
    });
  });
});
