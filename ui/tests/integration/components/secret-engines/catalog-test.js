/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click, fillIn } from '@ember/test-helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { filterEnginesByMountCategory } from 'vault/utils/all-engines-metadata';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

import hbs from 'htmlbars-inline-precompile';
import sinon from 'sinon';

const SELECTORS = {
  badgeLegend: '[data-test-badge-legend]',
  searchInput: GENERAL.inputSearch('search by keywords'),
  secretTypeToggle: GENERAL.toggleInput('filter-by-secret-type'),
  rotationTypeToggle: GENERAL.toggleInput('filter-by-rotation-type'),
  platformToggle: GENERAL.toggleInput('filter-by-platform'),
  checkmark: (value) => `[data-test-checkmark="${value}"]`,
  filterTag: (value) => GENERAL.button(value),
  keywordTag: GENERAL.button('keyword-tag'),
  clearAllButton: GENERAL.button('Clear all'),
};

module('Integration | Component | secret-engines/catalog', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.setMountType = sinon.spy();
    this.pluginCatalogData = null;
    this.pluginCatalogError = false;
  });

  test('it renders secret engines catalog', async function (assert) {
    const expectedEngines = filterEnginesByMountCategory({
      mountCategory: 'secret',
      isEnterprise: false,
    }).filter((engine) => engine.type !== 'cubbyhole');

    // Dynamic assertion count: 1 for title + number of engines
    assert.expect(1 + expectedEngines.length);

    await render(
      hbs`<SecretEngines::Catalog
        @setMountType={{this.setMountType}}
        @pluginCatalogData={{this.pluginCatalogData}}
        @pluginCatalogError={{this.pluginCatalogError}}
      />`
    );

    assert.dom(GENERAL.breadcrumbs).exists('renders breadcrumbs');

    for (const engine of expectedEngines) {
      assert.dom(GENERAL.cardContainer(engine.type)).exists(`renders ${engine.displayName} engine card`);
    }
  });

  test('it calls setMountType when engine is selected', async function (assert) {
    await render(
      hbs`<SecretEngines::Catalog
        @setMountType={{this.setMountType}}
        @pluginCatalogData={{this.pluginCatalogData}}
        @pluginCatalogError={{this.pluginCatalogError}}
      />`
    );

    await click(GENERAL.cardContainer('kv'));
    await click(GENERAL.button('next'));

    assert.true(this.setMountType.calledOnce, 'setMountType was called');
    assert.true(this.setMountType.calledWith('kv'), 'setMountType was called with kv');
  });

  test('it shows plugin catalog error when provided', async function (assert) {
    this.pluginCatalogError = true;

    await render(
      hbs`<SecretEngines::Catalog
        @setMountType={{this.setMountType}}
        @pluginCatalogData={{this.pluginCatalogData}}
        @pluginCatalogError={{this.pluginCatalogError}}
      />`
    );

    assert.dom(GENERAL.inlineAlert).exists('shows plugin catalog error alert');
    assert
      .dom(GENERAL.inlineAlert)
      .hasText(
        'Plugin information unavailable Unable to fetch current plugin information. Using static plugin data instead. Some plugins may not show current details.',
        'shows correct error title'
      );
  });

  test('it shows flyout when clicking disabled plugin', async function (assert) {
    // Set up plugin catalog data that creates both enabled and disabled engines
    // An engine is disabled when it's not found in the plugin catalog detailed array
    this.pluginCatalogData = {
      detailed: [
        // Include only some engines, leaving others as "disabled"
        {
          name: 'kv',
          type: 'secret',
          builtin: true,
          deprecation_status: 'supported',
          version: 'v1.0.0',
        },
        // AWS engine is NOT included, so it will be marked as isAvailable: false
      ],
    };

    await render(
      hbs`<SecretEngines::Catalog
        @setMountType={{this.setMountType}}
        @pluginCatalogData={{this.pluginCatalogData}}
        @pluginCatalogError={{this.pluginCatalogError}}
      />`
    );

    // Initially, flyout should not be visible
    assert.dom(GENERAL.flyout).doesNotExist('flyout is not shown initially');

    // Find a disabled plugin card - since AWS is not in our catalog data,
    // it should be rendered as disabled
    const awsCard = document.querySelector(GENERAL.cardContainer('aws'));

    // Look for any disabled cards regardless of AWS card presence
    const disabledCards = document.querySelectorAll(
      '.selectable-engine-card.disabled, .selectable-engine-card[style*="opacity"]'
    );

    let clickedCard = false;

    if (awsCard) {
      await click(awsCard);
      clickedCard = true;

      // After clicking disabled plugin, flyout should appear
      assert.dom(GENERAL.flyout).exists('flyout appears after clicking disabled plugin');
    } else if (disabledCards.length > 0) {
      await click(disabledCards[0]);
      clickedCard = true;
      assert.dom(GENERAL.flyout).exists('flyout appears after clicking any disabled plugin');
    }

    // Always verify we completed the test successfully
    assert.ok(clickedCard, 'successfully clicked a disabled plugin card');
  });

  test('it toggles the badge legend and displays all entries with correct descriptions', async function (assert) {
    await render(
      hbs`<SecretEngines::Catalog
          @setMountType={{this.setMountType}}
          @pluginCatalogData={{this.pluginCatalogData}}
          @pluginCatalogError={{this.pluginCatalogError}}
        />`
    );

    assert.dom(GENERAL.revealButton('badge legend')).exists('badge legend toggle button is rendered');
    assert.dom(SELECTORS.badgeLegend).doesNotExist('badge legend is hidden by default');

    await click(GENERAL.revealButton('badge legend'));

    assert.dom(SELECTORS.badgeLegend).exists('badge legend is visible after clicking toggle');

    const entries = [
      { badge: 'Static', description: 'Can store static credentials' },
      { badge: 'Rotating', description: 'Can perform auto-rotation' },
      { badge: 'Dynamic', description: 'Can generate temporary credentials' },
      { badge: 'Encryption', description: 'Can perform encryption as a service' },
      { badge: 'Signing', description: 'Can sign keys' },
      { badge: 'Tokenization', description: 'Can tokenize sensitive data' },
      { badge: 'Certificate authority', description: 'Can issue and manage certificates' },
    ];

    for (const { badge, description } of entries) {
      assert.dom(GENERAL.badge(badge)).exists(`${badge} badge is shown`);
      assert.dom(GENERAL.textBody(badge)).hasText(description, `${badge} description is correct`);
    }

    await click(GENERAL.revealButton('badge legend'));
    assert.dom(SELECTORS.badgeLegend).doesNotExist('badge legend is hidden after toggling closed');
  });

  // --- Keyword search ---

  test('keyword search filters engines by display name', async function (assert) {
    assert.expect(3);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    assert.dom('[data-test-no-filters-applied]').exists('starts with no filters applied');

    await fillIn(SELECTORS.searchInput, 'transit');

    assert.dom(GENERAL.cardContainer('transit')).exists('Transit engine card is shown');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV engine card is hidden');
  });

  test('keyword search matches against engine type', async function (assert) {
    assert.expect(2);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await fillIn(SELECTORS.searchInput, 'ssh');

    assert.dom(GENERAL.cardContainer('ssh')).exists('SSH engine card is shown for type match');
    assert.dom(GENERAL.cardContainer('transit')).doesNotExist('Transit card is hidden');
  });

  test('keyword search matches against engine description', async function (assert) {
    assert.expect(2);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    // 'post-quantum' appears only in Transit's description
    await fillIn(SELECTORS.searchInput, 'post-quantum');

    assert.dom(GENERAL.cardContainer('transit')).exists('Transit card shown for description match');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV card hidden when description does not match');
  });

  // --- Secret type filter ---

  test('secret type filter - encryption keys shows only engines with that secretType', async function (assert) {
    assert.expect(5);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('encryptionKeys'));

    assert.dom(SELECTORS.filterTag('encryptionKeys')).exists('active filter tag appears for encryption keys');
    assert.dom(GENERAL.cardContainer('transit')).exists('Transit (encryption keys) is shown');
    assert.dom(GENERAL.cardContainer('gcpkms')).exists('GCP KMS (encryption keys) is shown');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV (static storage only) is hidden');
    assert.dom(GENERAL.cardContainer('database')).doesNotExist('Databases (database credentials) is hidden');
  });

  test('secret type filter - database credentials shows only the database engine', async function (assert) {
    assert.expect(4);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('databaseCredentials'));

    assert
      .dom(SELECTORS.filterTag('databaseCredentials'))
      .exists('active filter tag appears for database credentials');
    assert.dom(GENERAL.cardContainer('database')).exists('Databases engine is shown');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV is hidden');
    assert.dom(GENERAL.cardContainer('transit')).doesNotExist('Transit is hidden');
  });

  test('secret type filter - cloud credentials shows only engines with that secretType', async function (assert) {
    assert.expect(5);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('cloudCredentials'));

    assert
      .dom(SELECTORS.filterTag('cloudCredentials'))
      .exists('active filter tag appears for cloud credentials');
    assert.dom(GENERAL.cardContainer('aws')).exists('AWS (cloud credentials) is shown');
    assert.dom(GENERAL.cardContainer('azure')).exists('Azure (cloud credentials) is shown');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV (static storage) is hidden');
    assert.dom(GENERAL.cardContainer('database')).doesNotExist('Databases (database credentials) is hidden');
  });

  // --- Rotation type filter ---

  test('rotation type filter shows only engines with matching capability', async function (assert) {
    assert.expect(5);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.rotationTypeToggle);
    await click(SELECTORS.checkmark('rotating'));

    assert.dom(SELECTORS.filterTag('rotating')).exists('active filter tag appears for rotating');
    assert.dom(GENERAL.cardContainer('aws')).exists('AWS (rotating) is shown');
    assert.dom(GENERAL.cardContainer('database')).exists('Databases (rotating) is shown');
    assert.dom(GENERAL.cardContainer('kv')).doesNotExist('KV (static only) is hidden');
    assert.dom(GENERAL.cardContainer('transit')).doesNotExist('Transit (no rotating capability) is hidden');
  });

  // --- Platform filter ---

  test('platform filter shows only engines in the selected category', async function (assert) {
    assert.expect(5);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.platformToggle);
    await click(SELECTORS.checkmark('common engines'));

    assert.dom(SELECTORS.filterTag('common engines')).exists('active filter tag appears for common engines');
    assert.dom(GENERAL.cardContainer('kv')).exists('KV (common engines) is shown');
    assert.dom(GENERAL.cardContainer('aws')).exists('AWS (common engines) is shown');
    assert.dom(GENERAL.cardContainer('transit')).doesNotExist('Transit (crypto category) is hidden');
    assert.dom(GENERAL.cardContainer('ldap')).doesNotExist('LDAP (identity category) is hidden');
  });

  // --- Combined filters ---

  test('combined keyword and secret type filter narrows results further', async function (assert) {
    assert.expect(2);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('encryptionKeys'));
    await fillIn(SELECTORS.searchInput, 'transit');

    assert.dom(GENERAL.cardContainer('transit')).exists('Transit matches both filters');
    assert.dom(GENERAL.cardContainer('gcpkms')).doesNotExist('GCP KMS hidden by keyword filter');
  });

  // --- Active filter tags ---

  test('active filter tags show human-readable labels', async function (assert) {
    assert.expect(2);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('encryptionKeys'));

    assert.dom(SELECTORS.filterTag('encryptionKeys')).exists('active filter tag is rendered');
    assert
      .dom(SELECTORS.filterTag('encryptionKeys'))
      .hasText('Encryption keys', 'tag displays human-readable label, not raw value');
  });

  test('active keyword tag shown and dismissed correctly', async function (assert) {
    assert.expect(4);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await fillIn(SELECTORS.searchInput, 'transit');

    assert.dom(SELECTORS.keywordTag).exists('keyword filter tag appears');
    assert.dom(SELECTORS.keywordTag).hasText('transit', 'keyword tag text matches search input');

    await click(`${SELECTORS.keywordTag} button`);

    assert.dom(SELECTORS.keywordTag).doesNotExist('keyword tag removed after dismiss');
    assert.dom(GENERAL.cardContainer('kv')).exists('KV engine visible again after clearing keyword');
  });

  test('dismissing a filter tag removes only that filter', async function (assert) {
    assert.expect(3);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('encryptionKeys'));
    await click(SELECTORS.rotationTypeToggle);
    await click(SELECTORS.checkmark('static'));

    await click(`${SELECTORS.filterTag('encryptionKeys')} button`);

    assert.dom(SELECTORS.filterTag('encryptionKeys')).doesNotExist('encryption keys tag removed');
    assert.dom(SELECTORS.filterTag('static')).exists('static rotation tag remains');
    assert
      .dom(GENERAL.cardContainer('kv'))
      .exists('KV (static capability) visible after encryption keys dismissed');
  });

  // --- No filters / clear all ---

  test('"No filters applied" shown when no filters are active', async function (assert) {
    assert.expect(1);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    assert
      .dom('[data-test-no-filters-applied]')
      .hasText('No filters applied.', 'shows "No filters applied." when no filters active');
  });

  test('"Clear all" resets all active filters', async function (assert) {
    assert.expect(5);
    await render(hbs`<SecretEngines::Catalog
      @setMountType={{this.setMountType}}
      @pluginCatalogData={{this.pluginCatalogData}}
      @pluginCatalogError={{this.pluginCatalogError}}
    />`);

    await fillIn(SELECTORS.searchInput, 'transit');
    await click(SELECTORS.secretTypeToggle);
    await click(SELECTORS.checkmark('encryptionKeys'));
    await click(SELECTORS.rotationTypeToggle);
    await click(SELECTORS.checkmark('rotating'));

    assert.dom(SELECTORS.clearAllButton).exists('"Clear all" button is shown');

    await click(SELECTORS.clearAllButton);

    assert.dom(SELECTORS.keywordTag).doesNotExist('keyword tag cleared');
    assert.dom(SELECTORS.filterTag('encryptionKeys')).doesNotExist('secret type filter tag cleared');
    assert.dom(SELECTORS.filterTag('rotating')).doesNotExist('rotation type filter tag cleared');
    assert.dom(GENERAL.cardContainer('kv')).exists('all engines visible again after clear all');
  });
});
