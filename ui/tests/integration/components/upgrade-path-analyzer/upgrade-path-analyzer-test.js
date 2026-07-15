/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

const SELECTORS = {
  cardDescription: `[data-test-card-description]`,
  cardTitle: `[data-test-card-title]`,
};

module('Integration | Component | Upgrade Path Analyzer', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer' },
    ];
    this.version = this.owner.lookup('service:version');
    this.version.version = '1.19.5';
    this.onSetUpgradeInfo = () => {};
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

  test('it detects the current version', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    assert.dom(GENERAL.selectByAttr('1.19.5')).exists('Current version is detected');
  });

  test('it displays the Known issues card with correct title, description, count, and link', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    await click(GENERAL.button('Analyze'));

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
      .hasText('16', 'Known issues count is correct');
    assert.dom(GENERAL.linkTo('Known issues')).exists('Known issues view link exists');
  });

  test('it displays the Breaking changes card with correct title, description, count, and link', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    await click(GENERAL.button('Analyze'));

    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${SELECTORS.cardTitle}`)
      .hasText('Breaking changes', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${SELECTORS.cardDescription}`)
      .hasText('These are functional changes from one version to the other.', 'Card description is correct');
    assert
      .dom(`${GENERAL.cardContainer('Breaking changes')} ${GENERAL.badge()}`)
      .hasText('6', 'Breaking changes count is correct');
    assert.dom(GENERAL.linkTo('Breaking changes')).exists('Breaking changes view link exists');
  });

  test('it displays the New behavior card with correct title, description, count, and link', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    await click(GENERAL.button('Analyze'));

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
      .hasText('5', 'New behavior count is correct');
    assert.dom(GENERAL.linkTo('New behavior')).exists('New behavior view link exists');
  });

  test('it displays the Rollback steps card with correct title, description, count, and link', async function (assert) {
    await render(
      hbs`<UpgradePathAnalyzer::UpgradePathAnalyzer @breadcrumbs={{this.breadcrumbs}} @onSetUpgradeInfo={{this.onSetUpgradeInfo}}/>`
    );
    await click(GENERAL.button('Analyze'));

    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${SELECTORS.cardTitle}`)
      .hasText('Rollback steps', 'Card title is correct');
    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${SELECTORS.cardDescription}`)
      .hasText('Follow these steps to safely rollback.', 'Card description is correct');
    assert
      .dom(`${GENERAL.cardContainer('Rollback steps')} ${GENERAL.badge()}`)
      .hasText('0', 'Rollback steps count is correct');
    assert.dom(GENERAL.linkTo('Rollback steps')).exists('Rollback steps view link exists');
  });
});
