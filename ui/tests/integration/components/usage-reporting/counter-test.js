/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/counter', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders the count', async function (assert) {
    this.set('count', 42);
    await render(hbs`<UsageReporting::Counter @title="Child namespaces" @count={{this.count}} />`);
    assert
      .dom('[data-test-vault-reporting-counter="Child namespaces"]')
      .exists('counter element is rendered')
      .includesText('42', 'renders the count value');
  });

  test('it renders with a suffix', async function (assert) {
    await render(hbs`
      <UsageReporting::Counter @title="PKI roles" @count={{5}} @suffix="roles" />
    `);
    assert
      .dom('[data-test-vault-reporting-counter="PKI roles"]')
      .includesText('5 roles', 'renders count with suffix');
  });

  test('it renders empty text when count is 0 and @emptyText is provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Counter
        @title="KV secrets"
        @count={{0}}
        @emptyText="No secrets stored"
      />
    `);
    assert
      .dom('[data-test-vault-reporting-counter="KV secrets"]')
      .includesText('No secrets stored', 'renders emptyText when count is 0');
  });

  test('it renders numeric count even when @emptyText is provided and count is non-zero', async function (assert) {
    await render(hbs`
      <UsageReporting::Counter
        @title="KV secrets"
        @count={{10}}
        @emptyText="No secrets stored"
      />
    `);
    assert
      .dom('[data-test-vault-reporting-counter="KV secrets"]')
      .includesText('10', 'renders count when non-zero');
  });

  test('it renders a tooltip when @tooltipMessage is provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Counter
        @title="Child namespaces"
        @count={{3}}
        @tooltipMessage="Total number of namespaces for this cluster."
      />
    `);
    assert.dom('[data-test-vault-reporting-counter-tooltip-button]').exists('tooltip button is rendered');
  });

  test('it does not render a tooltip when @tooltipMessage is not provided', async function (assert) {
    await render(hbs`<UsageReporting::Counter @title="PKI roles" @count={{0}} />`);
    assert.dom('[data-test-vault-reporting-counter-tooltip-button]').doesNotExist();
  });

  test('it renders a link when @link is provided', async function (assert) {
    await render(hbs`
      <UsageReporting::Counter
        @title="Child namespaces"
        @count={{5}}
        @link="vault.cluster.access.namespaces"
      />
    `);
    assert.dom('[data-test-vault-reporting-counter="Child namespaces"] a').exists('link is rendered');
  });
});
