/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { UPGRADE_INFO } from 'vault/constants/upgrade-info';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

// Flatten the raw UPGRADE_INFO array into the single DisplayedInfo shape that
// the UpgradeInfo component receives at runtime (produced by filterReleaseInfo).
const flatUpgradeInfo = {
  known_issues: UPGRADE_INFO.flatMap((e) => e.known_issues ?? []),
  breaking_changes: UPGRADE_INFO.flatMap((e) => e.breaking_changes ?? []),
  new_behavior: UPGRADE_INFO.flatMap((e) => e.new_behavior ?? []),
  rollback_steps: [],
  rollbackOrder: [],
  rollbackGuidanceMessage: '',
  targetVersion: '1.21.0',
};

module('Integration | Component | UpgradePathAnalyzer::UpgradeInfo', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer', route: 'vault.cluster.support.upgrade' },
      { label: 'Issues' },
    ];

    this.upgradeInfo = flatUpgradeInfo;
  });

  test('it renders the component with tabs and data', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradeInfo @breadcrumbs={{this.breadcrumbs}} @upgradeInfo={{this.upgradeInfo}}/>`
    );

    // Check that each tab is rendered
    assert.dom(GENERAL.tab('Known issues')).exists();
    assert.dom(GENERAL.tab('Breaking changes')).exists();
    assert.dom(GENERAL.tab('New behavior')).exists();
    assert.dom(GENERAL.tab('Rollback steps')).exists();

    // UPGRADE_INFO has 4 (1.21) + 12 (1.20) = 16 known issues across all entries
    assert.dom(GENERAL.badge('Known issues')).hasText('16', 'known issues badge count is correct');
    // 1 (1.21) + 5 (1.20) = 6 breaking changes
    assert.dom(GENERAL.badge('Breaking changes')).hasText('6', 'breaking changes badge count is correct');
    // 1 (1.21) + 4 (1.20) = 5 new behavior entries
    assert.dom(GENERAL.badge('New behavior')).hasText('5', 'new behavior badge count is correct');
    // rollback_steps is empty — the count badge is not rendered when count is falsy
    assert
      .dom(GENERAL.badge('Rollback steps'))
      .doesNotExist('rollback steps badge is hidden when count is 0');

    // The first visible panel item is the first known issue (1.21 entry)
    assert.dom(`[data-test-panel-item] ${GENERAL.badge()}`).exists();
    assert
      .dom('[data-test-panel-item] [data-test-panel-item-title]')
      .hasText('Missed events with multiple event clients');
    assert.dom('[data-test-panel-item] [data-test-panel-item-description]').hasText('Found in 1.21.0');
    assert.dom(`[data-test-panel-item] ${GENERAL.linkTo('Item details')}`).exists();
  });

  test('it paginates the panels', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradeInfo @breadcrumbs={{this.breadcrumbs}} @upgradeInfo={{this.upgradeInfo}}/>`
    );
    assert.dom(GENERAL.paginationInfo).hasText('1–10 of 16');
    await click(GENERAL.nextPage);
    assert.dom(GENERAL.paginationInfo).hasText('11–16 of 16');
  });
});
