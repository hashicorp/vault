/**
 * Copyright IBM Corp. 2025, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

module('Integration | Component | usage-reporting/global-lease', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders percentage and count text when quota is set', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{210000}} @quota={{420000}} />
    `);
    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('50%', 'renders correct percentage');
    assert
      .dom('[data-test-vault-reporting-global-lease-count-text]')
      .hasText('210K / 420K', 'renders compact formatted count and quota');
  });

  test('it renders 100% and meter chart when count meets quota', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{420000}} @quota={{420000}} />
    `);
    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('100%', 'renders 100% when at capacity');
    assert
      .dom('[data-test-vault-reporting-global-lease-meter-chart] [data-carbon-chart]')
      .exists('meter chart is rendered');
  });

  test('it does not render alert at 94%', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{94}} @quota={{100}} />
    `);

    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('94%', 'renders correct percentage at 94%');
    assert
      .dom('[data-test-vault-reporting-global-lease-alert]')
      .doesNotExist('alert is not shown below the 95% threshold');
  });

  test('it renders neutral alert at 95%', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{95}} @quota={{100}} />
    `);

    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('95%', 'renders correct percentage at 95%');
    assert
      .dom('[data-test-vault-reporting-global-lease-alert]')
      .exists('alert is shown at 95%')
      .includesText(
        'Approaching quota limit. Reduce usage or increase the lease limit to avoid blocking new leases.',
        'neutral threshold copy is shown at 95%'
      );
  });

  test('it renders neutral alert at 99%', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{99}} @quota={{100}} />
    `);

    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('99%', 'renders correct percentage at 99%');
    assert
      .dom('[data-test-vault-reporting-global-lease-alert]')
      .exists('alert is shown at 99%')
      .includesText(
        'Approaching quota limit. Reduce usage or increase the lease limit to avoid blocking new leases.',
        'neutral threshold copy is shown at 99%'
      );
  });

  test('it renders warning alert at 100%', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{100}} @quota={{100}} />
    `);

    assert
      .dom('[data-test-vault-reporting-global-lease-percentage-text]')
      .hasText('100%', 'renders correct percentage at 100%');
    assert
      .dom('[data-test-vault-reporting-global-lease-alert]')
      .exists('alert is shown at 100%')
      .includesText(
        'Global lease quota limit reached. If lease creation is blocked, reduce usage or increase the limit.',
        'warning threshold copy is shown at 100%'
      );
  });

  test('it renders the description link when quota is set', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease @count={{10}} @quota={{100}} />
    `);
    assert
      .dom('[data-test-vault-reporting-global-lease-description-link]')
      .exists('description leases link is present');
  });

  test('it renders the default empty state when quota is not set', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease />
    `);
    assert
      .dom('[data-test-vault-reporting-global-lease-empty-state]')
      .exists('empty state renders when no quota');
    assert
      .dom('[data-test-vault-reporting-global-lease-empty-state-description]')
      .exists('empty state description renders');
    assert
      .dom('[data-test-vault-reporting-global-lease-empty-state-link]')
      .exists('empty state docs link renders');
  });

  test('it renders a named block empty state', async function (assert) {
    await render(hbs`
      <UsageReporting::GlobalLease>
        <:empty as |A|>
          <A.Body @text="Custom lease empty message." />
        </:empty>
      </UsageReporting::GlobalLease>
    `);
    assert
      .dom('[data-test-vault-reporting-global-lease-empty-state]')
      .includesText('Custom lease empty message.', 'renders custom named block empty state');
    assert
      .dom('[data-test-vault-reporting-global-lease-empty-state-description]')
      .doesNotExist('default description is not rendered when named block is used');
  });
});
