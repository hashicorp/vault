/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | policy-allowed-actions', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders with aggregated policy data in card mode', async function (assert) {
    this.aggregatedPolicy = {
      policy: {
        'secrets/ci/*': ['read', 'list'],
        'secrets/db/*': ['read', 'list', 'update', 'delete'],
        'secrets/kv/config-*': ['read'],
        'pki/issue/admin-ci': ['read', 'update'],
      },
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{true}} />`);

    assert.dom(GENERAL.cardContainer('policy-allowed-actions')).exists('Component renders');
    assert
      .dom(GENERAL.cardContainer('policy-allowed-actions'))
      .includesText('Allowed actions', 'Title is correct');

    // Check that all paths are rendered
    assert.dom('[data-test-path-entry]').exists({ count: 4 }, 'All 4 paths are rendered');
    assert.dom('[data-test-path="secrets/ci/*"]').hasText('secrets/ci/*', 'First path is correct');
    assert.dom('[data-test-path="secrets/db/*"]').hasText('secrets/db/*', 'Second path is correct');
  });

  test('it renders without card when isCard is false', async function (assert) {
    this.aggregatedPolicy = {
      policy: {
        'test/path': ['read'],
      },
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{false}} />`);

    assert.dom(GENERAL.cardContainer('policy-allowed-actions')).exists('Component renders');
    assert
      .dom(GENERAL.cardContainer('policy-allowed-actions'))
      .doesNotIncludeText('Allowed actions', 'Title is not shown in non-card mode');
    assert.dom('[data-test-path-entry]').exists('Path entry is rendered');
  });

  test('it displays capabilities with correct colors', async function (assert) {
    this.aggregatedPolicy = {
      policy: {
        'test/path': ['read', 'update', 'delete'],
      },
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{true}} />`);

    const pathEntry = '[data-test-path-entry="test/path"]';

    // Check that capabilities are rendered
    assert.dom(`${pathEntry} [data-test-capability="read"]`).exists('Read capability exists');
    assert.dom(`${pathEntry} [data-test-capability="update"]`).exists('Update capability exists');
    assert.dom(`${pathEntry} [data-test-capability="delete"]`).exists('Delete capability exists');

    // Check badge text formatting (capitalized)
    assert.dom(`${pathEntry} [data-test-capability="read"]`).hasText('Read', 'Read is capitalized');
    assert.dom(`${pathEntry} [data-test-capability="update"]`).hasText('Update', 'Update is capitalized');
    assert.dom(`${pathEntry} [data-test-capability="delete"]`).hasText('Delete', 'Delete is capitalized');
  });

  test('it sorts paths alphabetically', async function (assert) {
    this.aggregatedPolicy = {
      policy: {
        'z/path': ['read'],
        'a/path': ['read'],
        'm/path': ['read'],
      },
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{true}} />`);

    const paths = this.element.querySelectorAll('[data-test-path]');
    assert.strictEqual(paths[0].textContent.trim(), 'a/path', 'First path is a/path');
    assert.strictEqual(paths[1].textContent.trim(), 'm/path', 'Second path is m/path');
    assert.strictEqual(paths[2].textContent.trim(), 'z/path', 'Third path is z/path');
  });

  test('it shows empty state when no paths', async function (assert) {
    this.aggregatedPolicy = {
      policy: {},
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{true}} />`);

    assert.dom('[data-test-empty-state]').exists('Empty state is shown');
    assert
      .dom('[data-test-empty-state]')
      .hasText('No policy paths to display', 'Empty state message is correct');
    assert.dom('[data-test-path-entry]').doesNotExist('No path entries are shown');
  });

  test('it handles all capability types', async function (assert) {
    this.aggregatedPolicy = {
      policy: {
        'test/all': ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'],
      },
      policyString: '',
    };

    await render(hbs`<PolicyAllowedActions @aggregatedPolicy={{this.aggregatedPolicy}} @isCard={{true}} />`);

    const capabilities = ['create', 'read', 'update', 'delete', 'list', 'patch', 'sudo'];
    capabilities.forEach((cap) => {
      assert.dom(`[data-test-capability="${cap}"]`).exists(`${cap} capability is rendered`);
    });
  });
});
