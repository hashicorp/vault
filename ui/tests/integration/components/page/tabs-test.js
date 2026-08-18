/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import sinon from 'sinon';

module('Integration | Component | page/tabs', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.tabs = [{ label: 'Details' }, { label: 'Settings' }];
    this.onClickTab = sinon.spy();
    this.selectedTabIndex = 0;
    this.renderComponent = () =>
      render(hbs`
        <Page::Tabs
          @tabs={{this.tabs}}
          @onClickTab={{this.onClickTab}}
          @selectedTabIndex={{this.selectedTabIndex}}
          as |T|
        >
          <T.Panel data-test-panel="Details">Details content</T.Panel>
          <T.Panel data-test-panel="Settings">Settings content</T.Panel>
        </Page::Tabs>
      `);
  });

  test('it renders a tab button for each entry in @tabs', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('Details')).exists('renders Details tab');
    assert.dom(GENERAL.hdsTab('Settings')).exists('renders Settings tab');
  });

  test('it uses tab.key as data-test-tab when provided', async function (assert) {
    this.tabs = [
      { key: 'details-key', label: 'Details' },
      { key: 'settings-key', label: 'Settings' },
    ];
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('details-key')).exists('renders tab with key as data-test-tab');
    assert.dom(GENERAL.hdsTab('settings-key')).exists('renders tab with key as data-test-tab');
    assert.dom(GENERAL.hdsTab('Details')).doesNotExist('does not use label when key is present');
  });

  test('it falls back to tab.label for data-test-tab when key is absent', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('Details')).exists('falls back to label for data-test-tab');
    assert.dom(GENERAL.hdsTab('Settings')).exists('falls back to label for data-test-tab');
  });

  test('it renders tab.label as the visible tab text', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsTab('Details')).hasText('Details', 'Details tab has correct label text');
    assert.dom(GENERAL.hdsTab('Settings')).hasText('Settings', 'Settings tab has correct label text');
  });

  test('it renders an icon when tab.icon is provided', async function (assert) {
    this.tabs = [{ label: 'Details', icon: 'info' }, { label: 'Settings' }];
    await this.renderComponent();

    assert.dom(`${GENERAL.hdsTab('Details')} [data-test-icon="info"]`).exists('renders icon in tab');
    assert
      .dom(`${GENERAL.hdsTab('Settings')} [data-test-icon]`)
      .doesNotExist('does not render icon when tab.icon is absent');
  });

  test('it renders a badge count when tab.count is provided', async function (assert) {
    this.tabs = [{ label: 'Details', count: '4' }, { label: 'Settings' }];
    await this.renderComponent();

    assert.dom(`${GENERAL.hdsTab('Details')} .hds-badge-count`).exists('renders badge count in tab');
    assert
      .dom(`${GENERAL.hdsTab('Settings')} .hds-badge-count`)
      .doesNotExist('does not render badge when tab.count is absent');
  });

  test('it yields the HDS contextual component so panel content renders', async function (assert) {
    await this.renderComponent();

    assert.dom(GENERAL.hdsTabPanel('Details')).exists('first panel content renders via yielded T');
    assert.dom(GENERAL.hdsTabPanel('Settings')).exists('second panel content renders via yielded T');
  });

  test('it calls @onClickTab when a tab is clicked', async function (assert) {
    await this.renderComponent();

    await click(GENERAL.hdsTab('Settings'));

    assert.true(this.onClickTab.calledOnce, '@onClickTab was called once');
  });

  test('it passes @selectedTabIndex to HDS to control the active tab', async function (assert) {
    this.selectedTabIndex = 1;
    await this.renderComponent();

    assert
      .dom(GENERAL.hdsTab('Settings'))
      .hasAttribute('aria-selected', 'true', 'second tab is selected when selectedTabIndex is 1');
    assert.dom(GENERAL.hdsTab('Details')).hasAttribute('aria-selected', 'false', 'first tab is not selected');
  });
});
