/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | secret-engine/list', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    // Clear localStorage before each test
    window.localStorage.clear();
  });

  hooks.afterEach(function () {
    // Clean up localStorage after each test
    window.localStorage.clear();
  });

  test('it displays favorite star button for each secret engine', async function (assert) {
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_123',
        path: 'secret/',
        type: 'kv',
        id: 'kv_123',
      },
    });

    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_456',
        path: 'pki/',
        type: 'pki',
        id: 'pki_456',
      },
    });

    this.secretEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngines}} />`);

    assert.dom('[data-test-favorite-engine="kv_123"]').exists('favorite button exists for kv engine');
    assert.dom('[data-test-favorite-engine="pki_456"]').exists('favorite button exists for pki engine');

    // Check initial state - should show empty star
    assert.dom('[data-test-favorite-engine="kv_123"] .hds-icon-star').exists('shows empty star initially');
  });

  test('it toggles favorite state when star button is clicked', async function (assert) {
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_123',
        path: 'secret/',
        type: 'kv',
        id: 'kv_123',
      },
    });

    this.secretEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngines}} />`);

    // Initially should show empty star
    assert.dom('[data-test-favorite-engine="kv_123"] .hds-icon-star').exists('shows empty star initially');

    // Click to add to favorites
    await click('[data-test-favorite-engine="kv_123"]');

    // Should now show filled star
    assert
      .dom('[data-test-favorite-engine="kv_123"] .hds-icon-star-fill')
      .exists('shows filled star after click');

    // Click again to remove from favorites
    await click('[data-test-favorite-engine="kv_123"]');

    // Should show empty star again
    assert
      .dom('[data-test-favorite-engine="kv_123"] .hds-icon-star')
      .exists('shows empty star after second click');
  });

  test('it persists favorites in localStorage', async function (assert) {
    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_123',
        path: 'secret/',
        type: 'kv',
        id: 'kv_123',
      },
    });

    this.secretEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngines}} />`);

    // Add to favorites
    await click('[data-test-favorite-engine="kv_123"]');

    // Check localStorage
    const stored = JSON.parse(localStorage.getItem('vault-favorite-engines') || '[]');
    assert.true(stored.includes('kv_123'), 'engine ID is stored in localStorage');

    // Remove from favorites
    await click('[data-test-favorite-engine="kv_123"]');

    // Check localStorage again
    const storedAfter = JSON.parse(localStorage.getItem('vault-favorite-engines') || '[]');
    assert.false(storedAfter.includes('kv_123'), 'engine ID is removed from localStorage');
  });

  test('it sorts favorites to the top', async function (assert) {
    // Pre-populate localStorage with one favorite
    localStorage.setItem('vault-favorite-engines', JSON.stringify(['pki_456']));

    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_123',
        path: 'secret/',
        type: 'kv',
        id: 'kv_123',
      },
    });

    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'pki_456',
        path: 'pki/',
        type: 'pki',
        id: 'pki_456',
      },
    });

    this.secretEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngines}} />`);

    // Get all secret engine links in order
    const links = this.element.querySelectorAll('[data-test-secrets-backend-link]');
    const firstLink = links[0];

    // The favorite (pki_456) should be first
    assert
      .dom(firstLink)
      .hasAttribute('data-test-secrets-backend-link', 'pki_456', 'favorite engine appears first');
  });

  test('it loads favorites from localStorage on component initialization', async function (assert) {
    // Pre-populate localStorage
    localStorage.setItem('vault-favorite-engines', JSON.stringify(['kv_123']));

    this.store.pushPayload('secret-engine', {
      modelName: 'secret-engine',
      data: {
        accessor: 'kv_123',
        path: 'secret/',
        type: 'kv',
        id: 'kv_123',
      },
    });

    this.secretEngines = this.store.peekAll('secret-engine', {});

    await render(hbs`<SecretEngine::List @secretEngines={{this.secretEngines}} />`);

    // Should show filled star since it's loaded from localStorage
    assert
      .dom('[data-test-favorite-engine="kv_123"] .hds-icon-star-fill')
      .exists('shows filled star for favorite loaded from localStorage');
  });
});
