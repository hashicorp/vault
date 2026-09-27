/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { click, currentURL, fillIn, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { login } from 'vault/tests/helpers/auth/auth-helpers';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import engineDisplayData from 'vault/helpers/engines-display-data';

// Two stub engines that should appear in the list (shouldIncludeInList = true for non-internal types)
const STUB_ENGINES = {
  'kv/': {
    type: 'kv',
    path: 'kv/',
    accessor: 'kv_abc123',
    description: 'A KV secrets engine',
    options: { version: 2 },
    running_plugin_version: 'v0.14.0',
    plugin_version: '',
    running_sha256: '',
    local: false,
    seal_wrap: false,
    external_entropy_access: false,
    config: {},
    uuid: 'uuid-1',
  },
  'aws/': {
    type: 'aws',
    path: 'aws/',
    accessor: 'aws_def456',
    description: 'An AWS secrets engine',
    options: {},
    running_plugin_version: 'v0.15.0',
    plugin_version: '',
    running_sha256: '',
    local: false,
    seal_wrap: false,
    external_entropy_access: false,
    config: {},
    uuid: 'uuid-2',
  },
};

module('Acceptance | secrets backends list view', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    await login();
    // Dismiss the wizard so it does not block the list view
    this.owner.lookup('service:wizard').dismiss(WIZARD_ID_MAP.secretEngines);
    this.server.get('/sys/internal/ui/mounts', () => ({
      data: { secret: STUB_ENGINES },
    }));
  });

  // ── Page structure ───────────────────────────────────────────────────────

  test('renders page title', async function (assert) {
    await visit('/vault/secrets-engines');
    assert.dom(GENERAL.title).hasText('Secrets engines');
  });

  test('renders breadcrumbs', async function (assert) {
    await visit('/vault/secrets-engines');
    assert.dom(GENERAL.breadcrumbLink('Vault')).exists('first breadcrumb links to Vault dashboard');
    assert.dom(GENERAL.currentBreadcrumb('Secrets engines')).exists('last breadcrumb is current page');
  });

  test('renders primary action button', async function (assert) {
    await visit('/vault/secrets-engines');
    assert.dom(GENERAL.button('Enable new engine')).exists();
  });

  // ── Table columns ────────────────────────────────────────────────────────

  test('renders expected column headers', async function (assert) {
    await visit('/vault/secrets-engines');
    // Column 1: Engine Type (index 1, 1-based)
    assert.dom(GENERAL.tableColumnHeader(1)).hasText('Engine Type');
    // Column 2: Engine path
    assert.dom(GENERAL.tableColumnHeader(2)).hasText('Engine path');
    // Column 3: Accessor (has tooltip icon)
    assert.dom(GENERAL.tableColumnHeader(3)).includesText('Accessor');
    // Column 4: Description
    assert.dom(GENERAL.tableColumnHeader(4)).hasText('Description');
    // Column 5: Version
    assert.dom(GENERAL.tableColumnHeader(5)).hasText('Version');
  });

  // ── Filter ───────────────────────────────────────────────────────────────

  test('filter input is rendered', async function (assert) {
    await visit('/vault/secrets-engines');
    assert.dom(GENERAL.filterInput).exists('filter input is rendered');
  });

  test('filter shows filtered empty state when no match', async function (assert) {
    await visit('/vault/secrets-engines');
    await fillIn(GENERAL.filterInput, 'nonexistent-engine-path');
    assert.dom(GENERAL.emptyStateTitle).includesText('No results for');
  });

  // ── Empty state ──────────────────────────────────────────────────────────

  test('shows empty state when list is empty', async function (assert) {
    this.server.get('/sys/internal/ui/mounts', () => ({ data: { secret: {} } }));
    await visit('/vault/secrets-engines');
    assert.dom(GENERAL.emptyStateTitle).hasText('No secrets engines');
  });

  // ── Row actions ──────────────────────────────────────────────────────────

  test('popup menu contains View configuration and Delete engine path', async function (assert) {
    await visit('/vault/secrets-engines');
    await click(GENERAL.menuTrigger);
    assert.dom(GENERAL.menuItem('view-configuration')).exists('View configuration action exists');
    assert.dom(GENERAL.menuItem('delete-engine-path')).exists('Delete engine path action exists');
  });

  test('confirm modal appears when Delete engine path is clicked', async function (assert) {
    await visit('/vault/secrets-engines');
    await click(GENERAL.menuTrigger);
    await click(GENERAL.menuItem('delete-engine-path'));
    assert.dom(GENERAL.confirmTitle).hasText('Delete engine path?');
  });

  // ── URL query params ─────────────────────────────────────────────────────

  test('advances the page query param when the next-page button is clicked', async function (assert) {
    // Verifies that clicking next updates the URL QP via Page::ListView onPageChange.
    // Stub > DEFAULT_PAGE_SIZE (15) engines so the pagination control renders.
    const manyEngines = {};
    for (let i = 0; i < 20; i++) {
      manyEngines[`kv-${i}/`] = {
        type: 'kv',
        path: `kv-${i}/`,
        accessor: `kv_${i}`,
        description: '',
        options: { version: 2 },
        running_plugin_version: 'v0.14.0',
        plugin_version: '',
        running_sha256: '',
        local: false,
        seal_wrap: false,
        external_entropy_access: false,
        config: {},
        uuid: `uuid-${i}`,
      };
    }
    this.server.get('/sys/internal/ui/mounts', () => ({ data: { secret: manyEngines } }));
    await visit('/vault/secrets-engines');
    await click(GENERAL.nextPage);
    assert.true(currentURL().includes('page=2'), 'URL contains page=2 after clicking next');
  });

  test('page query param resets to 1 on navigation away and back', async function (assert) {
    await visit('/vault/secrets-engines?page=2');
    // navigate away then back — resetController should reset page
    await visit('/vault/dashboard');
    await visit('/vault/secrets-engines');
    assert.false(currentURL().includes('page=2'), 'page param was reset');
  });

  // ── Dark mode icon resolution ─────────────────────────────────────────────

  module('engine type dropdown icons in dark mode', function (hooks) {
    hooks.beforeEach(function () {
      this.theme = this.owner.lookup('service:theme');
    });

    hooks.afterEach(function () {
      this.theme.setTheme('light');
    });

    test('engine type dropdown shows monochrome icons in dark mode', async function (assert) {
      // Verify that -color glyphs are stripped for aws and kv engine types.
      const awsGlyph = engineDisplayData('aws')?.glyph ?? 'lock';
      const kvGlyph = engineDisplayData('kv')?.glyph ?? 'lock';

      this.theme.setTheme('dark');
      await visit('/vault/secrets-engines');
      await click(GENERAL.toggleInput('filter-by-engine-type'));

      // In dark mode the -color suffix should be stripped.
      const expectedAws = awsGlyph.replace(/-color$/, '');
      const expectedKv = kvGlyph.replace(/-color$/, '');

      assert
        .dom(`[data-test-checkbox="aws"] [data-test-icon="${expectedAws}"]`)
        .exists(`aws engine uses monochrome icon "${expectedAws}" in dark mode`);
      assert
        .dom(`[data-test-checkbox="kv"] [data-test-icon="${expectedKv}"]`)
        .exists(`kv engine uses monochrome icon "${expectedKv}" in dark mode`);
    });

    test('engine type dropdown shows colour icons in light mode', async function (assert) {
      const awsGlyph = engineDisplayData('aws')?.glyph ?? 'lock';
      const kvGlyph = engineDisplayData('kv')?.glyph ?? 'lock';

      this.theme.setTheme('light');
      await visit('/vault/secrets-engines');
      await click(GENERAL.toggleInput('filter-by-engine-type'));

      assert
        .dom(`[data-test-checkbox="aws"] [data-test-icon="${awsGlyph}"]`)
        .exists(`aws engine uses colour icon "${awsGlyph}" in light mode`);
      assert
        .dom(`[data-test-checkbox="kv"] [data-test-icon="${kvGlyph}"]`)
        .exists(`kv engine uses colour icon "${kvGlyph}" in light mode`);
    });
  });
});
