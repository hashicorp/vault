/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { UPGRADE_INFO } from 'vault/constants/upgrade-info';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | UpgradePathAnalyzer::UpgradeInfo', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer', route: 'vault.cluster.support.upgrade' },
      { label: 'Issues' },
    ];

    this.upgradeInfo = UPGRADE_INFO;
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

    // Check tab counts
    assert.dom(GENERAL.badge('Known issues')).hasText('16', 'badge count is correct');
    assert.dom(GENERAL.badge('Breaking changes')).hasText('6', 'badge count is correct');
    assert.dom(GENERAL.badge('New behavior')).hasText('5', 'badge count is correct');
    assert.dom(GENERAL.badge('Rollback steps')).doesNotExist();

    // Check issue description
    assert.dom(`[data-test-panel-item] ${GENERAL.badge()}`).exists();
    assert
      .dom('[data-test-panel-item] [data-test-panel-item-title]')
      .hasText('Missed events with multiple event clients');
    assert.dom('[data-test-panel-item] [data-test-panel-item-description]').hasText('Found in 1.21.0');
    assert.dom(`[data-test-panel-item] ${GENERAL.linkTo('Item details')}`).exists();
  });
});
